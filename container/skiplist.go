package container

import ()

/*
  it's not concurrent safe and not ensure unique
*/
type skipList struct {
	header *skiplistNode
	tail   *skiplistNode
	update []*skiplistNode
	rank   []int
	length int
	level  int

	//两个对象进行比较 ，相等返回0，obj1>obj2返回一个大于0的数，obj1<obj2返回一个小于0的数,根据函数计算值由大到小排序
	compare func(obj1 interface{}, obj2 interface{}) int
}

func newSkiplist(cf func(obj1, obj2 interface{}) int) *skipList {
	return &skipList{
		header:  newSkiplistNode(MAXLEVEL, nil, 0),
		tail:    nil,
		update:  make([]*skiplistNode, MAXLEVEL),
		rank:    make([]int, MAXLEVEL),
		length:  0,
		level:   1,
		compare: cf,
	}
}

func (s *skipList) flush() {
	s.header = newSkiplistNode(MAXLEVEL, nil, 0)
	s.tail = nil
	s.update = make([]*skiplistNode, MAXLEVEL)
	s.rank = make([]int, MAXLEVEL)
	s.length = 0
	s.level = 1
}

func (s *skipList) len() int {
	return s.length
}

func (s *skipList) insert(obj interface{}) error {
	//更新插入前缀节点列表及其rank
	n := s.header
	for i := s.level - 1; i >= 0; i-- {
		if i == s.level-1 {
			s.rank[i] = 0
		} else {
			s.rank[i] = s.rank[i+1]
		}
		var (
			rtv int
			e   error
		)
		for n.levels[i].forward != nil {
			if rtv, e = compareWithRecover(s.compare, n.levels[i].forward.obj, obj); e == nil && rtv > 0 {
				s.rank[i] += n.levels[i].span
				n = n.levels[i].forward
			} else if e != nil {
				return e
			} else if rtv <= 0 {
				break
			}
		}
		s.update[i] = n
	}

	level := randomLevel()
	if level > s.level { //增加层数
		for l := s.level; l < level; l++ {
			s.rank[l] = 0
			s.update[l] = s.header
			s.update[l].levels[l].span = s.length
		}
		s.level = level
	}

	n = newSkiplistNode(level, obj)
	for l := 0; l < level; l++ {
		//更新l层的链表
		n.levels[l].forward = s.update[l].level[l].forward
		s.update[l].level[i].forward = n
		//更新l层的跨度记录
		n.levels[l].span = s.update[l].level[l].span - (s.rank[0] - s.rank[l])
		s.update[l] = s.rank[0] - s.rank[l] + 1
	}

	//增加未达到层数的跨度记录
	for l := level; l < s.level; l++ {
		s.update[l].level[l].span++
	}

	//更新最底层链前缀元素
	if s.update[0] == s.header {
		n.backward = nil
	} else {
		n.backward = s.update[0]
	}
	if n.level[0].forward != nil {
		n.level[0].forward.backward = n
	} else {
		s.tail = n
	}
	s.length++
	return nil
}

func (s *skipList) delete(obj interface{}) error {
	node, err := findPreNodes(obj)
	if err != nil {
		return err
	}

	if node != nil && s.compare(node.obj, obj) == 0 {
		for l := 0; l < s.level; l++ {
			if s.update[l].levels[l].forward == node {
				s.update[l].levels[l].span += node.levels[l].span - 1
				s.update[l].levels[l].forward = node.levels[l].forward
			} else {
				s.update[l].levels.span -= 1
			}
		}
	} else {
		return nil //未找到
	}

	//更新前缀节点
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		s.tail = node.backward
	}

	//清空最上层
	for s.level > 1 && s.header.level[s.level-1].forward == nil {
		s.level--
	}

	//减小长度
	s.length--
	return nil
}

