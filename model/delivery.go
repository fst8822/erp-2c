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
	ID        int64          `db:"id"`
	Recipient string         `db:"recipient"`
	Address   string         `db:"address"`
	Status    deliveryStatus `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
}

type DeliveryItemsDB struct {
	ID         int64 `db:"id"`
	DeliveryID int64 `db:"delivery_id"`
	ProductID  int64 `db:"product_id"`
	ItemPrice  int64 `db:"item_price"`
	Quantity   int64 `db:"quantity"`
}

type DeliveryWithItems struct {
	DeliveryDB      DeliveryDB
	DeliveryItemsDB []DeliveryItemsDB
}

type DeliveryItemDomain struct {
	ProductId int64 `json:"product_id"`
	ItemPrice int64 `json:"item_price"`
	Quantity  int64 `json:"quantity"`
}

func (i *DeliveryItemDomain) totalAmount() int64 {
	return i.ItemPrice * i.Quantity
}

type DeliveryDomain struct {
	ID        int64                `json:"id"`
	Recipient string               `json:"recipient"`
	Address   string               `json:"address"`
	Status    deliveryStatus       `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	Items     []DeliveryItemDomain `json:"delivery_items"`
}

func (d *DeliveryDomain) CalculateTotal() int64 {
	var total int64
	for _, item := range d.Items {
		total += item.totalAmount()
	}
	return total
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
	Status     deliveryStatus `json:"status" validate:"required,min=1"`
}

func (d *DeliveryToSave) MapToDomain() DeliveryDomain {
	var items = make([]DeliveryItemDomain, 0, len(d.Items))
	for _, item := range d.Items {
		items = append(items, DeliveryItemDomain{
			ProductId: item.ProductId,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		})
	}
	return DeliveryDomain{
		Items:     items,
		Recipient: d.Recipient,
		Address:   d.Address,
		Status:    CREATED,
		CreatedAt: time.Now(),
	}
}

func (d *DeliveryDomain) MapToDB() DeliveryDB {
	return DeliveryDB{
		Recipient: d.Recipient,
		Address:   d.Address,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
	}
}
func (d *DeliveryDomain) MapToDBWithItems() DeliveryWithItems {
	var items = make([]DeliveryItemsDB, 0, len(d.Items))
	for _, item := range d.Items {
		items = append(items, DeliveryItemsDB{
			ProductID: item.ProductId,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		})
	}
	return DeliveryWithItems{
		DeliveryDB: DeliveryDB{
			Recipient: d.Recipient,
			Address:   d.Address,
			Status:    d.Status,
			CreatedAt: d.CreatedAt,
		},
		DeliveryItemsDB: items,
	}
}

func (d *DeliveryDB) MapToDomain() DeliveryDomain {
	return DeliveryDomain{
		Recipient: d.Recipient,
		Address:   d.Address,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
	}
}

func (d *DeliveryWithItems) MapToDomain() DeliveryDomain {
	var items = make([]DeliveryItemDomain, 0, len(d.DeliveryItemsDB))
	for _, item := range d.DeliveryItemsDB {
		items = append(items, DeliveryItemDomain{
			ProductId: item.ProductID,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		})
	}
	return DeliveryDomain{
		ID:        d.DeliveryDB.ID,
		Recipient: d.DeliveryDB.Recipient,
		Address:   d.DeliveryDB.Address,
		Status:    d.DeliveryDB.Status,
		CreatedAt: d.DeliveryDB.CreatedAt,
		Items:     items,
	}
}
