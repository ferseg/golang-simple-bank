package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func GrpcLogger(
  ctx context.Context, 
  req interface{}, 
  info *grpc.UnaryServerInfo, 
  handler grpc.UnaryHandler,
) (resp interface{}, err error) {
  startTime := time.Now()
  result, err := handler(ctx, req)
  duration := time.Since(startTime)
  log.Info().
    Str("protocol", "grpc").
    Str("method", info.FullMethod).
    Dur("duration", duration).
    Msg("Received a gRPC request")
  return result, err
}
