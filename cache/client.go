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
	tTL  time.Duration
}

func NewMapCache(tTL time.Duration) MapCache {
	return MapCache{tTL: tTL}
}

func (c *MapCache) GetAll() ([]any, bool) {
	res := make([]any, 0, 100)

	c.data.Range(func(key any, value any) bool {
		if ic, ok := value.(itemCache); ok {
			res = append(res, ic.value)
			return true
		}
		return false
	})
	if len(res) == 0 {
		return res, false
	}
	return res, true
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
		ExpiresAt: time.Now().Add(c.tTL),
	})
}

func (c *MapCache) CleanTTL(ctx context.Context) {
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
