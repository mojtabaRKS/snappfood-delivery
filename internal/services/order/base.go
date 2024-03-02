package order

import (
	"context"
	"snappfood/internal/domain"
	"snappfood/internal/dto"
	"time"

	"github.com/redis/go-redis/v9"
)

type orderService struct {
	orderRepository       orderRepositoryInterface
	delayReportRepository delayReportRepositoryInterface
	tripRepository        tripRepositoryInterface
	redis                 *redis.Client
}

type orderRepositoryInterface interface {
	FindById(ctx context.Context, orderID uint) (*domain.Order, error)
	UpdateDeliveryTime(ctx context.Context, orderId uint, deliveryAt time.Time) error
}

type delayReportRepositoryInterface interface {
	Create(ctx context.Context, orderID uint) error
	FindByOrder(ctx context.Context, orderID uint) ([]domain.DelayReport, error)
	Proccess(ctx context.Context, orderID, employeeID uint) error
	IsProccessing(ctx context.Context, employeeID int) (bool, error)
	GetDelayReport(ctx context.Context) ([]dto.VendorDelay, error)
}

type tripRepositoryInterface interface {
	FindByOrder(ctx context.Context, orderID uint) (domain.Trip, error)
}

func New(
	orderRepo orderRepositoryInterface,
	delayReportRepo delayReportRepositoryInterface,
	tripRepo tripRepositoryInterface,
	redis *redis.Client,
) *orderService {
	return &orderService{
		orderRepository:       orderRepo,
		delayReportRepository: delayReportRepo,
		tripRepository:        tripRepo,
		redis:                 redis,
	}
}
