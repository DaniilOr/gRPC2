package main

import (
	"context"
	"github.com/DaniilOr/gRPC2/cmd/flights/server/app"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	serverpb "github.com/DaniilOr/gRPC2/pkg/server"
	"log"
	"net"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const defaultDSN1 = "postgres://app:pass@localhost:5432/flights1"
const defaultDSN2 = "postgres://app:pass@localhost:5433/flights2"
const defaultDSN3 = "postgres://app:pass@localhost:5434/flights3"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}
	dsn1, ok := os.LookupEnv("APP_DSN1")
	if !ok {
		dsn1 = defaultDSN1
	}
	dsn2, ok := os.LookupEnv("APP_DSN2")
	if !ok {
		dsn2 = defaultDSN2
	}
	dsn3, ok := os.LookupEnv("APP_DSN3")
	if !ok {
		dsn3 = defaultDSN3
	}


	if err := execute(net.JoinHostPort(host, port), dsn1, dsn2, dsn3); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn1 string, dsn2 string, dsn3 string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	db1, err := pgxpool.Connect(context.Background(), dsn1)
	if err != nil {
		log.Println(err)
		return err
	}
	db2, err := pgxpool.Connect(context.Background(), dsn2)
	if err != nil {
		log.Println(err)
		return err
	}
	db3, err := pgxpool.Connect(context.Background(), dsn3)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx := context.Background()
	grpcServer := grpc.NewServer()
	server := app.NewServer(db1, db2, db3, ctx)
	serverpb.RegisterAgrigatorServiceServer(grpcServer, server)
	return grpcServer.Serve(listener)
}