package main

import (
	api "github.com/KhetwalDevesh/book-my-seat/server/internal/apis"
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

var address = "0.0.0.0:50051"

func main() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	log.Printf("Listening on %s\n", address)
	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, api.NewBookingServiceServer())
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve : %v\n", err)
	}
}
