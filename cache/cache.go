package cache

import (
	"context"
	"erp-2c/model"
	"log/slog"
	"sync"
	"time"
)

type Cache interface {
	Get(id int64) (any, bool)
	GetAll() ([]any, bool)
	Add(delivery model.DeliveryWithItemsDB)
	cleanTTL(ctx context.Context)
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

func (c *MapCache) GetAll() ([]any, bool) {
	var res []any
	c.data.Range(func(key any, value any) bool {
		res = append(res, value)
		return true
	})
	if len(res) == 0 {
		return res, true
	}
	return nil, false
}

func (c *MapCache) Get(id int64) (any, bool) {
	return c.data.Load(id)
}

func (c *MapCache) Add(deliveryWithItems model.DeliveryWithItemsDB) {
	c.data.Store(deliveryWithItems.DeliveryDB.ID, itemCache{
		any:       deliveryWithItems,
		ExpiresAt: time.Now().Add(c.tTL),
	})
}

func (c *MapCache) cleanTTL(ctx context.Context) {
	const op = "cache.cleanTTL"
	logger := slog.With("op", op)

	ticket := time.NewTicker(c.tTL)
	defer ticket.Stop()

	select {
	case <-ctx.Done():
		logger.Info("Stop clear ttl, context is done")
		return
	case <-ticket.C:
		logger.Info("Start clear ttl")
		c.data.Range(func(key, value any) bool {
			item := value.(itemCache)
			if time.Now().After(item.ExpiresAt) {
				c.data.Delete(key)
			}
			return true
		})
	default:
		logger.Info("skip clear iteration")
	}
}
