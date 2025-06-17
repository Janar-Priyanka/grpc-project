package service

import (
	"context"
	"fmt"
	pb "grpc-project/booking/proto"
	"grpc-project/cmd/server/models"
	dataStore "grpc-project/pkg/store"

	"github.com/google/uuid"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer
	Store *models.Store
}

func (s *BookingServer) PurchaseBooking(ctx context.Context, req *pb.PurchaseBookingRequest) (*pb.PurchaseBookingResponse, error) {

	//check if request is valid
	if req == nil || req.User == nil || req.From == "" || req.To == "" {
		return nil, fmt.Errorf("Invalid Booking Request")
	}
	var finalTicketPrice float32
	var valiDiscountCode bool

	if req.DisocuntCoupon == "" {
		return nil, fmt.Errorf("please provide valid Discount coupon code")
	}

	valiDiscountCode = dataStore.CheckValidCouponCode(s.Store, req.DisocuntCoupon)

	if valiDiscountCode {
		discountRate := s.Store.DiscountCodes[req.DisocuntCoupon]
		finalTicketPrice = req.PricePaid - discountRate
	} else {
		return nil, fmt.Errorf("please provide valid Discount code")
	}
	user := s.ParseUser(req.User)

	//Allocate the seat in the available section
	seatId, sectionId := s.AllocateSeat(user)
	if seatId == "" || sectionId == "" {
		return nil, fmt.Errorf("No available seats found")
	}

	//Create a receipt for the booking
	receipt := &models.Receipt{
		Id:            uuid.New().String(),
		From:          req.From,
		To:            req.To,
		Email:         user.Email,
		UserId:        user.Id,
		SeatNumber:    dataStore.GetSeat(s.Store, seatId, sectionId).SeatNumber,
		SeatId:        seatId,
		SectionId:     sectionId,
		SectionName:   dataStore.GetSection(s.Store, sectionId).Name,
		Price:         finalTicketPrice,
		BookingStatus: "Confirmed",
	}
	if finalTicketPrice > 0 {
		receipt.Price = finalTicketPrice
	} else {
		receipt.Price = 0.0
	}

	//Update User Store Receipts
	dataStore.UpdateUserReceipts(s.Store, user.Id, receipt)
	user.Receipts = append(user.Receipts, receipt)
	s.Store.Receipts[receipt.Id] = *receipt

	//Response structure
	response := &pb.PurchaseBookingResponse{
		Receipt: &pb.Receipt{
			ReceiptId: receipt.Id,
			From:      receipt.From,
			To:        receipt.To,
			User: &pb.User{
				UserId:    receipt.UserId,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
			Seat:          receipt.SeatNumber,
			Section:       receipt.SectionName,
			PricePaid:     receipt.Price,
			BookingStatus: receipt.BookingStatus,
		},
	}

	return response, nil
}
func (s *BookingServer) ShowReceipt(ctx context.Context, req *pb.ShowReceiptRequest) (*pb.ShowReceiptResponse, error) {
	//check if request is valid
	if req == nil || req.UserId == "" {
		return nil, fmt.Errorf("Invalid Receipt Request")
	}

	//Get the user details
	user := dataStore.GetUser(s.Store, req.UserId)
	if user == nil {
		return nil, fmt.Errorf("User not found")
	}
	//map the user receipts to response struct
	response := s.MapUserReceipts(user.Receipts, user)

	return response, nil
}
func (s *BookingServer) GetSectionBookingDetails(ctx context.Context, req *pb.GetSectionBookingDetailsRequest) (*pb.GetSectionBookingDetailsResponse, error) {
	if req == nil || req.SectionId == "" {
		return nil, fmt.Errorf("invalid Show Section-Bookings Request")
	}
	section := dataStore.GetSection(s.Store, req.SectionId)
	if section == nil {
		return nil, fmt.Errorf("section not found for the given Section ID: %s", req.SectionId)
	}
	seatsList := section.Seats

	// Map the section seats to response struct
	var pbSeats []*pb.SeatBooking
	for _, seat := range seatsList {
		seatDetails := &pb.SeatBooking{
			SeatId:        seat.Id,
			SeatNumber:    seat.SeatNumber,
			SectionId:     seat.SectionId,
			SectionName:   seat.SectionName,
			SeatAvailable: seat.SeatAvailable,
		}
		if seat.User != nil {
			seatDetails.User = &pb.User{
				UserId:    seat.User.Id,
				FirstName: seat.User.FirstName,
				LastName:  seat.User.LastName,
				Email:     seat.User.Email,
			}
		}
		pbSeats = append(pbSeats, seatDetails)
	}

	// Create the response structure
	response := &pb.GetSectionBookingDetailsResponse{
		SeatBookings: pbSeats,
	}
	return response, nil
}
func (s *BookingServer) DeleteBooking(ctx context.Context, req *pb.DeleteBookingRequest) (*pb.DeleteBookingResponse, error) {

	//check if request is valid
	if req == nil || req.ReceiptId == "" {
		return nil, fmt.Errorf("invalid Delete Booking Request")
	}

	//validate the receipt
	receipt, err := dataStore.CheckValidReceipt(s.Store, req.ReceiptId)
	if err != nil {
		return nil, fmt.Errorf("receipt not found: %v", err)

	}
	if receipt.BookingStatus == "Cancelled" {
		return nil, fmt.Errorf("your booking is already cancelled")
	}

	//Get the seat details from the receipt details
	seat := dataStore.GetSeat(s.Store, receipt.SeatId, receipt.SectionId)

	//Reset Seat status to available
	seat.SeatAvailable = true
	seat.User = nil

	//Update the section Seat Availability count
	section := dataStore.GetSection(s.Store, receipt.SectionId)
	section.AvailableSeats--

	//Update users store Receipts for cancellation
	user := dataStore.GetUser(s.Store, receipt.UserId)
	if user != nil {
		for _, receipt := range user.Receipts {
			if receipt.Id == req.ReceiptId {
				receipt.BookingStatus = "Cancelled"
				break
			}
		}
	}

	//Mark booking status in the receipts store
	dataStore.CancelReceiptsFromStore(s.Store, req.ReceiptId)

	//Response structure
	response := &pb.DeleteBookingResponse{
		DeleteStatus: true,
	}

	return response, nil
}
func (s *BookingServer) UpdateSeatBooking(ctx context.Context, req *pb.UpdateSeatBookingRequest) (*pb.UpdateSeatBookingResponse, error) {

	if req == nil || req.ReceiptId == "" || req.NewSeatId == "" || req.NewSectionId == "" {
		return nil, fmt.Errorf("Invalid Update-Seat Booking Request")
	}
	receipt, err := dataStore.CheckValidReceipt(s.Store, req.ReceiptId)
	if err != nil {
		return nil, fmt.Errorf("receipt not found: %v", err)
	}
	if receipt.BookingStatus == "Cancelled" {
		return nil, fmt.Errorf("your booking is already cancelled, hence cannot update user seat")
	}
	user := dataStore.GetUser(s.Store, receipt.UserId)
	//Check if the new seat is available
	newSeat := dataStore.GetSeat(s.Store, req.NewSeatId, req.NewSectionId)
	if newSeat == nil || !newSeat.SeatAvailable {
		return nil, fmt.Errorf("requested seat is not available")
	}
	newSeat.SeatAvailable = false
	newSeat.User = user
	newSeatSection := dataStore.GetSection(s.Store, newSeat.SectionId)
	if newSeatSection != nil {
		newSeatSection.AvailableSeats--
	}

	//Update the old seat to available
	oldSeat := dataStore.GetSeat(s.Store, receipt.SeatId, receipt.SectionId)
	if oldSeat != nil {
		oldSeat.SeatAvailable = true
		oldSeat.User = nil
	}

	oldSeatSection := dataStore.GetSection(s.Store, receipt.SectionId)
	oldSeatSection.AvailableSeats++

	//Update the receipt with new seat details in the Store
	receipt.SeatId = newSeat.Id
	receipt.SeatNumber = newSeat.SeatNumber
	receipt.SectionId = newSeat.SectionId
	receipt.SectionName = newSeat.SectionName
	s.Store.Receipts[receipt.Id] = *receipt

	//Update the Store Receipt receipts in the Store
	for i, r := range s.Store.Receipts {
		if r.Id == receipt.Id {
			s.Store.Receipts[i] = *receipt
			break
		}
	}

	//Update the users store receipts
	for i, user := range s.Store.Users {
		if user.Id == user.Id {
			for j, userReceipt := range user.Receipts {
				if userReceipt.Id == receipt.Id {
					s.Store.Users[i].Receipts[j] = receipt
					break
				}
			}
			break
		}
	}

	//Response structure
	response := &pb.UpdateSeatBookingResponse{
		UpdatedReceipt: &pb.Receipt{
			ReceiptId: receipt.Id,
			From:      receipt.From,
			To:        receipt.To,
			User: &pb.User{
				UserId:    user.Id,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
			Seat:          receipt.SeatNumber,
			Section:       receipt.SectionName,
			PricePaid:     dataStore.GetPriceFromReceipts(s.Store, receipt.Id),
			BookingStatus: receipt.BookingStatus,
		},
	}
	return response, nil
}

/*Helper Methods*/
func (s *BookingServer) MapUserReceipts(userReceipts []*models.Receipt, user *models.User) *pb.ShowReceiptResponse {
	var responseStruct *pb.ShowReceiptResponse
	var pbReceipts []*pb.Receipt
	for _, receipt := range userReceipts {
		pbReceipts = append(pbReceipts, &pb.Receipt{
			ReceiptId: receipt.Id,
			From:      receipt.From,
			To:        receipt.To,
			User: &pb.User{
				UserId:    user.Id,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
			Seat:          receipt.SeatNumber,
			Section:       receipt.SectionName,
			PricePaid:     dataStore.GetPriceFromReceipts(s.Store, receipt.Id),
			BookingStatus: receipt.BookingStatus,
		})
	}
	responseStruct = &pb.ShowReceiptResponse{
		Receipt: pbReceipts,
	}
	return responseStruct
}
func (s *BookingServer) AllocateSeat(user *models.User) (string, string) {
	var sections []*models.Section = dataStore.GetSectionStore(s.Store)

	for _, section := range sections {
		if section.AvailableSeats > 0 {
			nextAvailableSeatId := s.GetNextAvailableSeat(section)
			if nextAvailableSeatId != "" {
				seat := dataStore.GetSeat(s.Store, nextAvailableSeatId, section.Id)
				if seat != nil {
					seat.SeatAvailable = false
					seat.User = user
				}
				section.AvailableSeats--
				return nextAvailableSeatId, section.Id
			}
		}
	}
	return "", ""
}
func (s *BookingServer) GetNextAvailableSeat(section *models.Section) string {

	for _, seat := range section.Seats {
		if seat.SeatAvailable {
			return seat.Id // Return the first available seat
		}
	}
	return ""
}
func (s *BookingServer) ParseUser(user *pb.User) *models.User {
	if user == nil {
		return nil
	}
	return &models.User{
		Id:        user.GetUserId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
	}
}
