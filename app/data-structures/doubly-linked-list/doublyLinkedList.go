package doublylinkedlist

type Node struct {
	Val  string
	Next *Node
	Prev *Node
}

type List struct {
	Head *Node
	Tail *Node
	Len  int
}

func InsertInTheStart(list *List, n *Node) {
	if list.Head == nil {
		list.Head = n
		list.Tail = n
	} else {
		n.Next = list.Head
		list.Head.Prev = n
		list.Head = n
	}
	list.Len++
}

func InsertInTheEnd(list *List, n *Node) {
	if list.Tail == nil {
		list.Tail = n
		list.Head = n
	} else {
		list.Tail.Next = n
		n.Prev = list.Tail
		list.Tail = n
	}
	list.Len++
}

func DeleteFromStart(list *List) *Node {
	if list.Head == nil {
		return nil
	}
	deleted := list.Head
	Next := list.Head.Next
	if Next == nil {
		list.Head = nil
		list.Tail = nil
	} else {
		list.Head.Next = nil
		Next.Prev = nil
		list.Head = Next
	}
	list.Len--
	return deleted
}

func DeleteFromEnd(list *List) *Node {
	if list.Tail == nil {
		return nil
	}
	deleted := list.Tail
	Prev := list.Tail.Prev
	if Prev == nil {
		list.Head = nil
		list.Tail = nil
	} else {
		list.Tail.Prev = nil
		Prev.Next = nil
		list.Tail = Prev
	}
	list.Len--
	return deleted
}
