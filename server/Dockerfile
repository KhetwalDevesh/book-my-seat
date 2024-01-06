# Use an official Golang runtime as a parent image
FROM golang:1.21.3-bookworm

# Set the working directory to /go/src/app
WORKDIR /go/src/app

RUN go mod init github.com/KhetwalDevesh/book-my-seat

# Install the Protocol Buffers compiler (protoc)
RUN apt-get update &&  apt install -y protobuf-compiler

# Install the protoc-gen-go plugin
RUN apt-get update && apt-get install -y protoc-gen-go
RUN apt-get update && apt-get install -y protoc-gen-go-grpc

RUN go get google.golang.org/grpc
RUN go get google.golang.org/grpc/codes
RUN go get google.golang.org/grpc/status

ENV GOBIN /go/bin

ENV PATH="$PATH:/go/bin"

# Copy the proto file from the root directory into the container
COPY proto/ ./proto/

# Create an empty stubs directory
RUN mkdir -p ./stubs

# Copy your server code
COPY server/ ./server/

# Generate Go code from the proto file
RUN protoc -I ./proto \
    --go_out=paths=import:./stubs \
    --go_opt=paths=source_relative \
    --go-grpc_out=paths=import:./stubs \
    --go-grpc_opt=paths=source_relative \
    ./proto/booking-service/v1/booking.proto

# Run go mod tidy to update go.mod and go.sum
RUN go mod tidy

# Build your Go server application
RUN go build -o bin/server ./server 2>&1

# Expose port 50051 for gRPC
EXPOSE 50051

# Command to run the server
CMD ["./bin/server"]