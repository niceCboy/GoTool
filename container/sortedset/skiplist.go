package sortedset

type SkipList struct {
	header *skiplistNode
	tail   *skiplistNode
	update []*skiplistNode
	rank   []int
	length int
	level  int
	score  func(v interface{})float64 //对象获取分值函数
}

func New( f func(v interface{})float64) *SkipList{
   return &SkipList{
     header :  newSkiplistNode(MAXLEVEL, nil, 0),
	 tail : nil,
	 update: make([]*skiplistNode,MAXLEVEL),
	 rank : make([]int,MAXLEVEL),
     length : 0,
	 level:1,
	 score:f,
   }
}

func (s *SkipList)Flush(){
   s.header = newSkiplistNode(MAXLEVEL, nil, 0)
   s.tail = nil
   s.update = make([]*skiplistNode,MAXLEVEL)
   s.rank = make([]int,MAXLEVEL)
   s.length = 0 
   s.level = 1
   //score function unchange
}

func (s *SkipList)Len() int {
  return s.length
}


