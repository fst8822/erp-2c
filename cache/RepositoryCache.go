package cache

import (
	"context"
	"erp-2c/model"
)

type RepositoryCache struct {
	MapCache Cache
}

func NewRepositoryCache(mapCache Cache) *RepositoryCache {
	mapCache.cleanTTL()
	return &RepositoryCache{MapCache: mapCache}
}

func (c *RepositoryCache) GetAll(ctx context.Context) []model.DeliveryDB {
	return make([]model.DeliveryDB, 0)
}

func (c *RepositoryCache) Get(ctx context.Context) model.DeliveryDB {
	return model.DeliveryDB{}
}

func (c *RepositoryCache) Add(ctx context.Context) {

}
