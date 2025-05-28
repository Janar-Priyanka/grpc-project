package service

import (
	"context"
	"fmt"
	pb "grpc-project/booking/proto"
	"grpc-project/cmd/server/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func InitializeStore() *models.Store {
	// Initialize the store with some data
	store := &models.Store{
		Train: models.Train{
			Id:    "123-4567-8901-2345",
			From:  "London",
			To:    "France",
			Price: 20.0,
			Sections: []*models.Section{
				{
					Id:             "S1",
					Name:           "Section 1",
					AvailableSeats: 20,
					Seats: []*models.Seat{
						{Id: uuid.New().String(), SectionName: "Section 1", SectionId: "S1", SeatNumber: "1", SeatAvailable: false,
							User: &models.User{
								Id:        "1",
								FirstName: "Alice",
								LastName:  "Smith",
								Email:     "AliceSmith@gmaiil.com",
							},
						},
						{Id: uuid.New().String(), SectionName: "Section 1", SectionId: "S1", SeatNumber: "2", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 1", SectionId: "S1", SeatNumber: "3", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 1", SectionId: "S1", SeatNumber: "4", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 1", SectionId: "S1", SeatNumber: "5", SeatAvailable: true},
					},
				},
				{
					Id:             uuid.New().String(),
					Name:           "Section 2",
					AvailableSeats: 20,
					Seats: []*models.Seat{
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "1", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "2", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "3", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "4", SeatAvailable: true},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "5", SeatAvailable: true},
					},
				},
			},
		},
		Users: []*models.User{
			{
				Id:        "1",
				FirstName: "Alice",
				LastName:  "Smith",
				Email:     "AliceSmith@gmaiil.com",
			},
			{
				Id:        "2",
				FirstName: "Bob",
				LastName:  "Johnson",
				Email:     "BobJohnson@gmail.com",
			},
		},
		Receipts: make(map[string]models.Receipt),
	}
	aliceReceipts := []*models.Receipt{
		{
			Id:            "11",
			From:          "London",
			To:            "France",
			Email:         "AliceSmith@gmaiil.com",
			SeatNumber:    "1",
			SectionName:   "Section 1",
			SectionId:     store.Train.Sections[0].Id,
			SeatId:        store.Train.Sections[0].Seats[0].Id,
			UserId:        "1",
			BookingStatus: "Confirmed",
		},
	}
	// Assign the receipts to the user and the store
	store.Users[0].Receipts = aliceReceipts
	store.Train.Sections[0].Seats[0].User.Receipts = aliceReceipts
	store.Receipts[aliceReceipts[0].Id] = *aliceReceipts[0]
	return store
}

