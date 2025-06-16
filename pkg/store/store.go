package store

import (
	"fmt"
	"grpc-project/cmd/server/models"
)

func GetSelectedTrainFromStore(store *models.Store, trainId string) *models.Train {
	trains := store.Trains
	if len(trains) > 0 {
		for _, train := range trains {
			if train.Id == trainId {
				return train
			}
		}
	}
	return nil
}
func GetSectionStore(store *models.Store, trainId string) []*models.Section {
	train := GetSelectedTrainFromStore(store, trainId)
	if train != nil {
		return train.Sections
	}

	return nil
}
func GetSeat(store *models.Store, trainId string, seatId string, sectionId string) *models.Seat {
	section := GetRequestedSection(store, sectionId, trainId)

	for _, seat := range section.Seats {
		if seat.Id == seatId {
			return seat
		}
	}
	return nil // Seat not found
}
func GetRequestedSection(store *models.Store, id string, trainId string) *models.Section {
	train := GetSelectedTrainFromStore(store, trainId)
	fmt.Println("****** train :", train)
	for _, section := range train.Sections {
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
	train := GetSelectedTrainFromStore(store, trainId)
	if store == nil || trainId == "" || train == nil {
		return 0.0 // Return 0 if train not found or price not set
	}

	return train.Price
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
func UpdateUserReceipts(store *models.Store, userId string, receipt *models.Receipt) *models.Store {

	for _, user := range store.Users {
		if user.Id == userId {
			user.Receipts = append(user.Receipts, receipt)
		}
	}
	return store

}
