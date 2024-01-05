package main

import (
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

var address = "0.0.0.0:50051"

type BookingServiceServer struct {
	pb.BookingServiceServer
	tickets     map[string]pb.Ticket          // emailId is the key here
	seatMapping map[string]map[string]pb.User // seat_section is the key to outer map, emailId is the key to inner map
	seatCounter map[string]uint32             // seat_section is the key and seat_counts is the value here
}

func main() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	log.Printf("Listening on %s\n", address)
	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, &BookingServiceServer{
		tickets:     make(map[string]pb.Ticket),
		seatMapping: map[string]map[string]pb.User{pb.SeatSection_A.String(): {}, pb.SeatSection_B.String(): {}},
		seatCounter: make(map[string]uint32),
	})
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve : %v\n", err)
	}
}
