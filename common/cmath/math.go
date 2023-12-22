package cmath

type Number interface {
	int | int32 | int64 | float32 | float64
}

func Factorial[V Number](v V) V {
	if v == V(1) {
		return v
	}
	return v * Factorial(v-V(1))
}
