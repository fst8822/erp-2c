package cache

import (
	"context"
	"erp-2c/model"
	"erp-2c/store"

	"github.com/jmoiron/sqlx"
)

type RepositoryCache struct {
	store    store.Store
	mapCache Cache
	done     chan struct{}
}

func NewRepositoryCache(
	ctx context.Context,
	store *store.Store,
	mapCache Cache,
	done chan struct{}) *RepositoryCache {
	repo := RepositoryCache{store: *store, mapCache: mapCache, done: done}

	go mapCache.cleanTTL(ctx, repo.done)
	return &repo
}

func (c *RepositoryCache) GetAll(ctx context.Context, tx *sqlx.Tx) (*model.DeliverListDB, error) {
	deliveries := make([]model.DeliveryDB, 100)
	itemsDB := make([]model.ItemsDB, 100)

	if vAny, ok := c.mapCache.GetAll(); ok {
		for _, v := range vAny {
			if itemCaches, ok := v.(itemCache); ok {
				if cached, ok := itemCaches.any.(model.DeliveryWithItemsDB); ok {
					deliveries = append(deliveries, cached.DeliveryDB)
					itemsDB = append(itemsDB, cached.DeliveryItemsDB...)
				}
			}
		}

		return &model.DeliverListDB{
			DeliveriesDB: deliveries,
			ItemsDB:      itemsDB,
		}, nil
	}

	deliverListDB, err := c.store.Delivery.GetAll(tx)
	if err != nil {
		return nil, err
	}
	c.addList(deliverListDB)

	return deliverListDB, nil
}

func (c *RepositoryCache) GetById(ctx context.Context, tx *sqlx.Tx, id int64) (model.DeliveryWithItemsDB, error) {
	if v, ok := c.mapCache.Get(id); ok {
		if v2, ok2 := v.(itemCache); ok2 {
			if cached, ok3 := v2.any.(model.DeliveryWithItemsDB); ok3 {
				return cached, nil
			}
		}
	}
	deliveryWithItemsDB, err := c.store.Delivery.GetWithItemsById(tx, id)
	if err != nil {
		return model.DeliveryWithItemsDB{}, err
	}
	c.Add(deliveryWithItemsDB)

	return deliveryWithItemsDB, nil
}

func (c *RepositoryCache) Add(delivery model.DeliveryWithItemsDB) {
	c.mapCache.Add(delivery)
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
		c.Add(deliveryWithItemsDB)
	}
}
