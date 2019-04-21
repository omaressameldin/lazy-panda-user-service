package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)

// RunServer runs service to publish User service
func RunServer(ctx context.Context, v1API v1.UserServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterUserServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down user server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting user server...")
	return server.Serve(listen)
}
