package array

import (
	"strconv"
	"strings"
	"unsafe"
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

func ToNumbers[V Number](s, sep string) []V {
	tokens := strings.Split(s, sep)
	nums := make([]V, len(tokens), len(tokens))
	var t V
	stbits := 8 * int(unsafe.Sizeof(t))
	for i := range tokens {
		val, _ := strconv.ParseInt(tokens[i], 10, stbits)
		nums[i] = V(val)
	}
	return nums
}

func CloneNumbers[V Number](a []V) []V {
	n := make([]V, len(a), len(a))
	copy(n, a)
	return n
}

func ReverseNumbers[V Number](a []V) []V {
	b := CloneNumbers[V](a)
	lh := len(a) / 2
	for i, j := 0, len(b)-1; i < lh; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func SumNumbers[V Number](a []V) V {
	s := V(0)
	for _, v := range a {
		s += v
	}
	return s
}

func AllSameNumbers[V Number](a []V) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			return false
		}
	}
	return true
}
