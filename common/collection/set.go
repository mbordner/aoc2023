package collection

type Set[T comparable] struct {
	values map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{values: make(map[T]bool)}
}

func (set *Set[T]) Merge(o *Set[T]) {
	for _, v := range o.Values() {
		set.Add(v)
	}
}

func (set *Set[T]) Add(value T) {
	set.values[value] = true
}
func (set *Set[T]) Remove(value T) {
	delete(set.values, value)
}
func (set *Set[T]) Contains(value T) bool {
	b, ok := set.values[value]
	return ok && b
}
func (set *Set[T]) Len() int {
	return len(set.values)
}
func (set *Set[T]) Clear() {
	set.values = make(map[T]bool)
}
func (set *Set[T]) Values() []T {
	values := make([]T, 0, len(set.values))
	for v := range set.values {
		values = append(values, v)
	}
	return values
}
