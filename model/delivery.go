package model

import (
	"erp-2c/lib/types"
	"time"
)

type deliveryStatus string

func (s deliveryStatus) IsValid() error {
	switch s {
	case CREATED, SHIPPED, DELIVERED, CANCELLED, ACCEPTED:
		return nil
	default:
		return types.NewAppErr(string(s), types.ErrUnknownStatus)
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
	Address        string         `db:"address"`
	StatusDelivery deliveryStatus `db:"status_delivery"`
	CreatedAt      time.Time      `db:"created_at"`
	TotalAmount    int64          `db:"total_amount"`
}

type DeliveryProductDB struct {
	Id          int64 `db:"id"`
	DeliveryId  int64 `db:"delivery_id"`
	ProductId   int64 `db:"product_id"`
	Quantity    int64 `db:"quantity"`
	ItemPrice   int64 `db:"item_price"`
	TotalAmount int64 `db:"total_amount"`
}

type DeliveryItem struct {
	Product   ProductDomain
	ItemPrice int64
	Quantity  int64
}

func (i *DeliveryItem) TotalAmount() int64 {
	return i.ItemPrice * i.Quantity
}

type DeliveryDomain struct {
	Id             int64
	DeliveryItems  []DeliveryItem
	RecipientGoods string
	Address        string
	StatusDelivery deliveryStatus
	CreatedAt      time.Time
}

func (d *DeliveryDomain) CalculateTotal() int64 {
	var total int64
	for _, item := range d.DeliveryItems {
		total += item.TotalAmount()
	}
	return total
}

type DeliveryToSave struct {
	Items          []DeliveryItemToSave `json:"items" validate:"required"`
	RecipientGoods string               `json:"recipient_goods" validate:"required"`
	Address        string               `json:"address" validate:"required"`
}

type DeliveryItemToSave struct {
	ProductId int64 `json:"product_id" validate:"required"`
	ItemPrice int64 `json:"itemPrice" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required"`
}

type UpdateStatus struct {
	DeliveryId     int64          `json:"id" validate:"required"`
	StatusDelivery deliveryStatus `json:"status_delivery" validate:"required"`
}

func (d *DeliveryDomain) MapToDB() *DeliveryDB {
	return nil
}
func (d *DeliveryDB) MapToDomain() *DeliveryDomain {
	return nil
}
