API_URL=http://0.0.0.0:3000
HEADERS=Content-Type:application/json

install:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/cosmtrek/air@latest

generate:
	protoc -I ./proto \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
		./proto/chat/rooms.proto

build:
	go mod tidy && \
	go build -o api main.go

start:
	./api

dev: install generate build
	air

# Commands to interact with the API
get-rooms:
	http ${API_URL}/rooms

create-room:
	echo '{"name": "test"}' | http POST ${API_URL}/rooms ${HEADERS}
