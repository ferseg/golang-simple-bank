package gapi

import (
	"context"
	"database/sql"

	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/pb"
	"github.com/ferseg/golang-simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "The user doesn't exists")
		}
		return nil, status.Errorf(codes.Internal, "Could not get the user")
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Could not verify the password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		req.GetUsername(),
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating token")
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		// UserAgent:    ctx.Value,
		// ClientIp:     ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating session")
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		User:                  mapUserResponse(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}

func mapUserResponse(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
