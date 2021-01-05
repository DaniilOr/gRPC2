package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	serverbp "github.com/DaniilOr/gRPC2/pkg/server"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := serverbp.NewAgrigatorServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)
	t, err := time.Parse("2/1/2006", "28/09/2021")
	if err != nil {
		log.Println(err)
		return err
	}
	stream, err := client.Get(ctx, &serverbp.GetRequest{From: "Moscow", To: "Omsk", Date: timestamppb.New(t)})
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		log.Println(response)
	}
	log.Println("finished")
	return nil
}