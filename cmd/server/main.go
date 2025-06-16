package main

import (
	"fmt"
	pb "grpc-project/booking/proto"
	"grpc-project/cmd/server/models"
	"grpc-project/cmd/server/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var Store = &models.Store{}

func init() {
	//Initialize Store Data structure
	//Assume Train has 2 sections with  20 seats each

	sectionCount := 2
	seatCount := 20
	// price := 20
	from := "London"
	to := "France"

	Store, err := service.GenerateNewTrain(Store, sectionCount, seatCount, from, to, 20.0)
	if err != nil {
		fmt.Errorf("Error creating new train %w ", err)
	}

	//Create User GenerateNewTrain(
	alice := &models.User{
		Id:        "1",
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alicewonderland@gmal.com",
		Receipts:  []*models.Receipt{},
	}
	bob := &models.User{
		Id:        "2",
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     "bobthebuilder@gmail.com",
		Receipts:  []*models.Receipt{},
	}
	Store.Users = append(Store.Users, alice, bob)
	Store.Receipts = make(map[string]models.Receipt)
}

func main() {

	//Listen on port 8080
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//create a new gRPC server
	s := grpc.NewServer()

	//Register the booking service with the server
	bookingService := &service.BookingServer{
		Store: Store,
	}
	pb.RegisterBookingServiceServer(s, bookingService)

	reflection.Register(s)

	log.Println("Server is running on port %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
