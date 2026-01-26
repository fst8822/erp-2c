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
	Recipient      string         `db:"recipient"`
	Address        string         `db:"address"`
	StatusDelivery deliveryStatus `db:"status"`
	CreatedAt      time.Time      `db:"created_at"`
}

type DeliveryProductDB struct {
	Id          int64 `db:"id"`
	DeliveryId  int64 `db:"delivery_id"`
	ProductId   int64 `db:"product_id"`
	ItemPrice   int64 `db:"unit_price"`
	Quantity    int64 `db:"quantity"`
	TotalAmount int64 `db:"total_amount"`
}

type DeliveryItem struct {
	ProductId int64 `json:"product_id"`
	ItemPrice int64 `json:"item_price"`
	Quantity  int64 `json:"quantity"`
}

func (i *DeliveryItem) TotalAmount() int64 {
	return i.ItemPrice * i.Quantity
}

type DeliveryDomain struct {
	Id             int64          `json:"id"`
	DeliveryItems  []DeliveryItem `json:"delivery_items"`
	Recipient      string         `json:"recipient"`
	Address        string         `json:"address"`
	StatusDelivery deliveryStatus `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (d *DeliveryDomain) CalculateTotal() int64 {
	var total int64
	for _, item := range d.DeliveryItems {
		total += item.TotalAmount()
	}
	return total
}

type DeliveryToSave struct {
	Items     []DeliveryItemToSave `json:"items" validate:"required,min=1"`
	Recipient string               `json:"recipient" validate:"required,min=1"`
	Address   string               `json:"address" validate:"required,min=1"`
}

type DeliveryItemToSave struct {
	ProductId int64 `json:"product_id" validate:"gt=0"`
	ItemPrice int64 `json:"item_price" validate:"gt=0"`
	Quantity  int64 `json:"quantity" validate:"gt=0"`
}

type UpdateStatus struct {
	DeliveryId     int64          `json:"id" validate:"gt=0"`
	StatusDelivery deliveryStatus `json:"status" validate:"required,min=1"`
}

func (d *DeliveryToSave) MapToDomain() *DeliveryDomain {
	var items = make([]DeliveryItem, 0, len(d.Items))
	for _, item := range d.Items {
		items = append(items, DeliveryItem{
			ProductId: item.ProductId,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		})
	}
	return &DeliveryDomain{
		DeliveryItems:  items,
		Recipient:      d.Recipient,
		Address:        d.Address,
		StatusDelivery: CREATED,
		CreatedAt:      time.Now(),
	}
}

func (d *DeliveryDomain) MapToDeliveryProductsDB(productId int64) []DeliveryProductDB {
	deliveryProductsDB := make([]DeliveryProductDB, 0, len(d.DeliveryItems))
	for _, item := range d.DeliveryItems {
		deliveryProductsDB = append(deliveryProductsDB, DeliveryProductDB{
			DeliveryId:  productId,
			ProductId:   item.ProductId,
			ItemPrice:   item.ItemPrice,
			Quantity:    item.Quantity,
			TotalAmount: item.TotalAmount(),
		})
	}
	return deliveryProductsDB
}

func (d *DeliveryDomain) MapToDB() DeliveryDB {
	return DeliveryDB{
		Recipient:      d.Recipient,
		Address:        d.Address,
		StatusDelivery: d.StatusDelivery,
		CreatedAt:      d.CreatedAt,
	}
}

func (d *DeliveryDB) MapToDomain() DeliveryDomain {
	return DeliveryDomain{
		Recipient:      d.Recipient,
		Address:        d.Address,
		StatusDelivery: d.StatusDelivery,
		CreatedAt:      d.CreatedAt,
	}
}
