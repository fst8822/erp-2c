package model

import (
	"time"
)

type Notification struct {
	DeliveryId int64     `json:"delivery_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     int64     `json:"user_id"`
}
