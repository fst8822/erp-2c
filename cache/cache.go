package cache

import (
	"erp-2c/model"
	"sync"
	"time"
)

type Cache interface {
	Get(delivery model.DeliveryDB) (any, bool)
	Add(delivery model.DeliveryDB)
	cleanTTL()
}

type itemCache struct {
	any
	ExpiresAt time.Time
}

type MapCache struct {
	data sync.Map
	tTL  time.Duration
}

func NewMapCache(TTL time.Duration) *MapCache {
	return &MapCache{tTL: TTL}
}

func (c *MapCache) Get(delivery model.DeliveryDB) (any, bool) {
	return c.data.Load(delivery.ID)

}

func (c *MapCache) Add(delivery model.DeliveryDB) {

	c.data.Store(delivery.ID, itemCache{
		any:       delivery,
		ExpiresAt: time.Now().Add(c.tTL),
	})
}

func (c *MapCache) cleanTTL() {
	c.data.Range(func(key, value any) bool {
		item := value.(itemCache)
		if time.Now().After(item.ExpiresAt) {
			c.data.Delete(key)
		}
		return true
	})
}
