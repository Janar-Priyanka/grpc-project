package store

import (
	"fmt"
	"grpc-project/cmd/server/models"
)

func GetSectionStore(store *models.Store) []*models.Section {
	return store.Train.Sections
	// No sections found
}
func GetSeat(store *models.Store, seatId string, sectionId string) *models.Seat {
	section := GetSection(store, sectionId)

	for _, seat := range section.Seats {
		if seat.Id == seatId {
			return seat
		}
	}
	return nil // Seat not found
}
func GetSection(store *models.Store, id string) *models.Section {
	for _, section := range store.Train.Sections {
		if section.Id == id {
			return section
		}
	}
	return nil // Section not found
}

func GetUser(store *models.Store, userId string) *models.User {
	if store != nil && store.Users != nil {
		for _, user := range store.Users {
			if user.Id == userId {
				return user
			}
		}
	}
	return nil
}

func GetPrice(store *models.Store, trainId string) float32 {
	if store != nil && store.Train.Id == trainId {
		return store.Train.Price
	}
	return 0.0 // Return 0 if train not found or price not set
}
func CheckValidReceipt(store *models.Store, receiptId string) (*models.Receipt, error) {
	if receipt, exists := store.Receipts[receiptId]; exists {
		return &receipt, nil
	}
	return nil, fmt.Errorf("receipt not found for the given Receipt ID : %s", receiptId)
}
func CancelReceiptsFromStore(store *models.Store, receiptId string) {
	if receipt, exists := store.Receipts[receiptId]; exists {
		receipt.BookingStatus = "Cancelled"
		store.Receipts[receiptId] = receipt
	}
}
func UpdateUserReceipts(store *models.Store, userId string, receipt *models.Receipt) {

	for _, user := range store.Users {
		if user.Id == userId {
			user.Receipts = append(user.Receipts, receipt)
			return
		}
	}

}