/*找到各层最后一个compare计算不大于0的node,记录在update中，返回第一层的update的下一节点*/
func (s *skipList) findPreNodes(obj interface{}) (node *skiplistNode, err error) {
	n := s.header
	start := n
	for i := s.level - 1; i >= 0; i-- {
		flaged := false
		for n.levels[i].forward != nil {
			if rtv, e := compareWithRecover(s.compare, n.levels[i].forward.obj, obj); e == nil && rtv > 0 {
				n = n.levels[i].forward
			} else if e != nil {
				err = e
				return
			} else if rtv == 0 && n.levels[i].forward.obj != obj {
				if !flaged {
					start = n //标记下一层起始搜索节点
					flaged = true
				}
				n = n.levels[i].forward
			} else { // rtv < 0  or ==0 and forward.obj == obj
				s.update[i] = n
				if flaged {
					n = start
				}
				break
			}
		}
	}
	return s.update[0].forward, nil
}

func (s *skipList) contain(obj interface{}) (bool, error) {
	node, err := findPreNodes(obj)
	if err != nil {
		return false, err
	}
	if node != nil && s.compare(node.obj, obj) == 0 {
		return true, nil
	}
	return false, nil
}

/*返回有序的元素列表*/
func (s *skipList) list() []interface{} {
	ds := []interface{}{}
	node := s.header
	for node.next(0) != nil {
		ds = append(ds, node.obj)
		node = node.next(0)
	}
	return ds
}

/*返回节点持有的元素对象*/
func (s *skipList) getByRank(rank int) interface{} {
	n := s.header
	traversed := 0
	for l := s.level - 1; l >= 0; l-- {
		for n.levels[l].forward != nil && traversed+n.levels[l].span <= rank {
			traversed += n.levels[l].span
			n = n.levels[l].forward
		}
		if traversed == rank {
			return n.obj
		}
	}
	return nil
}

func (s *skipList) getRank(obj interface{}) (rank int, err error) {
	n := s.header
	for l := s.level - 1; l >= 0; l-- {
		for n.levels[l].forward != nil {
			rtv, e := compareWithRecover(s.compare, n.levels[i].forward.obj, obj)
			if e != nil {
				rank = 0
				err = e
				return
			}
			if rtv > 0 {
				rank += n.levels[l].span
				n = n.levels[l].forward
			} else if rtv < 0 {
				break
			} else if rtv == 0 && n.levels[i].forward.obj == obj {
				rank += n.levels[l].span
				return
			} else { // rtv ==0 && n.levels[i].forward.obj != obj
				if l == 0 { //最底层不再break
					rank += 1
					n = n.next(0)
				} else {
					break
				}
			}
		}
	}
	return 0
}

func (s *skipList) getTops(num int) ([]interface{}, error) {
	if s.length < num {
		return nil, errors.New("giving number is bigger than the size of set")
	}
	ds := []interface{}{}
	count := 0
	n := s.header
	for count < num {
		ds = append(ds, n.next(0).obj)
		n = n.next(0)
		count++
	}
	return ds, nil
}

/*得到start 到end 范围内节点数据*/
func (s *skipList) getRangeByRank(start, end int) ([]interface{}, error) {
	if start >= end || s.length < end {
		return nil, errors.New("wrong index of set")
	}
	n := s.header
	traversed := 0
	//找到start 的节点
	for l := s.level - 1; l >= 0; l-- {
		for n.levels[l].forward != nil && traversed+n.levels[l].span <= start {
			traversed += n.levels[l].span
			n = n.levels[l].forward
		}
		if traversed == start {
			break //此时n为start序号的节点
		}
	}

	ds := []interface{}{}
	count := 0
	num := end - start + 1
	for count < num {
		ds = append(ds, n.obj)
		n = n.next(0)
		count++
	}
	return ds, nil
}

func compareWithRecover(f func(interface{}, interface{}) int, obj1 interface{}, obj2 interface{}) (rtv int, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	rtv = f(obj1, obj2)
	return
}
