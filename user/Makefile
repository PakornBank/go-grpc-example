.PHONY: setup build test run proto generate

setup:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate:
	go generate ./...

build:
	go build -o bin/server ./cmd/grpc

test:
	go test ./...

docker-up:
	docker-compose --env-file .env -f docker-compose.yml up -d

run:
	go run ./cmd/grpc

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/v1/user.proto

