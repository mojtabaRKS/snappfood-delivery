package vendor

import (
	"context"
	"snappfood/internal/dto"
)

type VendorHandler struct {
	orderService orderServiceInterface
}

type orderServiceInterface interface {
	GetDelayReport(ctx context.Context) ([]dto.VendorDelay, error)
}

func New(orderService orderServiceInterface) *VendorHandler {
	return &VendorHandler{
		orderService: orderService,
	}
}
