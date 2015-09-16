package mContainer

import ()

type Set struct {
	data      map[interface{}]interface{}
	typeCheck func(d interface{}) bool
}

func NewSet(tc func(interface{}) bool) *Set {
	return &Set{
		data:      make(map[interface{}]interface{}),
		typeCheck: tc,
	}
}

func (s *Set) Add(d interface{}) bool {
	if s.typeCheck != nil && !withRecover(s.typeCheck, d) || s.Contain(d) {
		return false
	}
	s.data[d] = nil
	return true
}

/*return the numbers been added*/
func (s *Set) MultiAdd(ds []interface{}) int {
	if s.typeCheck != nil {
		for d := range ds {
			if !withRecover(s.typeCheck, d) {
				return 0
			}
		}
	}
	count := 0
	for d := range ds {
		if !s.Contain(d) {
			s.data[d] = nil
			count++
		}
	}
	return count
}

func (s *Set) Delete(d interface{}) bool {
	if s.typeCheck != nil && !withRecover(s.typeCheck, d) || !s.Contain(d) {
		return false
	}
	s.DeleteRaw(d)
	return true
}

/*不对类型和存在进行检查，提高效率*/
func (s *Set) DeleteRaw(d interface{}) {
	delete(s.data, d)
}

func (s *Set) Contain(d interface{}) bool {
	_, ok := s.data[d]
	return ok
}

func (s *Set) Length() int {
	return len(s.data)
}

func (s *Set) Flush() []interface{} {
	rtd := []interface{}{}
	for d, _ := range s.data {
		rtd = append(rtd, d)
	}
	s.data = make(map[interface{}]interface{})
	return rtd
}

func (s *Set) List() []interface{} {
	rtd := []interface{}{}
	for d, _ := range s.data {
		rtd = append(rtd, d)
	}
	return rtd
}

func (s *Set) Empty() bool {
	if len(s.data) == 0 {
		return true
	}
	return false
}

func withRecover(f func(interface{}) bool, d interface{}) (rtv bool) {
	defer func() {
		if r := recover(); r != nil {
			rtv = false
		}
	}()
	rtv = f(d)
	return
}
