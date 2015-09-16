package container

import (
	"errors"
	"math/rand"
)

type SortedSet struct {
	set *Set
	sl  *skipList
}

func NewSortedSet(cf func(obj1, obj2 interface{}) int, tc func(d interface{}) bool) *SortedSet {
	if cf == nil {
		return nil
	}
	return &SortedSet{
		set: NewSet(tc),
		sl:  newSkiplist(cf),
	}
}

func (ss *SortedSet) Add(d interface{}) (bool, error) {
	if ss.set.Add(d) {
		if err := ss.sl.insert(d); err != nil {
			ss.set.Delete(d)
			return false, err
		} else {
			return true, nil
		}
	} else {
		return false, errors.New("duplicated element or type error.")
	}
}

func (ss *SortedSet) Delete(d interface{}) (bool, error) {
	if ss.set.Delete(d) {
		if err := ss.sl.delete(d); err != nil {
			ss.set.Add(d)
			return false, err
		} else {
			return true, nil
		}
	} else {
		return false, errors.New("element doesn't exist or type error.")
	}
}

func (ss *SortedSet) PopByRank(rank int) interface{} {
	obj := ss.sl.deleteByRank(rank)
	if obj != nil {
		ss.set.DeleteRaw(obj)
	}
	return obj
}

/*随机取出一个元素*/
func (ss *SortedSet) PopRandom() interface{} {
	return ss.PopByRank(rand.Int() % ss.Length())
}

func (ss *SortedSet) PopByRanges(start, end int) ([]interface{}, error) {
	objs, err := ss.sl.deleteRangeByRank(start, end)
	for obj := range objs {
		ss.set.DeleteRaw(obj)
	}
	return objs, err
}

func (ss *SortedSet) GetTopN(n int) ([]interface{}, error) {
	return ss.sl.getTopN(n)
}

func (ss *SortedSet) Contain(d interface{}) bool {
	return ss.set.Contain(d)
}

func (ss *SortedSet) Length() int {
	return ss.set.Length()
}

func (ss *SortedSet) Flush() []interface{} {
	ds := ss.set.Flush()
	ss.sl.flush()
	return ds
}

/*返回集合的有序列表*/
func (ss *SortedSet) List() []interface{} {
	return ss.sl.list()
}

func (ss *SortedSet) Empty() bool {
	return ss.set.Empty()
}
