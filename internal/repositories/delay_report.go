package repositories

import (
	"context"
	"database/sql"
	"errors"
	"snappfood/internal/domain"
	"snappfood/internal/dto"
	"time"

	"gorm.io/gorm"
)

type delayReportRepository struct {
	db *gorm.DB
}

func NewDelayRepo(db *gorm.DB) *delayReportRepository {
	return &delayReportRepository{
		db: db,
	}
}

func (r *delayReportRepository) Create(ctx context.Context, orderID uint) error {
	model := domain.DelayReport{
		OrderID: orderID,
		Status:  domain.Pending,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *delayReportRepository) FindByOrder(ctx context.Context, orderId uint) ([]domain.DelayReport, error) {
	var reports []domain.DelayReport

	if err := r.db.WithContext(ctx).Model(&domain.DelayReport{}).Where("order_id = ?", orderId).Scan(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *delayReportRepository) Proccess(ctx context.Context, orderID, employeeID uint) error {
	var delay domain.DelayReport

	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Where("status = ? ", domain.Pending).
		First(&delay).
		Error

	if err != nil {
		return nil
	}

	return r.db.WithContext(ctx).
		Model(&domain.DelayReport{}).
		Where("id = ?", delay.ID).
		Update("proccessor_id", employeeID).
		Update("status", domain.Proccessing).
		Error
}

func (r *delayReportRepository) IsProccessing(ctx context.Context, employeeID int) (bool, error) {
	var delay domain.DelayReport
	err := r.db.WithContext(ctx).
		Where("proccessor_id = ?", employeeID).
		Where("status = ?", domain.Proccessing).
		First(&delay).
		Error

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows)) {
		return false, err
	}

	return delay.ID != 0, nil
}

func (r *delayReportRepository) GetDelayReport(ctx context.Context) ([]dto.VendorDelay, error) {

	lastWeek := time.Now().AddDate(0, 0, -7)

	// Query to get vendors with their total delay since last week
	var vendorDelays []dto.VendorDelay
	err := r.db.Table("vendors").
		Select("vendors.id AS vendor_id, vendors.name AS vendor_name, COALESCE(SUM(delay_reports.delay_time), 0) AS total_delay").
		Joins("LEFT JOIN orders ON vendors.id = orders.vendor_id").
		Joins("LEFT JOIN delay_reports ON orders.id = delay_reports.order_id").
		Where("delay_reports.created_at >= ?", lastWeek).
		Or("orders.created_at >= ?", lastWeek).
		Group("vendors.id, vendors.name").
		Order("total_delay DESC").
		Scan(&vendorDelays).
		Error

	if err != nil {
		return nil, err
	}

	return vendorDelays, nil
}
