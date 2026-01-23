package model

import (
	"erp-2c/lib/types"
	"strings"
	"time"
)

type DeliveryStatus string

func (s DeliveryStatus) ParseDeliverStatus(status string) (DeliveryStatus, error) {
	switch strings.ToUpper(status) {
	case "CREATED":
		return CREATED, nil
	case "SHIPPED":
		return SHIPPED, nil
	case "DELIVERED":
		return DELIVERED, nil
	case "CANCELLED":
		return CANCELLED, nil
	case "ACCEPTED":
		return ACCEPTED, nil
	default:
		return "", types.NewAppErr(status, types.ErrUnknownStatus)
	}
}

const (
	CREATED   DeliveryStatus = "CREATED"
	SHIPPED   DeliveryStatus = "SHIPPED"
	DELIVERED DeliveryStatus = "DELIVERED"
	CANCELLED DeliveryStatus = "CANCELLED"
	ACCEPTED  DeliveryStatus = "ACCEPTED"
)

type DeliveryDB struct {
	Id             int64          `db:"id"`
	RecipientGoods string         `db:"recipient_goods"`
	StatusDelivery DeliveryStatus `db:"status_delivery"`
	CreatedAt      time.Time      `db:"createdAt"`
}

type DeliveryProductDB struct {
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
	StatusDelivery DeliveryStatus
	CreatedAt      time.Time
}

type DeliveryToSave struct {
	Items []struct {
		ProductId int64 `json:"product_id" validate:"required"`
		Quantity  int64 `json:"quantity" validate:"required"`
	} `json:"items" validate:"required"`
	RecipientGoods string `json:"recipient_goods" validate:"required"`
	StatusDelivery string `json:"status_delivery" validate:"required"`
}
