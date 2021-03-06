package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	chatpb "github.com/andreafalzetti/grpc-go-example/proto/chat"
)

var chatRooms = []*chatpb.ChatRoom{
	{
		Id:   0,
		Name: "Coding",
	},
	{
		Id:   1,
		Name: "Travel",
	},
	{
		Id:   2,
		Name: "Investing",
	},
	{
		Id:   3,
		Name: "Gaming",
	},
}

type server struct {
	chatpb.UnimplementedChatRoomsServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Get(ctx context.Context, req *chatpb.GetRequest) (*chatpb.GetResponse, error) {
	return &chatpb.GetResponse{
		Rooms: chatRooms,
	}, nil
}

func (s *server) Create(ctx context.Context, req *chatpb.CreateRequest) (*chatpb.CreateResponse, error) {
	fmt.Println("new room created: ", req)

	newRoomId := int32(len(chatRooms))
	newRoomName := req.GetName()

	chatRooms = append(chatRooms, &chatpb.ChatRoom{
		Id:   newRoomId,
		Name: newRoomName,
	})

	return &chatpb.CreateResponse{
		Id:   newRoomId,
		Name: newRoomName,
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

	// Attach the Chat service to the server
	chatpb.RegisterChatRoomsServer(s, &server{})

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

	// Register Chat API
	err = chatpb.RegisterChatRoomsHandler(context.Background(), gwmux, conn)
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
