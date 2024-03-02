package repositories

import (
	"context"
	"snappfood/internal/domain"
	"time"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindById(ctx context.Context, orderID uint) (*domain.Order, error) {
	var order domain.Order

	if err := r.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) UpdateDeliveryTime(ctx context.Context, orderId uint, deliveryAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.Order{}).Where("id = ?", orderId).Update("delivery_at", deliveryAt).Error
}
