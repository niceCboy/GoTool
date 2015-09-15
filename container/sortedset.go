package container

import (
	"errors"
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

func (ss *SortedSet) List() []interface{} {
	return ss.sl.list()
}

func (ss *SortedSet) Empty() bool {
	return ss.set.Empty()
}
