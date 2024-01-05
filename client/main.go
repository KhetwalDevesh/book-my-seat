package main

import (
	"bufio"
	"fmt"
	pb "github.com/KhetwalDevesh/book-my-seat/stubs/booking-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"strings"
)

var serverAddress = "0.0.0.0:50051"

func PurchaseTicket(client pb.BookingServiceClient)           {}
func GetReceipt(client pb.BookingServiceClient)               {}
func GetUsersAndSeatAllocated(client pb.BookingServiceClient) {}
func RemoveUser(client pb.BookingServiceClient)               {}
func ModifyUserSeat(client pb.BookingServiceClient)           {}

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
