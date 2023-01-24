package gapi

import (
	"context"

	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/pb"
	"github.com/ferseg/golang-simple-bank/util"
	"github.com/ferseg/golang-simple-bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
  violations := validateCreateUserRequest(req)
  if len(violations) > 0 {
    return nil, toInvalidArgumentError(violations)
  }
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not hash password")
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not create user")
	}
	userResponse := mapUserResponse(user)
  return &pb.CreateUserResponse{
    User: userResponse,
  }, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
  if err := val.ValidateUsername(req.GetUsername()); err != nil {
    violations = append(violations, toFieldViolation("username", err))
  }
  if err := val.ValidateFullName(req.GetFullName()); err != nil {
    violations = append(violations, toFieldViolation("fullName", err))
  }
  if err := val.ValidateEmail(req.GetEmail()); err != nil {
    violations = append(violations, toFieldViolation("email", err))
  }
  if err := val.ValidatePassword(req.GetPassword()); err != nil {
    violations = append(violations, toFieldViolation("password", err))
  }
  return
}

func toFieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
  return &errdetails.BadRequest_FieldViolation{
    Field: field,
    Description: err.Error(),
  }
}

func toInvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
  badRequest := &errdetails.BadRequest{FieldViolations: violations}
  statusInvalid := status.New(codes.InvalidArgument, "Invalid parameter")
  statusDetails, err := statusInvalid.WithDetails(badRequest)
  if err!=nil {
    return statusInvalid.Err()
  }
  return statusDetails.Err()
}

func toUnauthenticatedError(err error) error {
  return status.Errorf(codes.Unauthenticated, "Unauthorized: %s", err)
}
