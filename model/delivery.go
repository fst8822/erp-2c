package model

import (
	"erp-2c/lib/types"
	"strings"
	"time"
)

type deliveryStatus string

func (s deliveryStatus) IsValid(status string) error {
	switch strings.ToUpper(status) {
	case "CREATED":
		return nil
	case "SHIPPED":
		return nil
	case "DELIVERED":
		return nil
	case "CANCELLED":
		return nil
	case "ACCEPTED":
		return nil
	default:
		return types.NewAppErr(status, types.ErrUnknownStatus)
	}
}

const (
	CREATED   deliveryStatus = "CREATED"
	SHIPPED   deliveryStatus = "SHIPPED"
	DELIVERED deliveryStatus = "DELIVERED"
	CANCELLED deliveryStatus = "CANCELLED"
	ACCEPTED  deliveryStatus = "ACCEPTED"
)

type DeliveryDB struct {
	Id             int64          `db:"id"`
	RecipientGoods string         `db:"recipient_goods"`
	StatusDelivery deliveryStatus `db:"status_delivery"`
	CreatedAt      time.Time      `db:"createdAt"`
}

type DeliveryDBProductDB struct {
	Id       int64      `db:"id"`
	Delivery DeliveryDB `db:"product_id"`
	Products ProductDB  `db:"delivery_id"`
	Quantity int64      `db:"quantity"`
}

type DeliveryItem struct {
	Products ProductDomain
	Quantity int64
}

type DeliveryDomain struct {
	Id             int
	DeliveryItem   []DeliveryItem
	RecipientGoods string
	StatusDelivery deliveryStatus
	CreatedAt      time.Time
}

type DeliveryToSave struct {
	Items []struct {
		ProductId int64 `json:"product_id" validate:"required"`
		Quantity  int64 `json:"quantity" validate:"required"`
	} `json:"items" validate:"required"`
	RecipientGoods string `json:"recipient_goods" validate:"required"`
}

type UpdateStatus struct {
	Id             int64  `json:"id"`
	StatusDelivery string `json:"status_delivery" validate:"required"`
}
