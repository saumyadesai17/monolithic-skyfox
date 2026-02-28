package model

type Booking struct {
	Id         int      `json:"id" gorm:"primaryKey"`
	Date       string   `json:"Date"`
	Show       Show     `json:"show"`
	ShowId     int      `gorm:"foreignKey:id"`
	Customer   Customer `json:"customer"`
	CustomerId int      `gorm:"foreignKey:id"`
	NoOfSeats  int      `json:"NoOfSeats"`
	AmountPaid float64  `json:"AmountPaid"`
}

func NewBooking(date string, show Show, customer Customer, noOfSeats int, amountPaid float64) Booking {
	return Booking{
		Date:       date,
		Show:       show,
		Customer:   customer,
		NoOfSeats:  noOfSeats,
		AmountPaid: amountPaid,
	}
}

func (Booking) TableName() string {
	return "booking"
}
