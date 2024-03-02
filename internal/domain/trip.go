package domain

import "time"

type TripStatus string

const (
	TRIP_STATUS_ASSIGNED  = "ASSIGNED"
	TRIP_STATUS_AT_VENDOR = "AT_VENDOR"
	TRIP_STATUS_PICKED    = "PICKED"
	TRIP_STATUS_DELIVERED = "DELIVERED"
)

type Trip struct {
	ID         uint
	AssigneeID uint
	Status     TripStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
