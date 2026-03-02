package repository

import (
	"context"
	"errors"
	"fmt"
	"skyfox/bookings/database/common"
	"skyfox/bookings/model"
	ae "skyfox/error"

	"gorm.io/gorm"
)

// QRBookingRepository is the minimal DB interface that QRService needs.
// It is kept separate from any future full BookingRepository so the QR feature
// can be developed and tested independently.
type QRBookingRepository interface {
	FindBookingByID(ctx context.Context, id string) (*model.BookingRecord, error)
	FindShowByID(ctx context.Context, id string) (*model.ShowRecord, error)
	FindSeatsByBookingID(ctx context.Context, bookingID string) ([]string, error)
	UpdateQRCodeURL(ctx context.Context, id string, qrURL string) error
}

type qrBookingRepository struct {
	*common.BaseDB
}

func NewQRBookingRepository(db *common.BaseDB) QRBookingRepository {
	return &qrBookingRepository{BaseDB: db}
}

func (r *qrBookingRepository) FindBookingByID(ctx context.Context, id string) (*model.BookingRecord, error) {
	var booking model.BookingRecord

	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Where("id = ?", id).First(&booking)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // caller checks nil → 404
		}
		return nil, ae.InternalServerError("InternalServerError", "something went wrong",
			fmt.Errorf("FindBookingByID failed: %w", result.Error))
	}
	return &booking, nil
}

func (r *qrBookingRepository) UpdateQRCodeURL(ctx context.Context, id string, qrURL string) error {
	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Model(&model.BookingRecord{}).Where("id = ?", id).Update("qr_code_url", qrURL)
	if result.Error != nil {
		return ae.InternalServerError("InternalServerError", "failed to update QR code URL",
			fmt.Errorf("UpdateQRCodeURL failed: %w", result.Error))
	}
	return nil
}

func (r *qrBookingRepository) FindShowByID(ctx context.Context, id string) (*model.ShowRecord, error) {
	var show model.ShowRecord

	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Where("id = ?", id).First(&show)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ae.InternalServerError("InternalServerError", "failed to fetch show",
			fmt.Errorf("FindShowByID failed: %w", result.Error))
	}
	return &show, nil
}

func (r *qrBookingRepository) FindSeatsByBookingID(ctx context.Context, bookingID string) ([]string, error) {
	var seats []model.BookingSeat

	db, cancel := r.WithContext(ctx)
	defer cancel()

	result := db.Where("booking_id = ?", bookingID).Find(&seats)
	if result.Error != nil {
		return nil, ae.InternalServerError("InternalServerError", "failed to fetch booking seats",
			fmt.Errorf("FindSeatsByBookingID failed: %w", result.Error))
	}

	seatIDs := make([]string, len(seats))
	for i, s := range seats {
		seatIDs[i] = s.SeatId
	}
	return seatIDs, nil
}
