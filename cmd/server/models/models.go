package models

type Receipt struct {
	Id            string
	From          string
	To            string
	TrainId       string
	Email         string
	SeatNumber    string
	SectionName   string
	SectionId     string
	SeatId        string
	UserId        string
	BookingStatus string
}

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Receipts  []*Receipt
}

type Seat struct {
	Id            string
	SectionName   string
	SectionId     string
	SeatNumber    string
	User          *User
	SeatAvailable bool
}

type Section struct {
	Id             string
	Name           string
	Seats          []*Seat
	AvailableSeats int
}

type Train struct {
	Id       string
	From     string
	To       string
	Sections []*Section
	Price    float32
}

type Store struct {
	Trains   []*Train
	Users    []*User
	Receipts map[string]Receipt
}