func Test_PurchaseBooking(t *testing.T) {
	store := InitializeStore()
	type test struct {
		PurchaseRequest  *pb.PurchaseBookingRequest
		ExpectedResponse *pb.PurchaseBookingResponse
		ExpectedError    error
		Store            *models.Store
		ExpectedStore    *models.Store
	}
	tests := map[string]test{
		"Happy Path - Valid Booking": {
			PurchaseRequest: &pb.PurchaseBookingRequest{
				From: "London",
				To:   "France",
				User: &pb.User{
					UserId:    store.Users[0].Id,
					FirstName: "Alice",
					LastName:  "Smith",
					Email:     "AliceSmith@gmaiil.com",
				},
				PricePaid: 20.0,
			},
			ExpectedResponse: &pb.PurchaseBookingResponse{
				Receipt: &pb.Receipt{
					From: "London",
					To:   "France",
					User: &pb.User{
						FirstName: "Alice",
						LastName:  "Smith",
						Email:     "AliceSmith@gmaiil.com",
					},
					BookingStatus: "Confirmed",
					Seat:          "1",
					Section:       "Section 1",
					PricePaid:     20.0,
				},
			},
		},
		"Sad Path - Invalid Booking Request when request does not have From details": {
			PurchaseRequest: &pb.PurchaseBookingRequest{
				From: "",
				To:   "France",
				User: &pb.User{
					UserId:    store.Users[0].Id,
					FirstName: "Alice",
					LastName:  "Smith",
					Email:     "AliceSmith@gmaiil.com",
				},
				PricePaid: 20.0,
			},
			ExpectedError:    fmt.Errorf("Invalid Booking Request"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Booking Request when request does not have To details": {
			PurchaseRequest: &pb.PurchaseBookingRequest{
				From: "London",
				To:   "",
				User: &pb.User{
					UserId:    store.Users[0].Id,
					FirstName: "Alice",
					LastName:  "Smith",
					Email:     "AliceSmith@gmaiil.com",
				},
				PricePaid: 20.0,
			},
			ExpectedError:    fmt.Errorf("Invalid Booking Request"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Booking Request when request does not have User details": {
			PurchaseRequest: &pb.PurchaseBookingRequest{
				From:      "London",
				To:        "",
				User:      nil,
				PricePaid: 20.0,
			},
			Store:            InitializeStore(),
			ExpectedError:    fmt.Errorf("Invalid Booking Request"),
			ExpectedResponse: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			bookingServer := &BookingServer{
				Store: store,
			}

			res, err := bookingServer.PurchaseBooking(ctx, tc.PurchaseRequest)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected Error but didnt got error failed: %v", err)
				}
				if fmt.Sprint(err) != fmt.Sprint(tc.ExpectedError) {
					t.Errorf("Expected Error: %v, but got: %v", tc.ExpectedError, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Expected no errirr but we got an error : %v", err)

				}
			}
			//assert.NotNil(t, res)
			assert.Equal(t, tc.ExpectedResponse.Receipt.From, res.Receipt.From)
			assert.Equal(t, tc.ExpectedResponse.Receipt.To, res.Receipt.To)
			assert.Equal(t, tc.ExpectedResponse.Receipt.User.FirstName, res.Receipt.User.FirstName)
			assert.Equal(t, tc.ExpectedResponse.Receipt.User.LastName, res.Receipt.User.LastName)
			assert.Equal(t, tc.ExpectedResponse.Receipt.User.Email, res.Receipt.User.Email)
			assert.Equal(t, tc.ExpectedResponse.Receipt.BookingStatus, res.Receipt.BookingStatus)
			assert.NotEmpty(t, res.Receipt.Seat, "Seat should not be empty")
			assert.NotEmpty(t, res.Receipt.Section, "Section should not be empty")
			assert.NotEmpty(t, res.Receipt.PricePaid, "PricePaid should not be empty")

		})
	}

}

func Test_PurchaseBooking_WhenNoAvailableSeatsOrSections(t *testing.T) {
	store := &models.Store{
		Train: models.Train{
			Id:    uuid.New().String(),
			From:  "London",
			To:    "France",
			Price: 20.0,
			Sections: []*models.Section{
				{
					Id:             uuid.New().String(),
					Name:           "Section 1",
					AvailableSeats: 0,
					Seats: []*models.Seat{
						{Id: uuid.New().String(), SectionName: "Section 1", SeatNumber: "1", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 1", SeatNumber: "2", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 1", SeatNumber: "3", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 1", SeatNumber: "4", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 1", SeatNumber: "5", SeatAvailable: false},
					},
				},
				{
					Id:             uuid.New().String(),
					Name:           "Section 2",
					AvailableSeats: 0,
					Seats: []*models.Seat{
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "1", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "2", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "3", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "4", SeatAvailable: false},
						{Id: uuid.New().String(), SectionName: "Section 2", SeatNumber: "5", SeatAvailable: false},
					},
				},
			},
		},
		Users: []*models.User{
			{
				Id:        uuid.New().String(),
				FirstName: "Alice",
				LastName:  "Smith",
				Email:     "AliceSmith@gmaiil.com",
			},
			{
				Id:        uuid.New().String(),
				FirstName: "Bob",
				LastName:  "Johnson",
				Email:     "BobJohnson@gmail.com",
			},
		},
		Receipts: make(map[string]models.Receipt),
	}
	type test struct {
		PurchaseRequest  *pb.PurchaseBookingRequest
		ExpectedError    error
		ExpectedResponse *pb.PurchaseBookingResponse
	}
	tests := map[string]test{
		"Send No Sections are available when all the seats are full  ": {
			PurchaseRequest: &pb.PurchaseBookingRequest{
				From: "London",
				To:   "France",
				User: &pb.User{
					UserId:    store.Users[0].Id,
					FirstName: "Alice",
					LastName:  "Smith",
					Email:     "AliceSmith@gmaiil.com",
				},
				PricePaid: 20.0,
			},
			ExpectedError: fmt.Errorf("No available seats found"),
		},
		"Send No Seats are Available when the train is full": {
			PurchaseRequest: &pb.PurchaseBookingRequest{
				From: "London",
				To:   "France",
				User: &pb.User{
					UserId:    store.Users[0].Id,
					FirstName: "Alice",
					LastName:  "Smith",
					Email:     "AliceSmith@gmaiil.com",
				},
				PricePaid: 20.0,
			},
			ExpectedError: fmt.Errorf("No available seats found"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			bookingServer := &BookingServer{
				Store: store,
			}

			_, err := bookingServer.PurchaseBooking(ctx, tc.PurchaseRequest)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected Error but didnt got error failed: %v", err)
				}
				if fmt.Sprint(err) != fmt.Sprint(tc.ExpectedError) {
					t.Errorf("Expected Error: %v, but got: %v", tc.ExpectedError, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Expected no errirr but we got an error : %v", err)

				}
			}
		})
	}
}

