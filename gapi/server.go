package gapi

import (
	"fmt"

	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/pb"
	"github.com/ferseg/golang-simple-bank/token"
	"github.com/ferseg/golang-simple-bank/util"
)

type Server struct {
  pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
  tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
  if err != nil {
    return nil, fmt.Errorf("Cannot create a new server")
  }

  server := &Server{
    config: config,
    store: store,
    tokenMaker: tokenMaker,
  }
  return server, nil
}
