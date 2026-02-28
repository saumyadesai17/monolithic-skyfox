package repository

import (
	// "skyfox/bookings/database/connection"
	"context"
	"skyfox/bookings/database/common"
	"skyfox/bookings/model"
	ae "skyfox/error"

	"gorm.io/gorm/clause"
)

type CustomerRepository interface {
	Create(context.Context, *model.Customer) error
}

type customerRepository struct {
	*common.BaseDB
}

func NewCustomerRepository(db *common.BaseDB) CustomerRepository {
	return &customerRepository{
		BaseDB: db,
	}
}

func (repo customerRepository) Create(ctx context.Context, c *model.Customer) error {
	dbCtx, cancel := repo.WithContext(ctx)
	defer cancel()

	// If we're creating a new customer (Id = 0), don't specify the Id field
	// to let the database generate it
	if c.Id == 0 {
		result := dbCtx.Omit("Id").Clauses(clause.OnConflict{DoNothing: true}).Create(c)
		if result.Error != nil {
			return ae.UnProcessableError("CustomerCreationFailed", "Customer creation failed due to unknown reason", result.Error)
		}
		if result.RowsAffected == 0 {
			return ae.UnProcessableError("CustomerAlreadyExist", "Customer already exist. Duplicate record", nil)
		}
	} else {
		// If Id is specified, use Save which can handle both create and update
		result := dbCtx.Clauses(clause.OnConflict{DoNothing: true}).Save(c)
		if result.Error != nil {
			return ae.UnProcessableError("CustomerCreationFailed", "Customer creation failed due to unknown reason", result.Error)
		}
		if result.RowsAffected == 0 {
			return ae.UnProcessableError("CustomerAlreadyExist", "Customer already exist. Duplicate record", nil)
		}
	}

	return nil
}