func Test_ShowReceipt(t *testing.T) {
	store := InitializeStore()
	type test struct {
		ShowReceiptRequest *pb.ShowReceiptRequest
		ExpectedResponse   *pb.ShowReceiptResponse
		ExpectedError      error
	}
	tests := map[string]test{
		"Happy Path - Valid Receipt": {
			ShowReceiptRequest: &pb.ShowReceiptRequest{
				UserId: "1",
			},
			ExpectedResponse: &pb.ShowReceiptResponse{
				Receipt: []*pb.Receipt{
					{
						ReceiptId: "11",
						From:      "London",
						To:        "France",
						User: &pb.User{
							UserId:    "1",
							FirstName: "Alice",
							LastName:  "Smith",
							Email:     "AliceSmith@gmaiil.com",
						},
						Seat:          store.Users[0].Receipts[0].SeatNumber,
						Section:       store.Users[0].Receipts[0].SectionName,
						PricePaid:     store.Train.Price,
						BookingStatus: store.Users[0].Receipts[0].BookingStatus,
					},
				},
			},
		},
		"Sad Path - Invalid Receipt Request when request does not have UserId": {
			ShowReceiptRequest: &pb.ShowReceiptRequest{
				UserId: "",
			},
			ExpectedError: fmt.Errorf("Invalid Receipt Request"),
		},
		"Sad Path - Invalid Receipt Request when User Id is invalid": {
			ShowReceiptRequest: &pb.ShowReceiptRequest{
				UserId: "22",
			},
			ExpectedError: fmt.Errorf("User not found"),
		},
		"Sad Path - Invalid Receipt Request when request is not present": {
			ShowReceiptRequest: nil,
			ExpectedError:      fmt.Errorf("Invalid Receipt Request"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			bookingServer := &BookingServer{
				Store: store,
			}

			res, err := bookingServer.ShowReceipt(ctx, tc.ShowReceiptRequest)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected Error but didnt got error failed: %v", err)
				}
				if fmt.Sprint(err) != fmt.Sprint(tc.ExpectedError) {
					t.Errorf("Expected Error: %v, but got: %v", tc.ExpectedError, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Expected no errirr but we got an error : %v", err)

				}
			}
			assert.NotNil(t, res)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].ReceiptId, res.Receipt[0].ReceiptId)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].From, res.Receipt[0].From)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].To, res.Receipt[0].To)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].User.FirstName, res.Receipt[0].User.FirstName)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].User.LastName, res.Receipt[0].User.LastName)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].User.Email, res.Receipt[0].User.Email)
			assert.Equal(t, tc.ExpectedResponse.Receipt[0].BookingStatus, res.Receipt[0].BookingStatus)

		})
	}
}

