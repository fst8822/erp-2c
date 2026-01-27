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

type ItemsDB struct {
	ID         int64 `db:"id"`
	DeliveryID int64 `db:"delivery_id"`
	ProductID  int64 `db:"product_id"`
	ItemPrice  int64 `db:"item_price"`
	Quantity   int64 `db:"quantity"`
}

type DeliveryWithItemsDB struct {
	DeliveryDB      DeliveryDB
	DeliveryItemsDB []ItemsDB
}

type DeliverListDB struct {
	DeliveriesDB []DeliveryDB
	ItemsDB      []ItemsDB
}

type DeliverDomain struct {
	ID            int64          `json:"id"`
	Recipient     string         `json:"recipient"`
	Address       string         `json:"address"`
	Status        deliveryStatus `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
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
	Status     deliveryStatus `json:"status" validate:"required,min=1"`
}

func (d *DeliveryToSave) MapToDomain() DeliveryItemsDomain {
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
		},
		Items: items,
	}
	delivery.CalculateTotalAmount()
	return delivery
}

func (i *DeliveryItemsDomain) MapToDBWithItems() DeliveryWithItemsDB {
	var items = make([]ItemsDB, 0, len(i.Items))
	for _, item := range i.Items {
		items = append(items, ItemsDB{
			ProductID: item.ProductID,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		})
	}
	return DeliveryWithItemsDB{
		DeliveryDB: DeliveryDB{
			Recipient: i.Recipient,
			Address:   i.Address,
			Status:    i.Status,
			CreatedAt: i.CreatedAt,
		},
		DeliveryItemsDB: items,
	}
}

func (d *DeliveryWithItemsDB) MapToDomain() DeliveryItemsDomain {
	var items = make([]ItemDomain, 0, len(d.DeliveryItemsDB))
	for _, item := range d.DeliveryItemsDB {
		itemDomain := ItemDomain{
			DeliveryID: item.DeliveryID,
			ProductID:  item.ProductID,
			ItemPrice:  item.ItemPrice,
			Quantity:   item.Quantity,
		}
		itemDomain.totalAmount()
		items = append(items, itemDomain)
	}
	delivery := DeliveryItemsDomain{
		DeliverDomain: DeliverDomain{
			ID:        d.DeliveryDB.ID,
			Recipient: d.DeliveryDB.Recipient,
			Address:   d.DeliveryDB.Address,
			Status:    d.DeliveryDB.Status,
			CreatedAt: d.DeliveryDB.CreatedAt,
		},
		Items: items,
	}
	delivery.CalculateTotalAmount()
	return delivery
}
