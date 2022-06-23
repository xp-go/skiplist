package skiplist

import (
	"fmt"
	"math/rand"
)

const (
	initLevel   = 1
	maxLevel    = 32   // 最高层
	probability = 0.25 //上一层概率
)

type Skiplist struct {
	Level int // 当前跳跃表最高层级
	Head  *Node
}
type Node struct {
	Value int
	Next  *Node
	Down  *Node
}

func NewSkiplist() Skiplist {
	return Skiplist{
		Level: initLevel,
		Head: &Node{
			Value: -1,
		},
	}
}

// 插入的层数
func (s *Skiplist) GetLevel() int {
	level := initLevel
	for rand.Float64() <= probability || level < maxLevel {
		level++
	}
	if level-s.Level > 1 {
		level = s.Level + 1
	}
	return level
}

func (s *Skiplist) Add(num int) {

	// 获取新节点层级
	level := s.GetLevel()

	// 获取最高层，第一个数
	p := s.Head

	var top *Node = nil // 新节点的第一个插入点
	var up *Node = nil  // 新节点的上层节点

	for l := s.Level; l >= 1; l-- {
		for p.Next != nil && p.Next.Value < num {
			// 节点不为空，且下一个节点数值小于value
			p = p.Next
		}

		// 1. p.next == nil
		///2. p.next.value < value
		// 判断新插入的节点层级
		if level >= l { // 新节点层级大于当前层级
			p.Next = &Node{
				Value: num,
				Next:  p.Next,
				Down:  nil,
			}
			if up != nil {
				up.Down = p.Next
			}
			up = p.Next

			// 新节点层级大于当前层级 ， 在这里赋值新节点最高层级的指针
			if top == nil && level > s.Level {
				top = p.Next
			}
		}
		p = p.Down
	}

	// 新节点层级大于当前层级，重置头节点
	if level > s.Level {
		s.Level = level
		s.Head = &Node{
			Value: -1,
			Next: &Node{
				Value: num,
				Next:  nil,
				Down:  top,
			},
			Down: s.Head,
		}
	}
	return
}

func (s *Skiplist) Prin() {
	p := s.Head
	top := s.Head
	for l := s.Level; l >= 1; l-- {
		for p != nil {
			fmt.Print(p.Value)
			fmt.Print(" ")
			p = p.Next
		}
		fmt.Println("")
		p = top.Down
		top = top.Down
	}
}

func (s *Skiplist) Get(target int) bool {
	node := s.Head
	for node != nil {
		for node.Next != nil && node.Next.Value <= target {
			if node.Next.Value == target {
				return true
			}
			node = node.Next
		}
		node = node.Down
	}
	return false
}

func (s *Skiplist) Delete(num int) bool {
	node := s.Head
	isDel := false
	for node != nil {
		for node.Next != nil && node.Next.Value <= num {
			if node.Next.Value == num {
				node.Next = node.Next.Next
				isDel = true
				continue
			}
			node = node.Next
		}
		node = node.Down
	}
	return isDel
}