func Test_GetSectionBookingDetails(t *testing.T) {
	store := InitializeStore()
	type test struct {
		GetSectionBookingDetailsRequest *pb.GetSectionBookingDetailsRequest
		ExpectedResponse                *pb.GetSectionBookingDetailsResponse
		ExpectedError                   error
	}
	tests := map[string]test{
		"Happy Path - Valid Section Booking Details": {
			GetSectionBookingDetailsRequest: &pb.GetSectionBookingDetailsRequest{
				SectionId: store.Train.Sections[0].Id,
			},
			ExpectedResponse: &pb.GetSectionBookingDetailsResponse{
				SeatBookings: []*pb.SeatBooking{
					{
						SeatId:        store.Train.Sections[0].Seats[0].Id,
						SeatNumber:    store.Train.Sections[0].Seats[0].SeatNumber,
						SeatAvailable: false,
						SectionId:     store.Train.Sections[0].Id,
						SectionName:   store.Train.Sections[0].Name,
						User: &pb.User{
							UserId:    "1",
							FirstName: store.Train.Sections[0].Seats[0].User.FirstName,
							LastName:  store.Train.Sections[0].Seats[0].User.LastName,
							Email:     store.Train.Sections[0].Seats[0].User.Email,
						},
					},
					{
						SeatId:        store.Train.Sections[0].Seats[1].Id,
						SeatNumber:    store.Train.Sections[0].Seats[1].SeatNumber,
						SectionId:     store.Train.Sections[0].Id,
						SectionName:   store.Train.Sections[0].Seats[1].SectionName,
						SeatAvailable: true,
					},
					{
						SeatId:        store.Train.Sections[0].Seats[2].Id,
						SeatNumber:    store.Train.Sections[0].Seats[2].SeatNumber,
						SectionId:     store.Train.Sections[0].Id,
						SectionName:   store.Train.Sections[0].Seats[2].SectionName,
						SeatAvailable: true,
					},
					{
						SeatId:        store.Train.Sections[0].Seats[3].Id,
						SeatNumber:    store.Train.Sections[0].Seats[3].SeatNumber,
						SectionId:     store.Train.Sections[0].Id,
						SectionName:   store.Train.Sections[0].Seats[3].SectionName,
						SeatAvailable: true,
					},
					{
						SeatId:        store.Train.Sections[0].Seats[4].Id,
						SeatNumber:    store.Train.Sections[0].Seats[4].SeatNumber,
						SectionId:     store.Train.Sections[0].Id,
						SectionName:   store.Train.Sections[0].Seats[4].SectionName,
						SeatAvailable: true,
					},
				},
			},
			ExpectedError: nil,
		},
		"Sad Path - Invalid Section Booking Details Request when request does not have SectionId": {
			GetSectionBookingDetailsRequest: &pb.GetSectionBookingDetailsRequest{
				SectionId: "",
			},
			ExpectedError:    fmt.Errorf("invalid Show Section-Bookings Request"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Section Booking Details Request when SectionId is invalid": {
			GetSectionBookingDetailsRequest: &pb.GetSectionBookingDetailsRequest{
				SectionId: "invalid-section-id",
			},
			ExpectedError:    fmt.Errorf("section not found for the given Section ID: invalid-section-id"),
			ExpectedResponse: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			bookingServer := &BookingServer{
				Store: store,
			}

			res, err := bookingServer.GetSectionBookingDetails(ctx, tc.GetSectionBookingDetailsRequest)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected Error but didnt got error failed: %v", err)
				}
				if fmt.Sprint(err) != fmt.Sprint(tc.ExpectedError) {
					t.Errorf("Expected Error: %v, but got: %v", tc.ExpectedError, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Expected no errirr but we got an error : %v", err)

				}
			}

			assert.NotNil(t, res)
			assert.Equal(t, len(tc.ExpectedResponse.SeatBookings), len(res.SeatBookings), "Number of seat bookings should match")
			assert.Equal(t, tc.ExpectedResponse.SeatBookings, res.SeatBookings, "Seat Bookings list should match")
		})
	}
}

func Test_DeleteBooking(t *testing.T) {
	store := InitializeStore()
	type test struct {
		DeleteBookingRequest *pb.DeleteBookingRequest
		ExpectedResponse     *pb.DeleteBookingResponse
		ExpectedError        error
	}
	tests := map[string]test{
		"Happy Path - Valid Delete Booking Request": {
			DeleteBookingRequest: &pb.DeleteBookingRequest{
				ReceiptId: "11",
			},
			ExpectedResponse: &pb.DeleteBookingResponse{
				DeleteStatus: true,
			},
			ExpectedError: nil,
		},
		"Sad Path - Invalid Delete Booking Request when request does not have ReceiptId": {
			DeleteBookingRequest: &pb.DeleteBookingRequest{
				ReceiptId: "",
			},
			ExpectedError:    fmt.Errorf("invalid Delete Booking Request"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Delete Booking Request when ReceiptId is invalid": {
			DeleteBookingRequest: &pb.DeleteBookingRequest{
				ReceiptId: "invalid-receipt-id",
			},
			ExpectedError:    fmt.Errorf("receipt not found: receipt not found for the given Receipt ID : invalid-receipt-id"),
			ExpectedResponse: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			bookingServer := &BookingServer{
				Store: store,
			}

			res, err := bookingServer.DeleteBooking(ctx, tc.DeleteBookingRequest)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected Error but didnt got error failed: %v", err)
				}
				if fmt.Sprint(err) != fmt.Sprint(tc.ExpectedError) {
					t.Errorf("Expected Error: %v, but got: %v", tc.ExpectedError, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Expected no error but we got an error : %v", err)

				}
			}

			assert.NotNil(t, res)
			assert.Equal(t, tc.ExpectedResponse.DeleteStatus, res.DeleteStatus, "Message should match")
			assert.Equal(t, "Cancelled", store.Receipts[tc.DeleteBookingRequest.ReceiptId].BookingStatus, "Booking status should be updated to Cancelled")

		})
	}
}

