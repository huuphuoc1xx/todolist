package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	todoserver "todolist/TodolistServer/proto"

	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, API todoserver.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	todoserver.RegisterToDoServiceServer(server, API)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)

}
