package main

import (
	"context"
	"fmt"
	"log"

	pb "grpc-project/booking/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create BookingService client
	client := pb.NewBookingServiceClient(conn)
	fmt.Println("Connected to gRPC server at localhost:8080")

	ctx := context.Background()

	// Step 1: Purchase a ticket for Bob
	fmt.Println("\n ********* Step 1: Purchasing a ticket for Bob  **********")
	purchaseReq := &pb.PurchaseBookingRequest{
		From: "London",
		To:   "France",
		User: &pb.User{
			UserId:    "2",
			FirstName: "Bob",
			LastName:  "Johnson",
			Email:     "bobthebuilder@gmail.com",
		},
		PricePaid: 20.0,
	}
	purchaseResp, err := client.PurchaseBooking(ctx, purchaseReq)
	if err != nil {
		log.Fatalf("PurchaseBooking failed: %v", err)
	}
	receiptId := purchaseResp.Receipt.ReceiptId
	fmt.Printf("Purchase successful!\n")
	fmt.Printf("Receipt ID: %s\n", receiptId)
	fmt.Printf("User: %s %s (%s)\n", purchaseResp.Receipt.User.FirstName, purchaseResp.Receipt.User.LastName, purchaseResp.Receipt.User.UserId)
	fmt.Printf("From: %s, To: %s\n", purchaseResp.Receipt.From, purchaseResp.Receipt.To)
	fmt.Printf("Seat: %s, Section: %s\n", purchaseResp.Receipt.Seat, purchaseResp.Receipt.Section)
	fmt.Printf("Price Paid: $%.2f, Status: %s\n", purchaseResp.Receipt.PricePaid, purchaseResp.Receipt.BookingStatus)

	// Step 2: Show Bob's receipts
	fmt.Println("\n ******* Step 2: Showing Bob's receipts *******")
	showReq := &pb.ShowReceiptRequest{
		UserId: "2",
	}
	showResp, err := client.ShowReceipt(ctx, showReq)
	if err != nil {
		log.Fatalf("ShowReceipt failed: %v", err)
	}
	fmt.Printf("Receipts for user %s:\n", showReq.UserId)
	for _, receipt := range showResp.Receipt {
		fmt.Printf("- Receipt ID: %s\n", receipt.ReceiptId)
		fmt.Printf("  User: %s %s (%s)\n", receipt.User.FirstName, receipt.User.LastName, receipt.User.UserId)
		fmt.Printf("  From: %s, To: %s\n", receipt.From, receipt.To)
		fmt.Printf("  Seat: %s, Section: %s\n", receipt.Seat, receipt.Section)
		fmt.Printf("  Price Paid: $%.2f, Status: %s\n", receipt.PricePaid, receipt.BookingStatus)
	}

	// Step 3: Get section booking details for Section 1
	fmt.Println("\n  ******* Step 3: Getting booking details for Section 1 *******")
	getSectionReq := &pb.GetSectionBookingDetailsRequest{
		SectionId: "S1",
	}
	getSectionResp, err := client.GetSectionBookingDetails(ctx, getSectionReq)
	if err != nil {
		log.Fatalf("GetSectionBookingDetails failed: %v", err)
	}
	fmt.Printf("Section Booking Details for Section ID %s:\n", getSectionReq.SectionId)
	for _, seat := range getSectionResp.SeatBookings {
		userInfo := "None"
		if seat.User != nil {
			userInfo = fmt.Sprintf("%s %s (%s)", seat.User.FirstName, seat.User.LastName, seat.User.UserId)
		}
		fmt.Printf("- Seat %s (ID: %s), Available: %v, User: %s\n",
			seat.SeatNumber, seat.SeatId, seat.SeatAvailable, userInfo)
	}

	// Step 4: Find available seat in Section 2 and update booking
	fmt.Println("\n ******* Step 4: Finding available seat in Section 2 ******")
	getSection2Req := &pb.GetSectionBookingDetailsRequest{
		SectionId: "S2",
	}
	getSection2Resp, err := client.GetSectionBookingDetails(ctx, getSection2Req)
	if err != nil {
		log.Fatalf("GetSectionBookingDetails for Section 2 failed: %v", err)
	}
	for _, seat := range getSection2Resp.SeatBookings {
		userInfo := "None"
		if seat.User != nil {
			userInfo = fmt.Sprintf("%s %s (%s)", seat.User.FirstName, seat.User.LastName, seat.User.UserId)
		}
		fmt.Printf("- Seat %s (ID: %s), Available: %v, User: %s\n",
			seat.SeatNumber, seat.SeatId, seat.SeatAvailable, userInfo)
	}
	var newSeatId string
	var newSeatNumber string
	var newSectionId string
	for _, seat := range getSection2Resp.SeatBookings {
		if seat.SeatAvailable {
			newSeatId = seat.SeatId
			newSeatNumber = seat.SeatNumber
			newSectionId = seat.SectionId
			fmt.Printf("Selected available seat ID: %s (Seat %s)\n", newSeatId, newSeatNumber)
			break
		}
	}
	if newSeatId == "" {
		log.Fatal("No available seats found in Section 2")
	}

	fmt.Println("\n ******** Step 4: Updating booking to new seat ********")
	updateReq := &pb.UpdateSeatBookingRequest{
		ReceiptId:    receiptId,
		NewSeatId:    newSeatId,
		NewSectionId: newSectionId,
	}
	updateResp, err := client.UpdateSeatBooking(ctx, updateReq)
	if err != nil {
		log.Fatalf("UpdateSeatBooking failed: %v", err)
	}
	fmt.Printf("Update successful!\n")
	fmt.Printf("Updated Receipt ID: %s\n", updateResp.UpdatedReceipt.ReceiptId)
	fmt.Printf("User: %s %s (%s)\n", updateResp.UpdatedReceipt.User.FirstName, updateResp.UpdatedReceipt.User.LastName, updateResp.UpdatedReceipt.User.UserId)
	fmt.Printf("From: %s, To: %s\n", updateResp.UpdatedReceipt.From, updateResp.UpdatedReceipt.To)
	fmt.Printf("Seat: %s, Section: %s\n", updateResp.UpdatedReceipt.Seat, updateResp.UpdatedReceipt.Section)
	fmt.Printf("Price Paid: $%.2f, Status: %s\n", updateResp.UpdatedReceipt.PricePaid, updateResp.UpdatedReceipt.BookingStatus)

	// Step 5: Delete booking
	fmt.Println("\n *********** Step 5: Deleting booking ***********")
	deleteReq := &pb.DeleteBookingRequest{
		ReceiptId: receiptId,
	}
	deleteResp, err := client.DeleteBooking(ctx, deleteReq)
	if err != nil {
		log.Fatalf("DeleteBooking failed: %v", err)
	}
	fmt.Printf("Deletion successful! Status: %v\n", deleteResp.DeleteStatus)
}
