package order

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"snappfood/internal/config"
	"snappfood/internal/domain"
	"snappfood/internal/dto"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func (s *orderService) newTime() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return int(rand.Float64() * 100)
}

func (s *orderService) ReportDelay(ctx context.Context, orderID uint) (*time.Time, error) {
	trip, err := s.tripRepository.FindByOrder(ctx, orderID)
	if err != nil && !(errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound)) {
		return nil, err
	}

	if trip.Status == domain.TRIP_STATUS_ASSIGNED || trip.Status == domain.TRIP_STATUS_AT_VENDOR || trip.Status == domain.TRIP_STATUS_PICKED {
		newDelivery := s.newTime()
		newTime := time.Now().Add(time.Duration(newDelivery) * time.Minute)

		if err := s.orderRepository.UpdateDeliveryTime(ctx, orderID, newTime); err != nil {
			return nil, err
		}

		return &newTime, nil
	} else {
		if err := s.redis.RPush(ctx, config.ORDER_QUEUE_KEY, orderID).Err(); err != nil {
			return nil, err
		}
	}

	if err := s.delayReportRepository.Create(ctx, orderID); err != nil {
		return nil, err
	}

	order, err := s.orderRepository.FindById(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &order.DeliveryAt, nil
}

func (s *orderService) IsReportable(ctx context.Context, orderID uint) (bool, error) {
	order, err := s.orderRepository.FindById(ctx, orderID)
	if err != nil {
		return false, err
	}

	return order.DeliveryAt.After(time.Now()) || order.DeliveryAt.Equal(time.Now()), nil
}

func (s *orderService) IsProcessing(ctx context.Context, orderID uint) (bool, error) {
	var flag bool

	reports, err := s.delayReportRepository.FindByOrder(ctx, orderID)
	if err != nil {
		return flag, err
	}

	for _, report := range reports {
		if report.Status == domain.Proccessing || report.Status == domain.Pending || report.ProccessorId == nil {
			flag = true
		}
	}

	return flag, nil
}

func (s *orderService) Proccess(ctx context.Context, employeeID int) (*domain.Order, error) {
	strOrderId, err := s.redis.LPop(ctx, config.ORDER_QUEUE_KEY).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if strOrderId == "" || errors.Is(err, redis.Nil) {
		return nil, nil
	}

	orderId, err := strconv.Atoi(strOrderId)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepository.FindById(ctx, uint(orderId))
	if err != nil {
		return nil, err
	}

	if err := s.delayReportRepository.Proccess(ctx, uint(orderId), uint(employeeID)); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) IsProccessing(ctx context.Context, employeeID int) (bool, error) {
	return s.delayReportRepository.IsProccessing(ctx, employeeID)
}

func (s *orderService) GetDelayReport(ctx context.Context) ([]dto.VendorDelay, error) {
	return s.delayReportRepository.GetDelayReport(ctx)
}
