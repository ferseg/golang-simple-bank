package gapi

import (
	"context"

	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/pb"
	"github.com/ferseg/golang-simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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
