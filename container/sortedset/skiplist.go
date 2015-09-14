package sortedset

import ()

/*
  it's not concurrent safe and not ensure unique
*/
type SkipList struct {
	header  *skiplistNode
	tail    *skiplistNode
	update  []*skiplistNode
	rank    []int
	length  int
	level   int
	compare func(obj1 interface{}, obj2 interface{}) int //两个对象进行比较 ，相等返回0，obj1>obj2返回一个大于0的数，obj1<obj2返回一个小于0的数
}

func newSkiplist(cf func(obj1, obj2 interface{}) int) *SkipList {
	if cf == nil {
		return nil
	}
	return &SkipList{
		header:  newSkiplistNode(MAXLEVEL, nil, 0),
		tail:    nil,
		update:  make([]*skiplistNode, MAXLEVEL),
		rank:    make([]int, MAXLEVEL),
		length:  0,
		level:   1,
		compare: cf,
	}
}

func (s *SkipList) Flush() {
	s.header = newSkiplistNode(MAXLEVEL, nil, 0)
	s.tail = nil
	s.update = make([]*skiplistNode, MAXLEVEL)
	s.rank = make([]int, MAXLEVEL)
	s.length = 0
	s.level = 1
}

func (s *SkipList) len() int {
	return s.length
}

func (s *SkipList) insert(objPointer interface{}) error {
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
			if rtv, e = compareWithRecover(s.compare, n.levels[i].forward.objPointer, objPointer); e == nil && rtv < 0 {
				s.rank[i] += n.levels[i].span
				n = n.levels[i].forward
				break
			}
			if e != nil {
				return e
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

	n = newSkiplistNode(level, obj, objScore)
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
	return n
}

func (s *SkipList) delete(objPointer interface{}) error {
	node, err := findPreNodes(objPointer)
	if err != nil {
		return err
	}
	if node != nil && s.compare(node.objPointer, objPointer) == 0 { //不再需要recover
		for l := 0; l < s.level; l++ {
			if s.update[l].levels[l].forward == node {
				s.update[l].levels[l].span += node.levels[l].span - 1
				s.update[l].levels[l].forward = node.levels[l].forward
			} else {
				s.update[l].levels.span -= 1
			}
		}
	} else { //未找到相等节点
		return nil
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
}

func (s *SkipList) findPreNodes(objPointer interface{}) (node *skiplistNode, err error) {
	n := s.header
	for i := s.level - 1; i >= 0; i-- {
		for n.levels[i].forward != nil {
			if rtv, e := compareWithRecover(s.compare, n.levels[i].forward.objPointer, objPointer); e == nil && rtv < 0 {
				n = n.levels[i].forward
				break
			}
			if e != nil {
				err = e
				return
			}
		}
		s.update[i] = n
	}
	return s.update[0].forward, nil
}

func (s *SkipList) contain(objPointer interface{})(bool,error) {
   node, err := findPreNodes(objPointer)
   if err!=nil{
      return false,err
   }
   if node != nil && s.compare(node.objPointer, objPointer) == 0 {
      return true,nil
   }else{
     return false,nil
   }
}

/*返回节点持有的元素对象指针*/
func (s *SkipList) getByRank(rank int)interface{} {
   n := s.header
   traversed := 0
   for l:=s.level - 1 ;l>=0; l-- {
      for n.level[l].forward !=nil && traversed + n.level[l].span <= rank{
	     traversed += n.level[l].span
		 n = n.level[l].forward
	  }
	  if traversed == rank {
	     return n.objPointer
	  }
   }
   return nil
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
