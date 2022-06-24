package skiplist

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	initLevel   = 1
	maxLevel    = 32
	probability = 0.25
)

type Skiplist struct {
	Level    int              // Current highest level of skip list
	RankName string           // rank name
	players  map[string]*Node // players node
	Head     *Node
	length   int64
}
type Node struct {
	playerId  string // player id
	Value     int64
	Timestamp int64 // millisecond when add
	Next      *Node
	Pre       *Node
	Down      *Node
}

func NewSkiplist(rankName string) Skiplist {
	return Skiplist{
		Level:    initLevel,
		RankName: rankName,
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

// 添加节点
// @param.player 玩家标识
// @param.num 玩家分数
func (s *Skiplist) Add(player string, num int64) error {

	node, ok := s.players[player]
	if ok { // 存在 ，先删除
		if err := s.DeleteByPlayerId(node); err != nil {
			return err
		}
	}

	// 增加
	s.add(player, num)
	return nil
}
func (s *Skiplist) add(playerId string, num int64) {
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
		// 2. p.next.value < value
		// 判断新插入的节点层级
		if level >= l { // 新节点层级大于当前层级
			node := s.NewNode(playerId, num, p, p.Next, nil)
			p.Next = node
			if up != nil {
				up.Down = p.Next
			}
			up = p.Next

			// 新节点层级大于当前层级 ， 在这里赋值新节点最高层级的指针
			if top == nil && level > s.Level {
				top = p.Next
			}
			// 没超过当前层级的第一层
			if level == l {
				s.players[playerId] = p.Next
			}
		}
		p = p.Down
	}

	// 新节点层级大于当前层级，重置头节点
	if level > s.Level {
		s.Level = level

		// 头节点后的第一个节点
		node := s.NewNode(playerId, num, s.Head, nil, top)
		// 头节点
		headNode := s.NewNode("", -1, nil, node, s.Head)
		s.Head = headNode

		s.players[playerId] = headNode.Next
	}
	return
}

func (s *Skiplist) NewNode(player string, num int64, pre, next, down *Node) *Node {
	return &Node{
		playerId:  player,
		Value:     num,
		Timestamp: time.Now().UnixNano() / 1e6,
		Next:      next,
		Pre:       pre,
		Down:      down,
	}
}
func (s *Skiplist) Print() {
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

func (s *Skiplist) Get(playerId string) bool {
	_, ok := s.players[playerId]
	if !ok {
		return false
	}
	return true
}

// 通过playerId删除
func (s *Skiplist) DeleteByPlayerId(node *Node) error {
	for node != nil {
		// 1.该层只有一个节点
		// 2.该层很多节点
		pre := node.Pre
		next := node.Next
		if pre.Value == -1 && next == nil { // 符合1.只有一个节点
			if s.Head.Down == nil { // 第一层
				s.Head.Next = nil
				break
			}
			s.Head = s.Head.Down
			s.Level--
		} else { // 符合2.该层很多节点
			pre.Next = next
		}
		node = node.Down
	}
	return nil
}

// 通过遍历跳跃表删除
func (s *Skiplist) Delete(playerId string, num int64) bool {
	node := s.Head
	isDel := false
	for node != nil {
		for node.Next != nil && node.Next.Value <= num {
			if node.Next.Value == num && node.Next.playerId == playerId {
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

func (s *Skiplist) GetRanges(start, end int64) {

}
