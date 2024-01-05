package apis

import pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"

type BookingServiceServer struct {
	pb.BookingServiceServer
	Tickets     map[string]pb.Ticket            // emailId is the key here
	SeatMapping map[string]map[string]pb.Ticket // seat_section is the key to outer map, emailId is the key to inner map
	SeatCounter map[string]uint32               // seat_section is the key and seat_counts is the value here
}
