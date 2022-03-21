package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"TestAll/grpc/greet"

	"google.golang.org/grpc"
)

type server struct {
	greet.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greet.HelloRequest) (*greet.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &greet.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 6789))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greet.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
