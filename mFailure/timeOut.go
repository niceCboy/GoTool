package mFailure

import (
	"sync"
	"time"
)

type TimeOutScheduler struct {
	tSet map[string]*time.Timer
	tRec map[string]time.Duration
	sync.RWMutex
}

func NewTimeOutScheduler() *TimeOutScheduler {
	return &TimeOutScheduler{
		tSet: make(map[string]*time.Timer),
		tRec: make(map[string]time.Duration),
	}
}

func (tos *TimeOutScheduler) AddActive(name string, period time.Duration, call func()) bool {
	tos.Lock()
	defer tos.Unlock()
	if _, ok := tos.tSet[name]; !ok {
		tos.tSet[name] = time.AfterFunc(period, func() {
			tos.DelActive(name)
			call()
		})
		tos.tRec[name] = period
		return true
	} else {
		return false
	}
}

/*最低限度添加功能实现，不保证goroutine安全，不保证定时器名的唯一性，不保证超时清理定时器集合相应名字的对象*/
func (tos *TimeOutScheduler) AddActiveRaw(name string, period time.Duration, call func()) {
	tos.tSet[name] = time.AfterFunc(period, func() {
		call()
	})
	tos.tRec[name] = period
}

func (tos *TimeOutScheduler) DelActive(name string) bool {
	tos.Lock()
	defer tos.Unlock()
	if _, ok := tos.tSet[name]; ok {
		tos.tSet[name].Stop()
		delete(tos.tSet, name)
		delete(tos.tRec, name)
		return true
	} else {
		return false
	}
}

/*最低限度删除功能实现，不保证goroutine安全，不保证定时器名的存在*/
func (tos *TimeOutScheduler) DelActiveRaw(name string) {
	tos.tSet[name].Stop()
	delete(tos.tSet, name)
	delete(tos.tRec, name)
}

func (tos *TimeOutScheduler) Reset(name string) bool {
	tos.Lock()
	defer tos.Unlock()
	if _, ok := tos.tSet[name]; ok {
		tos.tSet[name].Reset(tos.tRec[name])
		return true
	} else {
		return false
	}
}

/*最低限度重置时间功能实现，不保证goroutine安全，不保证定时器名的存在*/
func (tos *TimeOutScheduler) ResetRaw(name string) {
	tos.tSet[name].Reset(tos.tRec[name])
}

func (tos *TimeOutScheduler) ReSchedule(name string, period time.Duration, call func()) bool {
	tos.Lock()
	defer tos.Unlock()
	if _, ok := tos.tSet[name]; ok {
		tos.tSet[name].Stop()
		tos.tSet[name] = time.AfterFunc(period, func() {
			tos.DelActive(name)
			call()
		})
		tos.tRec[name] = period
		return true
	} else {
		return false
	}
}

/*最低限度重定时器功能实现，不保证goroutine安全，不保证定时器名的存在,不保证超时清理定时器集合相应名字的对象*/
func (tos *TimeOutScheduler) ReScheduleRaw(name string, period time.Duration, call func()) {
	tos.tSet[name].Stop()
	tos.tSet[name] = time.AfterFunc(period, func() {
		call()
	})
	tos.tRec[name] = period
}

func (tos *TimeOutScheduler) Clear() {
	tos.Lock()
	defer tos.Unlock()
	for name, timer := range tos.tSet {
		timer.Stop()
		delete(tos.tSet, name)
		delete(tos.tRec, name)
	}
}
