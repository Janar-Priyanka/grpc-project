package models

type Receipt struct {
	Id            string
	From          string
	To            string
	Email         string
	SeatNumber    string
	SectionName   string
	SectionId     string
	SeatId        string
	UserId        string
	BookingStatus string
	Price float32
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
	Id             string //A or B
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
	Train         Train
	Users         []*User
	DiscountCodes map[string]float32
	Receipts      map[string]Receipt
}

// We want a first class section of the train: section A.  
// The ticket price would  be $40.