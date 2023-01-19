package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/ferseg/golang-simple-bank/api"
	db "github.com/ferseg/golang-simple-bank/db/sqlc"
	"github.com/ferseg/golang-simple-bank/gapi"
	"github.com/ferseg/golang-simple-bank/pb"
	"github.com/ferseg/golang-simple-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Something has happend", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Could not create the server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot create listener")
	}
	
  err = grpcServer.Serve(listener)
	
  if err != nil {
		log.Fatal("Could not start server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	
  if err != nil {
		log.Fatal("Could not create the server", err)
	}
	
  err = server.Start(config.ServerAddress)
	
  if err != nil {
		log.Fatal("Cannot start the server", err)

	}
}
