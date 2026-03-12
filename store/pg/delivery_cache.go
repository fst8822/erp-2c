package pg

import (
	"context"
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type cacheInt interface {
	Get(key int64) (any, bool)
	GetAll() []any
	Add(key int64, value any)
	DeleteByKey(key int64)
	DeleteByKeys(keys []int64)
}

type DeliveryCache struct {
	deliveryRepo *DeliveryRepository
	cache        cacheInt
}

func NewDeliveryCacheCache(ctx context.Context, deliveryRepo *DeliveryRepository, cache cacheInt,
) *DeliveryCache {
	repo := DeliveryCache{deliveryRepo: deliveryRepo, cache: cache}
	return &repo
}

func (c *DeliveryCache) SaveWithItems(
	tx *sqlx.Tx, deliveryWithItems model.DeliveryWithItemsDB) (*model.DeliveryDB, error) {

	DeliveryDB, err := c.deliveryRepo.SaveWithItems(tx, deliveryWithItems)
	if err != nil {
		return nil, err
	}
	deliveryWithItems.DeliveryDB.ID = DeliveryDB.ID
	c.cache.Add(DeliveryDB.ID, deliveryWithItems)
	return DeliveryDB, nil
}

func (c *DeliveryCache) GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (model.DeliveryWithItemsDB, error) {
	if v, ok := c.cache.Get(deliveryId); ok {
		if deliveryItems, ok2 := v.(model.DeliveryWithItemsDB); ok2 {
			return deliveryItems, nil
		}
	}
	deliveryWithItemsDB, err := c.deliveryRepo.GetWithItemsById(tx, deliveryId)
	if err != nil {
		return model.DeliveryWithItemsDB{}, err
	}
	c.cache.Add(deliveryId, deliveryWithItemsDB)

	return deliveryWithItemsDB, nil
}

func (c *DeliveryCache) GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error) {
	deliveries := make([]model.DeliveryDB, 100)
	itemsDB := make([]model.ItemsDB, 100)
	isInterrupt := false

	if all := c.cache.GetAll(); len(all) > 0 {
		for _, v := range all {
			it, ok2 := v.(model.DeliveryWithItemsDB)
			if !ok2 {
				isInterrupt = true
				break
			}
			deliveries = append(deliveries, it.DeliveryDB)
			itemsDB = append(itemsDB, it.DeliveryItemsDB...)
		}
		if !isInterrupt {
			return &model.DeliverListDB{
				DeliveriesDB: deliveries,
				ItemsDB:      itemsDB,
			}, nil
		}
	}

	deliverListDB, err := c.deliveryRepo.GetAll(tx)
	if err != nil {
		return nil, err
	}
	c.addList(deliverListDB)

	return deliverListDB, nil
}

func (c *DeliveryCache) GetAllWithItemsByStatus(tx *sqlx.Tx, status model.DeliveryStatus) (*model.DeliverListDB, error) {
	return c.deliveryRepo.GetAllWithItemsByStatus(tx, status)
}

func (c *DeliveryCache) UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error {
	err := c.deliveryRepo.UpdateById(tx, deliveryId, status)
	if err != nil {
		return err
	}
	c.cache.DeleteByKey(deliveryId)
	return nil
}

func (c *DeliveryCache) DeleteById(tx *sqlx.Tx, deliveryId int64) error {
	err := c.deliveryRepo.DeleteById(tx, deliveryId)
	if err != nil {
		return err
	}
	c.cache.DeleteByKey(deliveryId)
	return nil
}

func (c *DeliveryCache) UpdateStatusById(tx *sqlx.Tx, deliveryId int64, status model.DeliveryStatus) error {
	err := c.deliveryRepo.UpdateStatusById(tx, deliveryId, status)
	if err != nil {
		return err
	}
	c.cache.DeleteByKey(deliveryId)
	return nil
}

func (c *DeliveryCache) UpdateStatusByIds(tx *sqlx.Tx, groups map[model.DeliveryStatus][]int64) error {
	err := c.deliveryRepo.UpdateStatusByIds(tx, groups)
	if err != nil {
		return err
	}
	for _, keys := range groups {
		c.cache.DeleteByKeys(keys)
	}
	return nil
}

func (c *DeliveryCache) addList(deliverListDB *model.DeliverListDB) {
	mapItem := make(map[int64][]model.ItemsDB)
	for _, v := range deliverListDB.ItemsDB {
		mapItem[v.DeliveryID] = append(mapItem[v.DeliveryID], v)
	}

	for _, v := range deliverListDB.DeliveriesDB {
		deliveryWithItemsDB := model.DeliveryWithItemsDB{
			DeliveryDB:      v,
			DeliveryItemsDB: mapItem[v.ID],
		}
		c.cache.Add(v.ID, deliveryWithItemsDB)
	}
}
