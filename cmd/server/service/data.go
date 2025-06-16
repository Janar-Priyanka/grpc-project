package service

import (
	"fmt"
	"grpc-project/cmd/server/models"

	"github.com/google/uuid"
)

func GenerateNewTrain(store *models.Store, sectionCount, seatCount int, fromDestination, toDestination string, price float32) (*models.Store, error) {
	if fromDestination == "" || toDestination == "" {
		return nil, fmt.Errorf(" from or to details are required to create a train")
	}
	if sectionCount == 0 {
		seatCount = 2
	}
	if seatCount == 0 {
		seatCount = 20
	}
	if price == 0.0 {
		price = float32(20.0)
	}
	fmt.Println("******", price)
	newTrain := &models.Train{
		Id:    "T1",
		From:  fromDestination,
		To:    toDestination,
		Price: float32(price),
	}

	newTrain.Sections = generateSectionsForTrain(store, sectionCount, seatCount)

	store.Trains = append(store.Trains, newTrain)
	return store, nil
}
func generateSeatForSection(sectionId string, sectionName string, seatCount int) []*models.Seat {
	seats := []*models.Seat{}

	for i := 0; i < seatCount; i++ {
		seat := &models.Seat{
			Id:            uuid.New().String(),
			SectionName:   sectionName,
			SectionId:     sectionId,
			SeatNumber:    "Seat" + fmt.Sprint(i+1),
			User:          nil,
			SeatAvailable: true,
		}
		seats = append(seats, seat)
	}
	return seats
}

func generateSectionsForTrain(store *models.Store, sectionCount int, seatCount int) []*models.Section {
	sections := []*models.Section{}

	for i := 0; i < sectionCount; i++ {
		section := &models.Section{
			Id:             "S" + fmt.Sprint(i+1),
			Name:           "S" + fmt.Sprint(i+1),
			AvailableSeats: seatCount,
		}
		section.Seats = generateSeatForSection(section.Id, section.Name, seatCount)
		sections = append(sections, section)
	}
	return sections
}
