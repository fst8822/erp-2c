package model

import (
	"api-gateway/lib/types"
	"time"
)

type StatusCount struct {
	Status string `db:"status"`
	Count  int    `db:"count"`
}

type DeliveryStatus string

func (s DeliveryStatus) IsValid() error {
	switch s {
	case CREATED, SHIPPED, DELIVERED, CANCELLED, ACCEPTED:
		return nil
	default:
		return types.NewAppErr(string(s), types.ErrUnknownStatus)
	}
}

const (
	CREATED   DeliveryStatus = "CREATED"
	SHIPPED   DeliveryStatus = "SHIPPED"
	DELIVERED DeliveryStatus = "DELIVERED"
	CANCELLED DeliveryStatus = "CANCELLED"
	ACCEPTED  DeliveryStatus = "ACCEPTED"
)

type DeliverDomain struct {
	ID            int64          `json:"id"`
	Recipient     string         `json:"recipient"`
	Address       string         `json:"address"`
	Status        DeliveryStatus `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UserID        int64          `json:"user_id"`
	DeliverAmount int64          `json:"deliver_amount"`
}

type ItemDomain struct {
	DeliveryID int64 `json:"-"`
	ProductID  int64 `json:"product_id"`
	ItemPrice  int64 `json:"item_price"`
	Quantity   int64 `json:"quantity"`
	ItemAmount int64 `json:"item_amount"`
}

func (i *ItemDomain) totalAmount() int64 {
	res := i.ItemPrice * i.Quantity
	i.ItemAmount = res
	return res
}

type DeliveryItemsDomain struct {
	DeliverDomain `json:"deliver_domain"`
	Items         []ItemDomain `json:"items"`
}

type DeliveryItemListDomain struct {
	DeliveryItemsDomain []DeliveryItemsDomain `json:"delivery_items_list"`
}

func (i *DeliveryItemsDomain) CalculateTotalAmount() {
	var total int64
	for _, item := range i.Items {
		total += item.totalAmount()
	}
	i.DeliverAmount = total
}

type DeliveryToSave struct {
	Recipient string               `json:"recipient" validate:"required,min=1"`
	Address   string               `json:"address" validate:"required,min=1"`
	Items     []DeliveryItemToSave `json:"items" validate:"required,min=1"`
}

type DeliveryItemToSave struct {
	ProductId int64 `json:"product_id" validate:"gt=0"`
	ItemPrice int64 `json:"item_price" validate:"gt=0"`
	Quantity  int64 `json:"quantity" validate:"gt=0"`
}

type UpdateStatus struct {
	DeliveryId int64          `json:"id" validate:"gt=0"`
	Status     DeliveryStatus `json:"status" validate:"required,min=1"`
}

func (d *DeliveryToSave) MapToDomain(UserID int64) DeliveryItemsDomain {
	var items = make([]ItemDomain, 0, len(d.Items))

	for _, item := range d.Items {
		itemDomain := ItemDomain{
			ProductID: item.ProductId,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		}
		itemDomain.totalAmount()
		items = append(items, itemDomain)
	}
	delivery := DeliveryItemsDomain{
		DeliverDomain: DeliverDomain{
			Recipient: d.Recipient,
			Address:   d.Address,
			Status:    CREATED,
			CreatedAt: time.Now(),
			UserID:    UserID,
		},
		Items: items,
	}
	delivery.CalculateTotalAmount()
	return delivery
}
