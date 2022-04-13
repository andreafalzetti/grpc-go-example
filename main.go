package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	recipepb "github.com/andreafalzetti/grpc-go-example/proto/recipe"
)

type server struct {
	recipepb.UnimplementedRecipesServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Get(ctx context.Context, in *recipepb.GetRequest) (*recipepb.GetResponse, error) {
	var recipes = []*recipepb.Recipe{
		{
			Id:   0,
			Name: "Carbonara Pasta",
		},
		{
			Id:   1,
			Name: "Chicken Tikka Masala",
		},
		{
			Id:   2,
			Name: "Potato rösti cakes with sage leaves",
		},
		{
			Id:   3,
			Name: "Salad",
		},
	}

	return &recipepb.GetResponse{
		Recipes: recipes,
	}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Attach the Recipes service to the server
	recipepb.RegisterRecipesServer(s, &server{})

	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	// Register Recipes API
	err = recipepb.RegisterRecipesHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":3000",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:3000")
	log.Fatalln(gwServer.ListenAndServe())
}
