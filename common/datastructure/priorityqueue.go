package datastructure

import "container/heap"

// Item is a generic struct to hold the value and its priority.
type Item[T any] struct {
	Value    T   // The item value
	Priority int // The item priority (lower is better for min-heap)
	index    int // The index in the heap array (managed by container/heap)
}

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T any] []*Item[T]

func (pq *PriorityQueue[T]) Len() int { return len(*pq) }

// Less is the custom comparison function. For a min-heap, we check if i's priority is less than j's.
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return (*pq)[i].Priority < (*pq)[j].Priority
}

func (pq *PriorityQueue[T]) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T]) // Type assertion
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	item.index = -1 // For safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue[T]) PushItem(value T, priority int) {
	item := &Item[T]{Value: value, Priority: priority}
	heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) PopItem() (T, bool) {
	if pq.Len() == 0 {
		var zero T // Return zero value for type T
		return zero, false
	}
	item := heap.Pop(pq).(*Item[T]) // Type assertion
	return item.Value, true
}
