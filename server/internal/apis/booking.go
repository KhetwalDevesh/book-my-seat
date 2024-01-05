package apis

import (
	"context"
	"fmt"
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"github.com/google/uuid"
	"log"
)

func (s *BookingServiceServer) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.PurchaseTicketResponse, error) {
	seatSection := req.SeatSection
	seatNumber := req.SeatNumber

	// Check if user already purchased a ticket
	_, exists := s.Tickets[req.User.Email]
	if exists {
		return nil, fmt.Errorf("User already have a ticket")
	}

	// Check if the section matches the available sections in the train
	if seatSection != pb.SeatSection_A && seatSection != pb.SeatSection_B {
		return nil, fmt.Errorf("invalid seat section")
	}

	// Check if there are available seats, 50 taken for each seatSection, just to put some limit
	if s.SeatCounter[seatSection.String()] >= 50 {
		return nil, fmt.Errorf("all seats in seatSection %s are occupied", seatSection)
	}

	seatAlreadyOccupied := false
	for _, ticket := range s.SeatMapping[seatSection.String()] {
		if ticket.SeatNumber == seatNumber {
			seatAlreadyOccupied = true
		}
	}
	if seatAlreadyOccupied {
		log.Fatalf("Seat already occupied, choose some other")
		return nil, fmt.Errorf("Seat already occupied, choose some other")
	}

	// Increment the seat counter and assign the seat number
	s.SeatCounter[seatSection.String()]++

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
		SeatSection: seatSection,
		SeatNumber:  seatNumber,
	}

	// Store the ticket and seat allocation
	s.Tickets[req.User.Email] = *ticket
	s.SeatMapping[seatSection.String()][req.User.Email] = *ticket

	return &pb.PurchaseTicketResponse{Ticket: ticket}, nil
}

func (s *BookingServiceServer) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	email := req.Email
	ticket := s.Tickets[email]
	return &pb.GetReceiptResponse{Ticket: &pb.Ticket{
		From:        ticket.From,
		To:          ticket.To,
		User:        ticket.User,
		PricePaid:   ticket.PricePaid,
		SeatSection: ticket.SeatSection,
		SeatNumber:  ticket.SeatNumber,
	}}, nil
}

func (s *BookingServiceServer) GetUsersAndSeatAllocated(ctx context.Context, req *pb.GetUsersAndSeatAllocatedRequest) (*pb.GetUsersAndSeatAllocatedResponse, error) {
	section := req.SeatSection.String()
	usersAndSeatAllocated := s.SeatMapping[section]

	pbUsersAndSeatAllocated := make(map[string]*pb.Ticket)
	for email, ticket := range usersAndSeatAllocated {
		pbUsersAndSeatAllocated[email] = &pb.Ticket{
			From:        ticket.From,
			To:          ticket.To,
			User:        ticket.User,
			PricePaid:   ticket.PricePaid,
			SeatSection: ticket.SeatSection,
			SeatNumber:  ticket.SeatNumber,
		}
	}
	return &pb.GetUsersAndSeatAllocatedResponse{SeatAllocated: pbUsersAndSeatAllocated}, nil
}

func (s *BookingServiceServer) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	_, exists := s.Tickets[req.Email]
	if !exists {
		return nil, fmt.Errorf("User not found")
	}

	// Remove user and seat allocation
	delete(s.Tickets, req.Email)
	for _, section := range []string{"A", "B"} {
		delete(s.SeatMapping[section], req.Email)
	}
	return &pb.RemoveUserResponse{Msg: "User removed successfully"}, nil
}

func (s *BookingServiceServer) ModifyUserSeat(ctx context.Context, req *pb.ModifyUserSeatRequest) (*pb.ModifyUserSeatResponse, error) {
	newSeatSection := req.NewSeatSection
	newSeatNumber := req.NewSeatNumber
	userEmail := req.Email

	_, exists := s.Tickets[userEmail]
	if !exists {
		return nil, fmt.Errorf("User not found!")
	}

	// Check if the section matches the available sections in the train
	if newSeatSection != pb.SeatSection_A && newSeatSection != pb.SeatSection_B {
		return nil, fmt.Errorf("invalid seat section")
	}

	// Check if there are available seats, 50 taken for each section, just to put some limit
	if newSeatNumber > 50 {
		log.Fatalf("Invalid seat number, only 50 seats exists")
		return nil, nil
	}

	seatAlreadyOccupied := false
	for _, ticket := range s.SeatMapping[newSeatSection.String()] {
		if ticket.SeatNumber == newSeatNumber {
			seatAlreadyOccupied = true
		}
	}
	if seatAlreadyOccupied {
		log.Fatalf("Seat already occupied, choose some other")
		return nil, nil
	}

	ticket, _ := s.Tickets[userEmail]
	// delete the old instance of seat allocated
	delete(s.SeatMapping[ticket.SeatSection.String()], userEmail)
	ticket.SeatSection = newSeatSection
	ticket.SeatNumber = newSeatNumber
	s.Tickets[userEmail] = ticket
	s.SeatMapping[newSeatSection.String()][userEmail] = ticket

	return &pb.ModifyUserSeatResponse{Msg: "User Seat successfully modified"}, nil
}
