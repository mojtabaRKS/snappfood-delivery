package domain

import "time"

type ProccessStatus string

const (
	Proccessing ProccessStatus = "PROCCESSING"
	Proccessed  ProccessStatus = "PROCCESSED"
	Pending     ProccessStatus = "PENDING"
)

type DelayReport struct {
	ID           uint
	OrderID      uint
	ProccessorId *uint
	Status       ProccessStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
