package apis_test

import (
	"context"
	api "github.com/KhetwalDevesh/book-my-seat/server/internal/apis"
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"testing"
)

func TestPurchaseTicket(t *testing.T) {
	server := api.NewBookingServiceServer()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("TestPurchaseTicket panicked: %v", r)
		}
	}()

	ctx := context.Background()
	user := &pb.User{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	request := &pb.PurchaseTicketRequest{
		User:        user,
		SeatSection: pb.SeatSection_A,
		SeatNumber:  1,
		TicketPrice: 50.0,
	}

	// Call the function being tested
	response, err := server.PurchaseTicket(ctx, request)
	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	// Assert the expected result
	expectedTicket := &pb.Ticket{
		From:        "London",
		To:          "France",
		User:        user,
		PricePaid:   50.0,
		SeatSection: pb.SeatSection_A,
		SeatNumber:  1,
	}

	if response.Ticket.From != expectedTicket.From ||
		response.Ticket.To != expectedTicket.To ||
		response.Ticket.User.Email != expectedTicket.User.Email ||
		response.Ticket.PricePaid != expectedTicket.PricePaid ||
		response.Ticket.SeatSection != expectedTicket.SeatSection ||
		response.Ticket.SeatNumber != expectedTicket.SeatNumber {
		t.Fatalf("Unexpected response. Expected %v, got %v", expectedTicket, response.Ticket)
	}
}

func TestGetReceipt(t *testing.T) {
	server := api.NewBookingServiceServer()
	email := "john.doe@example.com"
	expectedTicket := &pb.Ticket{
		From:        "London",
		To:          "France",
		User:        &pb.User{Id: 1, FirstName: "John", LastName: "Doe", Email: email},
		PricePaid:   50.0,
		SeatSection: pb.SeatSection_A,
		SeatNumber:  1,
	}

	server.Tickets[email] = *expectedTicket

	ctx := context.Background()
	request := &pb.GetReceiptRequest{Email: email}

	// Call the function being tested
	response, err := server.GetReceipt(ctx, request)
	if err != nil {
		t.Fatalf("GetReceipt failed: %v", err)
	}

	// Assert the expected result
	if response.Ticket.From != expectedTicket.From ||
		response.Ticket.To != expectedTicket.To ||
		response.Ticket.User.Email != expectedTicket.User.Email ||
		response.Ticket.PricePaid != expectedTicket.PricePaid ||
		response.Ticket.SeatSection != expectedTicket.SeatSection ||
		response.Ticket.SeatNumber != expectedTicket.SeatNumber {
		t.Fatalf("Unexpected response. Expected %v, got %v", expectedTicket, response.Ticket)
	}
}

func TestGetUsersAndSeatAllocated(t *testing.T) {
	server := api.NewBookingServiceServer()
	userEmail := "john.doe@example.com"
	expectedTicket := &pb.Ticket{
		From:        "London",
		To:          "France",
		User:        &pb.User{Id: 1, FirstName: "John", LastName: "Doe", Email: userEmail},
		PricePaid:   50.0,
		SeatSection: pb.SeatSection_A,
		SeatNumber:  1,
	}

	server.SeatMapping[pb.SeatSection_A.String()] = map[string]pb.Ticket{userEmail: *expectedTicket}

	ctx := context.Background()
	request := &pb.GetUsersAndSeatAllocatedRequest{SeatSection: pb.SeatSection_A}

	// Call the function being tested
	response, err := server.GetUsersAndSeatAllocated(ctx, request)
	if err != nil {
		t.Fatalf("GetUsersAndSeatAllocated failed: %v", err)
	}

	// Assert the expected result
	if len(response.SeatAllocated) != 1 {
		t.Fatalf("Unexpected response. Expected 1 user, got %d", len(response.SeatAllocated))
	}

	actualTicket := response.SeatAllocated[userEmail]
	if actualTicket.From != expectedTicket.From ||
		actualTicket.To != expectedTicket.To ||
		actualTicket.User.Email != expectedTicket.User.Email ||
		actualTicket.PricePaid != expectedTicket.PricePaid ||
		actualTicket.SeatSection != expectedTicket.SeatSection ||
		actualTicket.SeatNumber != expectedTicket.SeatNumber {
		t.Fatalf("Unexpected response. Expected %v, got %v", expectedTicket, actualTicket)
	}
}

func TestRemoveUser(t *testing.T) {
	server := api.NewBookingServiceServer()
	userEmail := "john.doe@example.com"
	expectedTicket := &pb.Ticket{
		From:        "London",
		To:          "France",
		User:        &pb.User{Id: 1, FirstName: "John", LastName: "Doe", Email: userEmail},
		PricePaid:   50.0,
		SeatSection: pb.SeatSection_A,
		SeatNumber:  1,
	}

	server.Tickets[userEmail] = *expectedTicket
	server.SeatMapping[pb.SeatSection_A.String()] = map[string]pb.Ticket{userEmail: *expectedTicket}

	ctx := context.Background()
	request := &pb.RemoveUserRequest{Email: userEmail}

	// Call the function being tested
	response, err := server.RemoveUser(ctx, request)
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	// Assert the expected result
	if _, exists := server.Tickets[userEmail]; exists {
		t.Fatalf("User ticket not removed")
	}

	if _, exists := server.SeatMapping[pb.SeatSection_A.String()][userEmail]; exists {
		t.Fatalf("User seat allocation not removed")
	}

	expectedMessage := "User removed successfully"
	if response.Msg != expectedMessage {
		t.Fatalf("Unexpected response. Expected '%s', got '%s'", expectedMessage, response.Msg)
	}
}

func TestModifyUserSeat(t *testing.T) {
	server := api.NewBookingServiceServer()

	userEmail := "john.doe@example.com"
	expectedTicket := &pb.Ticket{
		From:        "London",
		To:          "France",
		User:        &pb.User{Id: 1, FirstName: "John", LastName: "Doe", Email: userEmail},
		PricePaid:   50.0,
		SeatSection: pb.SeatSection_A,
		SeatNumber:  1,
	}

	server.Tickets[userEmail] = *expectedTicket
	server.SeatMapping[pb.SeatSection_A.String()] = map[string]pb.Ticket{userEmail: *expectedTicket}

	ctx := context.Background()
	request := &pb.ModifyUserSeatRequest{
		Email:          userEmail,
		NewSeatSection: pb.SeatSection_B,
		NewSeatNumber:  2,
	}

	// Call the function being tested
	response, err := server.ModifyUserSeat(ctx, request)
	if err != nil {
		t.Fatalf("ModifyUserSeat failed: %v", err)
	}

	// Assert the expected result
	if ticket, exists := server.Tickets[userEmail]; !exists || ticket.SeatSection != pb.SeatSection_B || ticket.SeatNumber != 2 {
		t.Fatalf("User ticket not modified")
	}

	if ticket, exists := server.SeatMapping[pb.SeatSection_B.String()][userEmail]; !exists || ticket.SeatSection != pb.SeatSection_B || ticket.SeatNumber != 2 {
		t.Fatalf("User seat allocation not modified")
	}

	expectedMessage := "User Seat successfully modified"
	if response.Msg != expectedMessage {
		t.Fatalf("Unexpected response. Expected '%s', got '%s'", expectedMessage, response.Msg)
	}
}
