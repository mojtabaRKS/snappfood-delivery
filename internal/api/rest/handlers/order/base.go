package order

import (
	"context"
	"snappfood/internal/domain"
	"time"
)

type OrderHandler struct {
	orderService orderServiceInterface
}

type orderServiceInterface interface {
	ReportDelay(ctx context.Context, orderID uint) (*time.Time, error)
	IsReportable(ctx context.Context, orderID uint) (bool, error)
	IsProcessing(ctx context.Context, orderID uint) (bool, error)
	Proccess(ctx context.Context, employeeID int) (*domain.Order, error)
	IsProccessing(ctx context.Context, employeeID int) (bool, error)
}

func New(orderService orderServiceInterface) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}
