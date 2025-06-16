package main

import (
	"context"
	"fmt"
	"log"

	pb "grpc-project/booking/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func PurchasingTicket(client pb.BookingServiceClient, ctx context.Context) string {
	purchaseReq := &pb.PurchaseBookingRequest{
		From:    "London",
		To:      "France",
		TrainId: "T1",
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
	fmt.Printf("Train ID: %s\n", purchaseResp.Receipt.TrainId)
	fmt.Printf("User: %s %s (%s)\n", purchaseResp.Receipt.User.FirstName, purchaseResp.Receipt.User.LastName, purchaseResp.Receipt.User.UserId)
	fmt.Printf("From: %s, To: %s\n", purchaseResp.Receipt.From, purchaseResp.Receipt.To)
	fmt.Printf("Seat: %s, Section: %s\n", purchaseResp.Receipt.Seat, purchaseResp.Receipt.Section)
	fmt.Printf("Price Paid: $%.2f, Status: %s\n", purchaseResp.Receipt.PricePaid, purchaseResp.Receipt.BookingStatus)
	return receiptId

}

func ShowReceipts(client pb.BookingServiceClient, ctx context.Context, userId string) {
	showReq := &pb.ShowReceiptRequest{
		UserId: userId,
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
		fmt.Printf("  Train ID: %s\n", receipt.TrainId)
		fmt.Printf("  Seat: %s, Section: %s\n", receipt.Seat, receipt.Section)
		fmt.Printf("  Price Paid: $%.2f, Status: %s\n", receipt.PricePaid, receipt.BookingStatus)
		fmt.Printf("\n")
	}
}

func getSectionBookingDetails(client pb.BookingServiceClient, ctx context.Context, sectionId string) []*pb.SeatBooking {
	getSectionReq := &pb.GetSectionBookingDetailsRequest{
		SectionId: "S1",
		TrainId:   "T1",
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
		fmt.Printf("- Seat Number %s (ID: %s), Available: %v, User: %s\n",
			seat.SeatNumber, seat.SeatId, seat.SeatAvailable, userInfo)
	}
	return getSectionResp.SeatBookings
}

func findNextAvailableSeatFromSection(seats []*pb.SeatBooking) (string, string) {
	var newSeatId string
	var newSeatNumber string
	var newSectionId string
	for _, seat := range seats {
		if seat.SeatAvailable {
			newSeatId = seat.SeatId
			newSeatNumber = seat.SeatNumber
			newSectionId = seat.SectionId
			fmt.Printf("Selected available seat ID: %s (Seat %s)\n", newSeatId, newSeatNumber)
			break
		}
	}

	return newSeatId, newSectionId
}

func updateBooking(client pb.BookingServiceClient, ctx context.Context, receiptId, newSeatId, newSectionId string) {
	updateReq := &pb.UpdateSeatBookingRequest{
		ReceiptId:    receiptId,
		NewSeatId:    newSeatId,
		NewSectionId: newSectionId,
		NewTrainId:   "T1",
	}
	updateResp, err := client.UpdateSeatBooking(ctx, updateReq)
	if err != nil {
		log.Fatalf("UpdateSeatBooking failed: %v", err)
	}
	fmt.Printf("Update successful!\n")
	fmt.Printf("Updated Receipt ID: %s\n", updateResp.UpdatedReceipt.ReceiptId)
	fmt.Printf("User: %s %s (%s)\n", updateResp.UpdatedReceipt.User.FirstName, updateResp.UpdatedReceipt.User.LastName, updateResp.UpdatedReceipt.User.UserId)
	fmt.Printf("From: %s, To: %s\n", updateResp.UpdatedReceipt.From, updateResp.UpdatedReceipt.To)
	fmt.Printf("Train Id: %s\n", updateResp.UpdatedReceipt.TrainId)
	fmt.Printf("Seat: %s, Section: %s\n", updateResp.UpdatedReceipt.Seat, updateResp.UpdatedReceipt.Section)
	fmt.Printf("Price Paid: $%.2f, Status: %s\n", updateResp.UpdatedReceipt.PricePaid, updateResp.UpdatedReceipt.BookingStatus)
}
func DeleteBooking(client pb.BookingServiceClient, ctx context.Context, receiptId string) {
	deleteReq := &pb.DeleteBookingRequest{
		ReceiptId: receiptId,
	}
	deleteResp, err := client.DeleteBooking(ctx, deleteReq)
	if err != nil {
		log.Fatalf("DeleteBooking failed: %v", err)
	}
	fmt.Printf("Deletion successful! Status: %v\n", deleteResp.DeleteStatus)
}

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
	receiptId := PurchasingTicket(client, ctx)
	// Step 2: Show Bob's receipts
	fmt.Println("\n ******* Step 2: Showing Bob's receipts *******")
	ShowReceipts(client, ctx, "2")

	// Step 3: Get section booking details for Section 1
	fmt.Println("\n  ******* Step 3: Getting booking details for Section 1 *******")
	_ = getSectionBookingDetails(client, ctx, "S1")
	// Step 4: Find available seat in Section 2 and update booking
	fmt.Println("\n ******* Step 4: Finding available seat in Section 2 ******")
	section2SeatBookings := getSectionBookingDetails(client, ctx, "S2")

	// Find an available seat in Section 2
	newSeatId, newSectionId := findNextAvailableSeatFromSection(section2SeatBookings)

	if newSeatId == "" || newSectionId == "" {
		log.Fatal("No available seats found in Section 2")
	}

	fmt.Println("\n ******** Step 4: Updating booking to new seat ********")
	updateBooking(client, ctx, receiptId, newSeatId, newSectionId)

	// Step 5: Delete booking
	fmt.Println("\n *********** Step 5: Deleting booking ***********")
	DeleteBooking(client, ctx, receiptId)

	/** Edge Cases to test when user is trying to modify users seat or cancel the booking again
	/* Uncomment the below use case one by one to test the edge cases
	*/
	// fmt.Println("\n *********** Update Seat booking with cancelled Receipt Id ! ***********")
	// updateBooking(client, ctx, receiptId, "a6a3aa98-822c-4fdf-bdcd-8578af42b825", "S2")

	// fmt.Println("\n *********** Deleting booking with cancelled Receipt Id ! ***********")
	// DeleteBooking(client, ctx, receiptId)

}
