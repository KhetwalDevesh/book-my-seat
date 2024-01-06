package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"strings"
)

var serverAddress = "book-my-seat:50051"

func PurchaseTicket(client pb.BookingServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	// Get input for User fields
	fmt.Println("Enter User Details : ")
	fmt.Print("Enter First Name : ")
	firstName, _ := reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)

	fmt.Print("Enter Last Name : ")
	lastName, _ := reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)

	fmt.Print("Enter Email : ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	user := &pb.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	// Get input for SeatSection field
	fmt.Print("Choose Section ( A or B ) : ")
	seatSectionStr, _ := reader.ReadString('\n')
	seatSectionStr = strings.TrimSpace(seatSectionStr)
	seatSection := pb.SeatSection(pb.SeatSection_value[strings.ToUpper(seatSectionStr)])
	if seatSection < 0 || seatSection > 1 {
		log.Fatalf("Invalid seat section : %s", seatSectionStr)
	}

	// Get input for SeatNumber field
	fmt.Print("Choose Seat Number from 1 to 50 : ")
	seatNumberStr, _ := reader.ReadString('\n')
	seatNumberStr = strings.TrimSpace(seatNumberStr)
	seatNumber, err := strconv.ParseUint(seatNumberStr, 10, 32)
	if err != nil {
		log.Fatalf("Invalid Seat Number : %v", err)
	}

	// Make the user to enter the price
	fmt.Print("Enter the price of the ticket ( $20 ) : ")
	ticketPriceStr, _ := reader.ReadString('\n')
	ticketPriceStr = strings.TrimSpace(ticketPriceStr)
	ticketPrice, err := strconv.ParseFloat(ticketPriceStr, 32)
	if ticketPrice != 20 {
		log.Fatalf("Please enter the correct ticket price i.e $20")
	}

	// Finally call the grpc method PurchaseTicket
	response, err := client.PurchaseTicket(context.Background(), &pb.PurchaseTicketRequest{
		User:        user,
		SeatSection: seatSection,
		SeatNumber:  uint32(seatNumber),
	})
	if err != nil {
		log.Fatalf("Error calling PurchaseTicket: %v", err.Error())
	}
	fmt.Printf("\nBooking Details:\n\n%+v\n\n", response)
}

func GetReceipt(client pb.BookingServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Email to get receipt : ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// Call the grpc method GetReceipt
	response, err := client.GetReceipt(context.Background(), &pb.GetReceiptRequest{Email: email})
	if err != nil {
		log.Fatalf("Error calling GetReceipt : %v", err)
	}
	fmt.Printf("\nReceipt Details:\n\n%+v\n\n", response)
}

func GetUsersAndSeatAllocated(client pb.BookingServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose Section ( A or B ) : ")
	seatSectionStr, _ := reader.ReadString('\n')
	seatSectionStr = strings.TrimSpace(seatSectionStr)
	seatSection := pb.SeatSection(pb.SeatSection_value[strings.ToUpper(seatSectionStr)])
	// Call the grpc method GetUsersAndSeatAllocated
	response, err := client.GetUsersAndSeatAllocated(context.Background(), &pb.GetUsersAndSeatAllocatedRequest{SeatSection: seatSection})
	if err != nil {
		log.Fatalf("Error calling GetUsersAndSeatAllocated : %v", err)
	}
	fmt.Printf("\nUser and Seat allocated Details : \n")
	for _, ticket := range response.SeatAllocated {
		fmt.Printf("%+v\n", ticket)
	}
}

func RemoveUser(client pb.BookingServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Email of the User to be removed : ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// Call the grpc method RemoveUser
	response, err := client.RemoveUser(context.Background(), &pb.RemoveUserRequest{Email: email})
	if err != nil {
		log.Fatalf("Error calling RemoveUser : %v", err)
	}
	fmt.Printf("\nMessage :\n\n%+v\n\n", response)
}

func ModifyUserSeat(client pb.BookingServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Email of the User whose seat is to be updated : ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	fmt.Print("Enter the new SeatSection ( A or B ) : ")
	seatSectionStr, _ := reader.ReadString('\n')
	seatSectionStr = strings.TrimSpace(seatSectionStr)
	seatSection := pb.SeatSection(pb.SeatSection_value[strings.ToUpper(seatSectionStr)])
	fmt.Print("Enter the new SeatNumber (anywhere from 1 to 50 ) : ")
	seatNumberStr, _ := reader.ReadString('\n')
	seatNumberStr = strings.TrimSpace(seatNumberStr)
	seatNumber, err := strconv.ParseUint(seatNumberStr, 10, 32)
	if err != nil {
		log.Fatalf("Invalid seat number")
	}

	// call the grpc method ModifyUserSeat
	response, err := client.ModifyUserSeat(context.Background(), &pb.ModifyUserSeatRequest{
		Email:          email,
		NewSeatSection: seatSection,
		NewSeatNumber:  uint32(seatNumber),
	})
	if err != nil {
		log.Fatalf("Error calling ModifyUserRequest : %v", err)
	}
	fmt.Printf("\nMessage : \n\n%+v\n\n", response)
}

func main() {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %v\n", err)
	}
	defer conn.Close()
	client := pb.NewBookingServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose an option :\n 1 to purchase a Ticket\n 2 to view the details of the receipt of a User\n 3 to view the Users and seat they are allocated" +
		"\n 4 to remove a user from the train\n 5 to modify a user's seat")
	inputOption, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error inputting the option : %v\n", err)
	}
	inputOption = strings.TrimSpace(inputOption)
	// below method converts the inputOption which is a string to integer
	option, err := strconv.Atoi(inputOption)
	if err != nil {
		log.Fatalf("Invalid option: %v", err)
	}
	switch option {
	case 1:
		PurchaseTicket(client)
	case 2:
		GetReceipt(client)
	case 3:
		GetUsersAndSeatAllocated(client)
	case 4:
		RemoveUser(client)
	case 5:
		ModifyUserSeat(client)
	default:
		log.Fatal("Please choose a valid option.")
	}
}
