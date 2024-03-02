package domain

import "time"

type Order struct {
	ID           uint      `json:"id"`
	VendorID     uint      `json:"vendor_id"`
	UserID       uint      `json:"user_id"`
	DeliveryTime int       `json:"delivery_time"`
	DeliveryAt   time.Time `json:"delivery_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
