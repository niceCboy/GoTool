package tool

import (
	"sync"
	"sync/atomic"
)

type Counter struct {
	sync.RWMutex
	v int64
}

func NewCounter() *Counter {
	return &Counter{
		v: 0,
	}
}

func (c *Counter) Incr(num int64) {
	c.Lock()
	defer c.Unlock()
	if num < 0 {
		num = -num
	}
	atomic.AddInt64(&c.v, num)
}

func (c *Counter) Decr(num int64) {
	c.Lock()
	defer c.Unlock()
	if num > 0 {
		num = -num
	}
	atomic.AddInt64(&c.v, num)
}

func (c *Counter) ReadV() int64 {
	c.RLock()
	defer c.RUnlock()
	return atomic.LoadInt64(&c.v)
}

/*计数器自减一，并判断自减后是否与给定值相等*/
func (c *Counter) DecrToEqual(target int64) bool {
	c.Lock()
	defer c.Unlock()
	atomic.AddInt64(&c.v, -1)
	if atomic.LoadInt64(&c.v) == target {
		return true
	} else {
		return false
	}
}

/*判断计数器是否等于给定值,相等则自加且返回成功*/
func (c *Counter) EqualAndIncr(target int64) bool {
	c.Lock()
	defer c.Unlock()
	if atomic.LoadInt64(&c.v) == target {
		c.v++
		return true
	} else {
		return false
	}
}
