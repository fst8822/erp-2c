package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type itemCache struct {
	value     any
	ExpiresAt time.Time
}

type MapCache struct {
	data sync.Map
	ttl  time.Duration
}

func NewMapCache(ctx context.Context, ttl time.Duration) *MapCache {
	m := &MapCache{ttl: ttl}
	go m.cleanTTL(ctx)
	return m
}

func (c *MapCache) Get(key int64) (any, bool) {
	vAny, ok := c.data.Load(key)
	var zero any
	if ok {
		if ic, ok2 := vAny.(itemCache); ok2 {
			return ic.value, true
		}
	}
	return zero, false
}

func (c *MapCache) DeleteByKey(key int64) {
	c.data.Delete(key)
}

func (c *MapCache) DeleteByKeys(keys []int64) {
	for _, key := range keys {
		c.data.Delete(key)
	}
}

func (c *MapCache) Add(key int64, value any) {
	c.data.Store(key, itemCache{
		value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	})
}

func (c *MapCache) cleanTTL(ctx context.Context) {
	const op = "cache.cleanTTL"
	logger := slog.With("op", op)

	ticket := time.NewTicker(c.ttl)
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
