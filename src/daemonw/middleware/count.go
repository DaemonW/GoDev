package middleware

import (
	"time"
	"sync"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"strconv"
	"daemonw/log"
)

var (
	fallbackCounter = newMemoryCounter()
)

type countItem struct {
	count int64
	limit int64
	begin time.Time
	end   time.Time
}

type memoryCounter struct {
	mu      sync.RWMutex
	counter map[string]countItem
}

func newMemoryCounter() *memoryCounter {
	c := &memoryCounter{}
	c.counter = make(map[string]countItem)
	return c
}

func (c *memoryCounter) setLimitTime(key string, limit int64, begin, end time.Time) {
	c.mu.Lock()
	if item, ok := c.counter[key]; ok {
		item.begin = begin
		item.end = end
		item.limit = limit
	} else {
		c.counter[key] = countItem{begin: begin, end: end, limit: limit}
	}
	c.mu.Unlock()
}

func (c *memoryCounter) Incr(key string) {
	c.mu.Lock()
	if item, ok := c.counter[key]; ok {
		item.count++
	}
	c.mu.Unlock()
}

func (c *memoryCounter) Count(key string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if item, ok := c.counter[key]; ok {
		return item.count
	}
	return 0
}

func (c *memoryCounter) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.counter, key)
}

type Counter struct {
	redis redisPiper
	limit int64
	count int64
	begin time.Time
	end   time.Time
}

func NewCounter(piper redisPiper) *Counter {
	c := &Counter{redis: piper}
	return c
}

func (c *Counter) DeleleLimit(key string) error {
	_, err := c.redis.Pipelined(func(pipe redis.Pipeliner) error {
		pipe.Del(key)
		_, e := pipe.Exec()
		return e
	})
	fallbackCounter.Delete(key)
	return err
}

func (c *Counter) SetLimit(key string, limit int64, begin, end time.Time) {
	c.limit = limit
	_, err := c.redis.Pipelined(func(pipe redis.Pipeliner) error {
		pipe.HMSet(key, map[string]interface{}{
			"begin": begin.Unix(),
			"end":   end.Unix(),
			"limit": limit,
			"count": 0,
		})
		//pipe.ExpireAt(key, end)
		_, e := pipe.Exec()
		return e
	})
	if err != nil {
		log.Error().Err(err).Msg("set limit failed by redis")
		fallbackCounter.setLimitTime(key, limit, begin, end)
	}
}

func (c *Counter) Allow(key string) (n int64, allow bool) {
	var varsCmd *redis.SliceCmd
	var countCmd *redis.IntCmd
	_, err := c.redis.Pipelined(func(pipe redis.Pipeliner) error {
		varsCmd = pipe.HMGet(key, "begin", "end", "limit", "count")
		countCmd = pipe.HIncrBy(key, "count", 1)
		_, e := pipe.Exec()
		return e
	})
	if err == nil {
		c.count, _ = countCmd.Result()
		vars, _ := varsCmd.Result()
		beginUnix, _ := strconv.ParseInt(vars[0].(string), 10, 64)
		c.begin = time.Unix(beginUnix, 0)
		endUnix, _ := strconv.ParseInt(vars[1].(string), 10, 64)
		c.end = time.Unix(endUnix, 0)
		c.limit, _ = strconv.ParseInt(vars[2].(string), 10, 64)
	} else {
		if item, ok := fallbackCounter.counter[key]; ok {
			item.count++
			c.begin = item.begin
			c.end = item.end
			c.limit = item.limit
			c.count = item.count
		}
	}
	if c.limit == 0 {
		return 0, true
	}
	current := time.Now()
	if current.After(c.end) || current.Before(c.begin) {
		return 0, true
	}

	return c.count, c.count <= c.limit
}

func (c *Counter) DeadLine(key string) (t *time.Time, err error) {
	var end = int64(0)
	var endCmd *redis.StringCmd
	_, err = c.redis.Pipelined(func(pipe redis.Pipeliner) error {
		endCmd = pipe.HGet(key, "end")
		_, e := pipe.Exec()
		return e
	})
	end, _ = endCmd.Int64()
	if err == nil {
		t := time.Unix(end, 0)
		return &t, nil
	}
	if item, ok := fallbackCounter.counter[key]; ok {
		return &item.end, nil
	}
	return nil, errors.New("not exist")
}
