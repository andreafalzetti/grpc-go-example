package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	recipepb "github.com/andreafalzetti/grpc-go-example/proto/recipe"
)

type server struct {
	recipepb.UnimplementedRecipesServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Get(ctx context.Context, in *recipepb.GetRequest) (*recipepb.GetResponse, error) {
	return &recipepb.GetResponse{Message: "hello"}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Attach the Greeter service to the server
	recipepb.RegisterRecipesServer(s, &server{})

	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatal(s.Serve(lis))
}
