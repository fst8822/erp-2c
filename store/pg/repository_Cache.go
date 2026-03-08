package pg

import (
	"context"
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type cacheInt interface {
	Get(key int64) (any, bool)
	GetAll() ([]any, bool)
	Add(key int64, value any)
	CleanTTL(ctx context.Context)
	DeleteByKey(key int64)
	DeleteByKeys(keys []int64)
}

//type deliveryRepositoryInt interface {
//	GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (model.DeliveryWithItemsDB, error)
//	GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error)
//	SaveWithItems(tx *sqlx.Tx, deliveryWithItems model.DeliveryWithItemsDB) (*model.DeliveryDB, error)
//}

type RepositoryCache struct {
	deliveryRepo *DeliveryRepository
	cache        cacheInt
}

func NewRepositoryCache(ctx context.Context, deliveryRepo *DeliveryRepository, cache cacheInt,
) *RepositoryCache {
	repo := RepositoryCache{deliveryRepo: deliveryRepo, cache: cache}
	go cache.CleanTTL(ctx)
	return &repo
}

func (c *RepositoryCache) GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error) {
	deliveries := make([]model.DeliveryDB, 100)
	itemsDB := make([]model.ItemsDB, 100)
	isInterrupt := false

	if all, ok := c.cache.GetAll(); ok {
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

func (c *RepositoryCache) GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (model.DeliveryWithItemsDB, error) {
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

func (c *RepositoryCache) SaveWithItems(
	tx *sqlx.Tx, deliveryWithItems model.DeliveryWithItemsDB) (*model.DeliveryDB, error) {

	DeliveryDB, err := c.deliveryRepo.SaveWithItems(tx, deliveryWithItems)
	if err != nil {
		return nil, err
	}
	deliveryWithItems.DeliveryDB.ID = DeliveryDB.ID
	c.cache.Add(DeliveryDB.ID, deliveryWithItems)
	return DeliveryDB, nil
}

func (c *RepositoryCache) GetAllWithItemsByStatus(tx *sqlx.Tx, status model.DeliveryStatus) (*model.DeliverListDB, error) {
	return c.deliveryRepo.GetAllWithItemsByStatus(tx, status)
}

func (c *RepositoryCache) LockAndGetDeliveries(status model.DeliveryStatus) ([]model.DeliveryDB, error) {
	return c.deliveryRepo.LockAndGetDeliveries(status)
}

func (c *RepositoryCache) UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error {
	c.cache.DeleteByKey(deliveryId)
	return c.deliveryRepo.UpdateById(tx, deliveryId, status)
}

func (c *RepositoryCache) DeleteById(tx *sqlx.Tx, deliveryId int64) error {
	c.cache.DeleteByKey(deliveryId)
	return c.deliveryRepo.DeleteById(tx, deliveryId)
}

func (c *RepositoryCache) UpdateStatusById(tx *sqlx.Tx, deliveryId int64, status model.DeliveryStatus) error {
	c.cache.DeleteByKey(deliveryId)
	return c.deliveryRepo.UpdateStatusById(tx, deliveryId, status)
}

func (c *RepositoryCache) UpdateStatusByIds(tx *sqlx.Tx, groups map[model.DeliveryStatus][]int64) error {
	for _, keys := range groups {
		c.cache.DeleteByKeys(keys)
	}
	return c.deliveryRepo.UpdateStatusByIds(tx, groups)
}

func (c *RepositoryCache) GetStatusCount(tx *sqlx.Tx) ([]model.StatusCount, error) {
	return c.deliveryRepo.GetStatusCount(tx)
}

func (c *RepositoryCache) BeginTxx(ctx context.Context) (*sqlx.Tx, error) {
	return c.deliveryRepo.BeginTxx(ctx)
}

func (c *RepositoryCache) addList(deliverListDB *model.DeliverListDB) {
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
