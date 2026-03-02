package model

import "time"

// BookingRecord maps to the `bookings` table (migration 018).
// This is the new schema model, distinct from the legacy Booking struct.
type BookingRecord struct {
	Id            string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ShowId        string    `json:"showId" gorm:"type:uuid;not null"`
	CustomerId    string    `json:"customerId" gorm:"type:uuid;not null"`
	BookingStatus string    `json:"bookingStatus" gorm:"default:RESERVED"`
	PaymentMode   string    `json:"paymentMode,omitempty"`
	QRCodeURL     string    `json:"qrCodeUrl,omitempty" gorm:"column:qr_code_url"`
	TotalAmount   float64   `json:"totalAmount" gorm:"not null"`
	ExpiresAt     time.Time `json:"expiresAt,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (BookingRecord) TableName() string {
	return "bookings"
}

// ShowRecord maps to the new `show` table (migration 015).
// Used specifically for QR payload enrichment.
type ShowRecord struct {
	Id          string    `json:"id" gorm:"primaryKey;type:uuid"`
	MovieImdbId string    `json:"movieImdbId" gorm:"column:movie_imdb_id"`
	ScreenId    string    `json:"screenId" gorm:"column:screen_id;type:uuid"`
	TheatreId   string    `json:"theatreId" gorm:"column:theatre_id;type:uuid"`
	StartTime   time.Time `json:"startTime" gorm:"column:start_time"`
	EndTime     time.Time `json:"endTime" gorm:"column:end_time"`
	Status      string    `json:"status"`
}

func (ShowRecord) TableName() string {
	return "show"
}

// BookingSeat maps to the `booking_seats` table (migration 019).
// Used to retrieve seat IDs associated with a booking for QR payload enrichment.
type BookingSeat struct {
	Id        string  `gorm:"primaryKey;type:uuid"`
	BookingId string  `gorm:"column:booking_id;type:uuid"`
	SeatId    string  `gorm:"column:seat_id;type:uuid"`
	ShowId    string  `gorm:"column:show_id;type:uuid"`
	Price     float64 `gorm:"column:price"`
}

func (BookingSeat) TableName() string {
	return "booking_seats"
}
