package sortedset

import (
	"math/rand"
)

const (
	MAXLEVEL = 32
	PART     = 4
)

type skiplistNode struct {
	value    interface{}
	backward *skiplistNode
	levels   []*skiplistLevel
	score    float64
}

func newSkiplistNode(level int, value interface{}, score float64) {
	levels := make([]*skiplistLevel, level)
	for i := 0; i < level; i++ {
		levels[i] = new(skiplistLevel)
	}
	return &skiplistNode{
		value:    value,
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

func randonLevel() int {
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
