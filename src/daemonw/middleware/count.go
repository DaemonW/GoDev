package middleware

import (
	"github.com/go-redis/redis"
	"time"
	"sync"
)

var (
	globolCounter = newMemoryCounter()
)

type memoryCounter struct {
	mu      sync.Mutex
	counter map[string]int64
}

func newMemoryCounter() *memoryCounter {
	c := &memoryCounter{}
	c.counter = make(map[string]int64)
	return c
}

func (c *memoryCounter) Incr(key string) {
	c.mu.Lock()
	if count, ok := c.counter[key]; ok {
		c.counter[key] = count + 1
	}else{
		c.counter[key] = 1
	}
	c.mu.Unlock()
}

func (c *memoryCounter) Count(key string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if count, ok := c.counter[key]; ok {
		return count
	}
	return 0
}

type Counter struct {
	rc       *redis.Client
	maxCount int64
	countAt  time.Time
	deadLine time.Time
	fallBack *memoryCounter
}

func NewCounter(rc *redis.Client, max int64, deadLine time.Time) *Counter {
	c := &Counter{rc: rc, maxCount: max, countAt: time.Now(), deadLine: deadLine}
	c.fallBack = globolCounter
	return c
}

func (c *Counter) Allow(key string, deadLine time.Time) (n int64, allow bool) {
	if deadLine.After(c.deadLine) {
		return 0, true
	}
	c.fallBack.Incr(key)
	count := c.fallBack.Count(key)
	return count, count < c.maxCount
}

func (c *Counter) DeadLine(key string) time.Time {
	return c.deadLine
}
