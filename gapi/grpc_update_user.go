package gapi

import (
	"context"
	"database/sql"

	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/pb"
	"github.com/ferseg/golang-simple-bank/util"
	"github.com/ferseg/golang-simple-bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
  authPayload, err := server.authorizeUser(ctx)
  if err != nil {
    return nil, toUnauthenticatedError(err)
  }

	violations := validateUpdateUserRequest(req)
	if len(violations) > 0 {
		return nil, toInvalidArgumentError(violations)
	}
  
  if authPayload.Username != req.Username {
    return nil, status.Errorf(codes.PermissionDenied, "Cannot update other user's info %s", err)
  }

	arg := db.UpdateUserParams{
		Username: req.Username,
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not hash password")
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not create user")
	}
	userResponse := mapUserResponse(user)
	return &pb.UpdateUserResponse{
		User: userResponse,
	}, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateFullName(req.GetFullName()); err != nil && req.FullName != nil {
		violations = append(violations, toFieldViolation("fullName", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil && req.Email != nil {
		violations = append(violations, toFieldViolation("email", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil && req.Password != nil {
		violations = append(violations, toFieldViolation("password", err))
	}
	return
}