func Test_UpdateSeatBooking(t *testing.T) {
	store := InitializeStore()
	type test struct {
		UpdateSeatBookingRequest *pb.UpdateSeatBookingRequest
		ExpectedResponse         *pb.UpdateSeatBookingResponse
		ExpectedError            error
	}
	tests := map[string]test{
		"Happy Path - Valid Update Seat Booking Request": {
			UpdateSeatBookingRequest: &pb.UpdateSeatBookingRequest{
				ReceiptId:    "11",
				NewSeatId:    store.Train.Sections[0].Seats[1].Id,
				NewSectionId: store.Train.Sections[0].Id,
			},
			ExpectedResponse: &pb.UpdateSeatBookingResponse{
				UpdatedReceipt: &pb.Receipt{
					ReceiptId: "11",
					From:      "London",
					To:        "France",
					User: &pb.User{
						UserId:    "1",
						FirstName: "Alice",
						LastName:  "Smith",
						Email:     "AliceSmith@gmaiil.com",
					},
					Seat:          store.Train.Sections[0].Seats[1].SeatNumber,
					Section:       store.Train.Sections[0].Name,
					BookingStatus: "Confirmed",
					PricePaid:     store.Train.Price,
				},
			},
			ExpectedError: nil,
		},
		"Sad Path - Invalid Update Seat Booking Request when request does not have ReceiptId": {
			UpdateSeatBookingRequest: &pb.UpdateSeatBookingRequest{
				ReceiptId:    "",
				NewSeatId:    store.Train.Sections[0].Seats[1].Id,
				NewSectionId: store.Train.Sections[0].Id,
			},
			ExpectedError:    fmt.Errorf("Invalid Update-Seat Booking Request"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Update Seat Booking Request when request does not have NewSectionId ": {
			UpdateSeatBookingRequest: &pb.UpdateSeatBookingRequest{
				ReceiptId:    "11",
				NewSeatId:    store.Train.Sections[0].Seats[1].Id,
				NewSectionId: "",
			},
			ExpectedError:    fmt.Errorf("Invalid Update-Seat Booking Request"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Update Seat Booking Request when request does has invalid ReceiptId": {
			UpdateSeatBookingRequest: &pb.UpdateSeatBookingRequest{
				ReceiptId:    "24",
				NewSeatId:    store.Train.Sections[0].Seats[1].Id,
				NewSectionId: store.Train.Sections[0].Id,
			},
			ExpectedError:    fmt.Errorf("receipt not found: receipt not found for the given Receipt ID : 24"),
			ExpectedResponse: nil,
		},
		"Sad Path - Invalid Update Seat Booking Request when request does not have NewSeatId": {
			UpdateSeatBookingRequest: &pb.UpdateSeatBookingRequest{
				ReceiptId: "11",
				NewSeatId: "",
			},
			ExpectedError:    fmt.Errorf("Invalid Update-Seat Booking Request"),
			ExpectedResponse: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			bookingServer := &BookingServer{
				Store: store,
			}

			res, err := bookingServer.UpdateSeatBooking(ctx, tc.UpdateSeatBookingRequest)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("Expected Error but didnt got error failed: %v", err)
				}
				if fmt.Sprint(err) != fmt.Sprint(tc.ExpectedError) {
					t.Errorf("Expected Error: %v, but got: %v", tc.ExpectedError, err)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Expected no error but we got an error : %v", err)

				}
			}

			assert.NotNil(t, res)
			assert.Equal(t, tc.ExpectedResponse.UpdatedReceipt, res.UpdatedReceipt, "Update status should be true")
		})
	}
}
