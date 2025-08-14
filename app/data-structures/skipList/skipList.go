package skiplist

import "math/rand"

const (
	MAX_HEIGHT  = 16
	PROBABILITY = 0.5
)

type Node struct {
	Member string
	Score  float64
	Tower  [MAX_HEIGHT]*Node
}

type List struct {
	Head   *Node
	Height int
	Len    int
}

func New() *List {
	return &List{
		Head:   &Node{},
		Height: 1,
	}
}

func (list *List) Insert(score float64, member string) int {
	found, update := list.search(score, member)
	if found != nil {
		return 0
	}

	nodeHeight := list.randomHeight()
	if nodeHeight > list.Height {
		list.increaseHeight(&update, nodeHeight)
	}

	newNode := &Node{
		Member: member,
		Score:  score,
	}

	for level := range nodeHeight {
		newNode.Tower[level] = update[level].Tower[level]
		update[level].Tower[level] = newNode
	}

	list.Len++
	return 1
}

func (list *List) Delete(score float64, member string) int {
	found, update := list.search(score, member)
	if found == nil {
		return 0
	}

	for level := range list.Height {
		update[level].Tower[level] = found.Tower[level]
		found.Tower[level] = nil
	}

	list.decreaseHeight()

	list.Len--
	return 1
}

func (list *List) search(score float64, member string) (*Node, [MAX_HEIGHT]*Node) {
	var found *Node
	var update [MAX_HEIGHT]*Node

	cur := list.Head
	for level := list.Height - 1; level >= 0; level-- {
		for cur.Tower[level] != nil {
			next := cur.Tower[level]

			if score == next.Score && next.Member == member {
				found = next
				// Don't return to fill all levels for update array
				break
			}

			if score < next.Score || (score == next.Score && member < next.Member) {
				break
			}

			cur = next
		}

		update[level] = cur
	}

	return found, update
}

func (list *List) increaseHeight(update *[MAX_HEIGHT]*Node, newHeight int) {
	for level := list.Height; level < newHeight; level++ {
		update[level] = list.Head
	}
	list.Height = newHeight
}

func (list *List) decreaseHeight() {
	for list.Height > 1 && list.Head.Tower[list.Height-1] == nil {
		list.Height--
	}
}

func (list *List) randomHeight() int {
	height := 1
	for rand.Float64() > PROBABILITY && height < MAX_HEIGHT {
		height++
	}
	return height
}
