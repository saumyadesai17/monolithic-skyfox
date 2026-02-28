package repository

import (
	"context"
	"fmt"
	"skyfox/bookings/database/common"

	"skyfox/bookings/model"
	ae "skyfox/error"
)

type showRepository struct {
	*common.BaseDB
}

func NewShowRepository(db *common.BaseDB) *showRepository {
	return &showRepository{
		BaseDB: db,
	}
}

func (repo showRepository) GetAllShowsOn(ctx context.Context, date string) ([]model.Show, error) {
	var shows []model.Show

	db, cancel := repo.WithContext(ctx)
	defer cancel()

	err := db.Model(&model.Show{}).Preload("Slot").Where("date=?", date).Find(&shows).Error

	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, ae.InternalServerError("InternalServerError", "query could not be processed", fmt.Errorf("error: %v", err))
		}
	}

	return shows, nil
}

func (repo showRepository) FindById(ctx context.Context, id int) (model.Show, error) {
	var show model.Show

	db, cancel := repo.WithContext(ctx)
	defer cancel()

	result := db.Model(&model.Show{}).Preload("Slot").Where("id=?", id).First(&show)

	if result.Error != nil {
		return model.Show{}, ae.NotFoundError("ShowNotFound", fmt.Sprintf("Show not found for id : %d", id), result.Error)
	}
	return show, nil
}
