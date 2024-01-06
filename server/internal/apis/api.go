package apis

import pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"

type BookingServiceServer struct {
	pb.BookingServiceServer
	Tickets     map[string]pb.Ticket            // emailId is the key here
	SeatMapping map[string]map[string]pb.Ticket // seat_section is the key to outer map, emailId is the key to inner map
}

// NewBookingServiceServer creates a new instance of BookingServiceServer with initialized maps.
func NewBookingServiceServer() *BookingServiceServer {
	return &BookingServiceServer{
		Tickets:     make(map[string]pb.Ticket),
		SeatMapping: make(map[string]map[string]pb.Ticket),
	}
}
