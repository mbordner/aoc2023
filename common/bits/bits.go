package bits

type Ints interface {
	int | int32 | int64 | uint | uint32 | uint64
}

func Toggle[T Ints](val T, pos uint) T {
	mask := T(1 << pos)
	return val ^ mask // xor
}
