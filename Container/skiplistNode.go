package container

import (
	"math/rand"
)

const (
	MAXLEVEL = 32
	PART     = 4
)

type skiplistNode struct {
	obj      interface{} //指向真实节点的指针
	backward *skiplistNode
	levels   []*skiplistLevel
}

func newSkiplistNode(level int, obj interface{}) *skiplistNode {
	levels := make([]*skiplistLevel, level)
	for i := 0; i < level; i++ {
		levels[i] = new(skiplistLevel)
	}
	return &skiplistNode{
		obj:      obj,
		backward: nil,
		levels:   levels,
	}
}

type skiplistLevel struct {
	forward *skiplistNode
	span    int
}

func (sn *skiplistNode) next(level int) *skiplistNode {
	return sn.levels[level].forward
}

func (sn *skiplistNode) prev() *skiplistNode {
	return sn.backward
}

func randomLevel() int {
	level := 1
	for (rand.Int63()&0xFFFF)%PART == 0 { //25%的概率增长层数
		level += 1
	}
	if level < MAXLEVEL {
		return level
	} else {
		return MAXLEVEL
	}
}
