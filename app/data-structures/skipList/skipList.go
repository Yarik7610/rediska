package skiplist

const (
	MAX_LEVEL   = 16
	PROBABILITY = 0.5
)

type Node struct {
	Member string
	Score  float64
	Levels [MAX_LEVEL]*Node
}

type List struct {
	Head   *Node
	Len    int
	Height int
}

func New() *List {
	return &List{
		Head:   &Node{},
		Height: 1,
		Len:    0,
	}
}
