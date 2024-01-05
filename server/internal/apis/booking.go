package apis

import (
	"context"
	"fmt"
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"github.com/google/uuid"
)

func (s *BookingServiceServer) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.PurchaseTicketResponse, error) {
	section := req.SeatSection
	// Check if the section matches the available sections in the train
	if section != pb.SeatSection_A || section != pb.SeatSection_B {
		return nil, fmt.Errorf("invalid seat section")
	}

	// Check if there are available seats, 50 taken for each section, just to put some limit
	if s.seatCounter[section.String()] >= 50 {
		return nil, fmt.Errorf("all seats in section %s are occupied", section)
	}

	// Increment the seat counter and assign the seat number
	seatNumber := s.seatCounter[section.String()] + 1
	s.seatCounter[section.String()]++

	// Generate a unique ID for the user
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user ID: %v", err)
	}

	// Create a new ticket with the unique ID
	ticket := &pb.Ticket{
		From: "London",
		To:   "France",
		User: &pb.User{
			Id:        uint64(userID.ID()),
			FirstName: req.User.FirstName,
			LastName:  req.User.LastName,
			Email:     req.User.Email,
		},
		PricePaid:   20.0,
		SeatSection: section,
		SeatNumber:  seatNumber,
	}

	// Store the ticket and seat allocation
	s.tickets[req.User.Email] = *ticket
	s.seatMapping[section.String()][req.User.Email] = *ticket.User

	return &pb.PurchaseTicketResponse{Ticket: ticket}, nil
}

func (s *BookingServiceServer) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	// Implement the logic for the GetReceipt RPC
	return nil, fmt.Errorf("not implemented")
}

func (s *BookingServiceServer) GetSeatAllocated(ctx context.Context, req *pb.GetSeatAllocatedRequest) (*pb.GetSeatAllocatedResponse, error) {
	// Implement the logic for the GetSeatAllocated RPC
	return nil, fmt.Errorf("not implemented")
}

func (s *BookingServiceServer) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	// Implement the logic for the RemoveUser RPC
	return nil, fmt.Errorf("not implemented")
}

func (s *BookingServiceServer) ModifyUserSeat(ctx context.Context, req *pb.ModifyUserSeatRequest) (*pb.ModifyUserSeatResponse, error) {
	// Implement the logic for the ModifyUserSeat RPC
	return nil, fmt.Errorf("not implemented")
}
