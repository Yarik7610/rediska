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
	Span   [MAX_HEIGHT]int
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
	found, update, rank := list.Search(score, member)
	if found != nil {
		return 0
	}

	nodeHeight := list.randomHeight()
	if nodeHeight > list.Height {
		list.increaseHeight(nodeHeight, update, rank)
	}

	newNode := &Node{
		Member: member,
		Score:  score,
	}

	for level := range nodeHeight {
		newNode.Tower[level] = update[level].Tower[level]
		update[level].Tower[level] = newNode

		// rank[0] - rank[i] shows count of elements between update[level] (not including) and newNode on 0 level
		newNode.Span[level] = update[level].Span[level] - (rank[0] - rank[level])
		update[level].Span[level] = (rank[0] - rank[level]) + 1
	}

	for i := nodeHeight; i < list.Height; i++ {
		update[i].Span[i]++
	}

	list.Len++
	return 1
}

func (list *List) Delete(score float64, member string) int {
	found, update, _ := list.Search(score, member)
	if found == nil {
		return 0
	}

	for level := range list.Height {
		if update[level].Tower[level] == found {
			update[level].Span[level] += found.Span[level] - 1
			update[level].Tower[level] = found.Tower[level]
		} else {
			update[level].Span[level]--
		}
	}

	list.shrinkHeight()

	list.Len--
	return 1
}

func (list *List) Search(score float64, member string) (*Node, *[MAX_HEIGHT]*Node, *[MAX_HEIGHT]int) {
	var found *Node
	var update [MAX_HEIGHT]*Node
	var rank [MAX_HEIGHT]int

	cur := list.Head
	for level := list.Height - 1; level >= 0; level-- {
		if level == list.Height-1 {
			rank[level] = 0
		} else {
			rank[level] = rank[level+1]
		}

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

			rank[level] += cur.Span[level]
			cur = next
		}

		update[level] = cur
	}

	return found, &update, &rank
}

func (list *List) increaseHeight(newHeight int, update *[MAX_HEIGHT]*Node, rank *[MAX_HEIGHT]int) {
	for level := list.Height; level < newHeight; level++ {
		rank[level] = 0
		update[level] = list.Head
		update[level].Span[level] = list.Len
	}
	list.Height = newHeight
}

func (list *List) shrinkHeight() {
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
