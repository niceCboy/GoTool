package detector

import (
	"sync"
	"time"
)

var (
	Default *TimeOutDetector
)

func InitDefault() {
	Default = NewTimeOutDetector()
}

type TimeOutDetector struct {
	tSet map[string]*time.Timer
	tRec map[string]time.Duration
	sync.RWMutex
}

func NewTimeOutDetector() *TimeOutDetector {
	return &TimeOutDetector{
		tSet: make(map[string]*time.Timer),
		tRec: make(map[string]time.Duration),
	}
}

func (tod *TimeOutDetector) AddActive(name string, period time.Duration, call func()) bool {
	tod.Lock()
	defer tod.Unlock()
	if _, ok := tod.tSet[name]; !ok {
		tod.tSet[name] = time.AfterFunc(period, func() {
			call()
		})
		tod.tRec[name] = period
		return true
	} else {
		return false
	}
}

/*最低限度添加功能实现，不保证goroutine安全，不保证定时器名的唯一性，不保证超时清理定时器集合相应名字的对象*/
func (tod *TimeOutDetector) AddActiveRaw(name string, period time.Duration, call func()) {
	tod.tSet[name] = time.AfterFunc(period, func() {
		call()
	})
	tod.tRec[name] = period
}

func (tod *TimeOutDetector) DelActive(name string) bool {
	tod.Lock()
	defer tod.Unlock()
	if _, ok := tod.tSet[name]; ok {
		tod.tSet[name].Stop()
		delete(tod.tSet, name)
		delete(tod.tRec, name)
		return true
	} else {
		return false
	}
}

/*最低限度删除功能实现，不保证goroutine安全，不保证定时器名的存在*/
func (tod *TimeOutDetector) DelActiveRaw(name string) {
	tod.tSet[name].Stop()
	delete(tod.tSet, name)
	delete(tod.tRec, name)
}

func (tod *TimeOutDetector) Reset(name string) bool {
	tod.Lock()
	defer tod.Unlock()
	if _, ok := tod.tSet[name]; ok {
		tod.tSet[name].Reset(tod.tRec[name])
		return true
	} else {
		return false
	}
}

/*最低限度重置时间功能实现，不保证goroutine安全，不保证定时器名的存在*/
func (tod *TimeOutDetector) ResetRaw(name string) {
	tod.tSet[name].Reset(tod.tRec[name])
}

func (tod *TimeOutDetector) ReSchedule(name string, period time.Duration, call func()) bool {
	tod.Lock()
	defer tod.Unlock()
	if _, ok := tod.tSet[name]; ok {
		tod.tSet[name].Stop()
		tod.tSet[name] = time.AfterFunc(period, func() {
			call()
		})
		tod.tRec[name] = period
		return true
	} else {
		return false
	}
}

/*最低限度重定时器功能实现，不保证goroutine安全，不保证定时器名的存在,不保证超时清理定时器集合相应名字的对象*/
func (tod *TimeOutDetector) ReScheduleRaw(name string, period time.Duration, call func()) {
	tod.tSet[name].Stop()
	tod.tSet[name] = time.AfterFunc(period, func() {
		call()
	})
	tod.tRec[name] = period
}

func (tod *TimeOutDetector) Contain(name string) bool {
	tod.Lock()
	defer tod.Unlock()
	_, ok := tod.tSet[name]
	return ok
}

func (tod *TimeOutDetector) Clear() {
	tod.Lock()
	defer tod.Unlock()
	for name, timer := range tod.tSet {
		timer.Stop()
		delete(tod.tSet, name)
		delete(tod.tRec, name)
	}
}
