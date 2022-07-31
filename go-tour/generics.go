import "fmt"

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func main() {
	root := List[int]{val: 1}
	node1 := List[int]{val: 2}
	node2 := List[int]{val: 3}
	root.next = &node1
	node1.next = &node2

	p := &root
	for {
		fmt.Println(p.val)
		if p.next == nil {
			break
		}
		p = p.next
	}
}
