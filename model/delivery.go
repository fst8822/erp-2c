package model

import "time"

const (
	CREATED   = "CREATED"
	SHIPPED   = "SHIPPED"
	DELIVERED = "DELIVERED"
	CANCELLED = "CANCELLED"
	ACCEPTED  = "ACCEPTED"
)

type DeliveryDB struct {
	Id             int64     `db:"id"`
	RecipientGoods string    `db:"recipient_goods"`
	StatusDelivery string    `db:"status_delivery"`
	CreatedAt      time.Time `db:"createdAt"`
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
	StatusDelivery string
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
