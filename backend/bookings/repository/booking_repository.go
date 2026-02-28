package repository

import (
	"context"
	"database/sql"
	"errors"
	"skyfox/bookings/database/common"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	ae "skyfox/error"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type bookingRepository struct {
	*common.BaseDB
}

func NewBookingRepository(db *common.BaseDB) *bookingRepository {
	return &bookingRepository{
		BaseDB: db,
	}
}

func (repo *bookingRepository) Create(ctx context.Context, b *model.Booking) error {
	dbCtx, cancel := repo.WithContext(ctx)
	defer cancel()

	if result := dbCtx.Clauses(clause.OnConflict{DoNothing: true}).Create(b); result.Error != nil {
		logger.Error("error occurred. %v", result.Error)
		if errors.Is(result.Error, gorm.ErrInvalidData) {
			return ae.UnProcessableError("BookingCreationFailed", "booking creation failed due to unknown reason", result.Error)
		} else {
			return ae.InternalServerError("BookingCreationFailed", "something went wrong", result.Error)
		}
	}
	return nil
}

func (repo *bookingRepository) BookedSeatsByShow(ctx context.Context, id int) int {
	var bookedSeats sql.NullInt32

	dbCtx, cancel := repo.WithContext(ctx)
	defer cancel()

	dbCtx.Table("booking").Select("sum(no_of_seats)").Where("show_id=?", id).Scan(&bookedSeats)

	return int(bookedSeats.Int32)
}

func (repo *bookingRepository) BookingAmountByShows(ctx context.Context, shows []int) float64 {
	var result sql.NullFloat64

	dbCtx, cancel := repo.WithContext(ctx)
	defer cancel()

	dbCtx.Table("booking").Select("sum(amount_paid)").Where("show_id IN ?", shows).Scan(&result)

	return result.Float64
}
