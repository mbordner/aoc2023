package ranges

/*
Package ranges provides a utility for managing ranges as efficiently as possible.

Ranges keeps track of l and r ranges.  and will merge overlaps when added. l and r are
inclusive so they are part of the range.  adding (l,r) to an empty range adds r-l+1 to the length.
*/

import (
	"sort"

	"github.com/pkg/errors"
)

type Number interface {
	int | int32 | float32 | int64 | uint64 | float64
}

type side int

const (
	leftSide side = iota
	rightSide
)

type el[T Number] struct {
	val T
	s   side
}

type elements[T Number] []el[T]

func (e elements[T]) Len() int {
	return len(e)
}

func (e elements[T]) Less(a, b int) bool {
	if e[a].val < e[b].val {
		return true
	}
	if e[a].val == e[b].val {
		if e[a].s == leftSide && e[b].s == rightSide {
			return true
		}
	}
	return false
}

func (e elements[T]) Swap(a, b int) {
	e[a], e[b] = e[b], e[a]
}

type Collection[T Number] struct {
	values []T
}

func (c *Collection[T]) Clone() *Collection[T] {
	o := new(Collection[T])
	o.values = make([]T, len(c.values))
	copy(o.values, c.values)
	return o
}

func (c *Collection[T]) ValuePairs() []T {
	return c.values
}

func (c *Collection[T]) Len() T {
	var l T
	for i := 0; i < len(c.values); i += 2 {
		l += c.values[i+1] - c.values[i] + 1
	}
	return l
}

// Remove l-r from the range (l,r are inclusive). returns the new range, or an error
func (c *Collection[T]) Remove(l, r T) ([]T, error) {
	if r < l {
		return nil, errors.New("right value cannot be less than left value")
	}

	start, end := -1, -1
	inStart, inEnd := -1, -1
	for i := 0; i < len(c.values); i += 2 {
		if l >= c.values[i] && l <= c.values[i+1] {
			start = i
		}
		if l < c.values[i] && inStart < 0 {
			inStart = i
		}
		if r >= c.values[i] && r <= c.values[i+1] {
			end = i
		}
		if r < c.values[i] {
			inEnd = i
		}
		if end != -1 || inEnd != -1 {
			break
		}
	}

	if start == -1 && end == -1 {

		if inStart != -1 && inEnd != -1 {
			c.values = append(c.values[0:inStart], c.values[inEnd:]...)
		} else if inStart != -1 {
			c.values = c.values[0:inStart]
		}

	} else {
		// found some overlap

		if start != -1 && end != -1 {
			// start and end overlap

			if start == end { // overlap in same range

				if l == c.values[start] && r == c.values[end+1] { // remove the range
					if start == 0 { // removing from start
						c.values = c.values[2:]
					} else if end+1 == len(c.values)-1 { // removing from end
						c.values = c.values[0 : len(c.values)-2]
					} else { // removing in the middle
						c.values = append(c.values[0:start], c.values[start+2:]...)
					}
				} else { // not removing the range, just adjusting current range

					if l == c.values[start] {

						c.values[start] = r + 1

					} else if r == c.values[end+1] {

						c.values[end+1] = l - 1

					} else { // splitting the range, as the ends don't overlap
						values := make([]T, 0, len(c.values)+2)
						values = append(values, c.values[0:start]...)
						values = append(values, []T{c.values[start], l - 1}...)
						values = append(values, []T{r + 1, c.values[end+1]}...)
						c.values = append(values, c.values[start+2:]...)
					}
				}

			} else { // not overlapping in same range

				c.values[start+1] = l - 1
				c.values[end] = r + 1

				if end-start > 1 {
					c.values = append(c.values[0:start+2], c.values[end:]...)
				}

			}

		} else if start != -1 {
			// start overlaps but end doesn't

			c.values[start+1] = l - 1

		} else {
			// end overlaps but start doesn't

			c.values[end] = r + 1
		}

	}

	// s, e == -1, -1 no overlap
	return c.values, nil
}

// Subtract returns a new collection with other subtracted out
func (c *Collection[T]) Subtract(other *Collection[T]) *Collection[T] {
	s := c.Clone()
	for i := 0; i < len(other.values); i += 2 {
		_, _ = s.Remove(other.values[i], other.values[i+1])
	}
	return s
}

// Add adds l,r to the range (l & r are inclusive), returns the new range or an error
func (c *Collection[T]) Add(l, r T) ([]T, error) {

	if r < l {
		return nil, errors.New("right value cannot be less than left value")
	}

	els := make(elements[T], 0, len(c.values)+2)

	if len(c.values) > 0 {
		for i, v := range c.values {
			s := leftSide
			if i%2 == 1 {
				s = rightSide
			}
			els = append(els, el[T]{val: v, s: s})
		}
	}

	els = append(els, el[T]{val: l, s: leftSide})
	els = append(els, el[T]{val: r, s: rightSide})

	sort.Sort(els)

	values := make([]T, 0, len(c.values)+2)

	leftValStack := make([]T, 0, 10)

	for i := 0; i < len(els); i++ {
		if els[i].s == leftSide {
			leftValStack = append(leftValStack, els[i].val)
		} else if els[i].s == rightSide {
			leftVal := leftValStack[len(leftValStack)-1]
			leftValStack = leftValStack[0 : len(leftValStack)-1]
			if len(leftValStack) == 0 {
				rightVal := els[i].val
				if len(values) > 0 {
					if leftVal == values[len(values)-1] || leftVal == values[len(values)-1]+1 {
						values[len(values)-1] = rightVal
						continue
					}
				}
				values = append(values, leftVal)
				values = append(values, rightVal)
			}
		}

	}

	c.values = values

	return values, nil
}
