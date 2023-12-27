package datastructure

import (
	"container/heap"
)

type AnyHeap[T any] struct {
	data []T
	cmp  func(a, b T) int
}

func (h *AnyHeap[T]) Less(i, j int) bool { return h.cmp(h.data[i], h.data[j]) < 0 }
func (h *AnyHeap[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *AnyHeap[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

func (h *AnyHeap[T]) Pop() any {
	l := len(h.data)
	v := h.data[l-1]
	h.data = h.data[0 : l-1]
	return v
}

// do not use the above methods

func (h *AnyHeap[T]) Len() int { return len(h.data) }

func (h *AnyHeap[T]) Shift() T {
	return heap.Pop(h).(T)
}

func (h *AnyHeap[T]) Unshift(v T) {
	heap.Push(h, v)
}

func (h *AnyHeap[T]) Peek() T {
	return h.data[0]
}

func NewAnyHeap[T any](cmp func(a, b T) int) *AnyHeap[T] {
	h := new(AnyHeap[T])
	h.cmp = cmp
	h.data = make([]T, 0)
	heap.Init(h)
	return h
}
