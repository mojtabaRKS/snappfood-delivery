package repositories

import (
	"context"
	"snappfood/internal/domain"

	"gorm.io/gorm"
)

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepo(db *gorm.DB) *tripRepository {
	return &tripRepository{
		db: db,
	}
}

func (r *tripRepository) FindByOrder(ctx context.Context, orderID uint) (domain.Trip, error) {
	var trip domain.Trip

	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&trip).Error

	if err != nil {
		return domain.Trip{}, err
	}

	return trip, err
}
