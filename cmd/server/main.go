package main

import (
	pb "grpc-project/booking/proto"
	"grpc-project/cmd/server/models"
	"grpc-project/cmd/server/service"
	"log"
	"net"

	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var Store = &models.Store{
	Train: models.Train{
		Id: uuid.New().String(),
	},
}

func init() {
	//Initialize Store Data structure
	//Assume Train has 2 sections with  20 seats each

	sectionCount := 2
	seatCount := 20
	price := 20

	for i := 0; i < sectionCount; i++ {
		section := &models.Section{
			Id:             uuid.New().String(),
			Name:           "Section " + fmt.Sprint(i+1),
			Seats:          make([]*models.Seat, seatCount),
			AvailableSeats: seatCount,
		}

		for j := 0; j < seatCount; j++ {
			seat := &models.Seat{
				Id:            uuid.New().String(),
				SectionName:   section.Name,
				SectionId:     section.Id,
				SeatNumber:    "Seat " + fmt.Sprint(j+1),
				SeatAvailable: true,
				User:          nil,
			}
			section.Seats[j] = seat
		}
		Store.Train.Sections = append(Store.Train.Sections, section)

	}

	//Create User
	alice := &models.User{
		Id:        uuid.New().String(),
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alicewonderland@gmal.com",
	}
	bob := &models.User{
		Id:        uuid.New().String(),
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     "bobthebuilder@gmail.com",
	}
	Store.Users = append(Store.Users, alice, bob)
	Store.Train.Price = float32(price)
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

	// ctx := context.Background()
	// purchasebookingRequest := &pb.PurchaseBookingRequest{
	// 	From: "London",
	// 	To:   "France",
	// 	User: &pb.User{
	// 		UserId:    bookingService.Store.Users[0].Id,
	// 		FirstName: "Alice",
	// 		LastName:  "Smith",
	// 		Email:     "AliceSmith@gmaiil.com",
	// 	},
	// 	PricePaid: 20.0,
	// }
	// res, _ := bookingService.PurchaseBooking(ctx, purchasebookingRequest)
	// fmt.Println("Response:", res.Receipt.Seat)
	// fmt.Println("Response  Section:", res.Receipt.Section)

	log.Println("Server is running on port %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
