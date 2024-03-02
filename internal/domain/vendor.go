package domain

import "time"

type Vendor struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
