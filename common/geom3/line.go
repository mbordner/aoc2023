package geom3

import (
	"fmt"
	"math"
	"strings"
)

const (
	epsilon = 1e-9
	zero    = float64(0)
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

type Point[T Number] struct {
	X T
	Y T
	Z T
}

func (p *Point[T]) String() string {
	return fmt.Sprintf("(%v, %v, %v)", p.X, p.Y, p.Z)
}

func (p *Point[T]) Equal(o *Point[T]) bool {
	return AlmostEqual(float64(p.X), float64(o.X)) && AlmostEqual(float64(p.Y), float64(o.Y)) && AlmostEqual(float64(p.Z), float64(o.Z))
}

type Bounds[T Number] struct {
	Min Point[T]
	Max Point[T]
}

func (b *Bounds[T]) Contains(p *Point[T]) bool {
	if p.X < b.Min.X || p.X > b.Max.X {
		return false
	}
	if p.Y < b.Min.Y || p.Y > b.Max.Y {
		return false
	}
	if p.Z < b.Min.Z || p.Z > b.Max.Z {
		return false
	}
	return true
}

type Lines[T Number] []*Line[T]

type Line[T Number] struct {
	Point     Point[T]
	Direction Vector[T]
}

func (l *Line[T]) String() string {
	vals := make([]string, 6)
	vals[0] = fmt.Sprintf("%v", l.Point.X)
	vals[1] = fmt.Sprintf("%v", l.Point.Y)
	vals[2] = fmt.Sprintf("%v", l.Point.Z)
	vals[3] = fmt.Sprintf("%v", l.Direction.X)
	vals[4] = fmt.Sprintf("%v", l.Direction.Y)
	vals[5] = fmt.Sprintf("%v", l.Direction.Z)
	return strings.Join(vals, ",")
}

func (l *Line[T]) IntersectionPoint(o *Line[T]) *Point[float64] {
	if l.Coincident(o) || l.Parallel(o) {
		return nil
	}

	v1 := l.Direction
	v2 := o.Direction
	// Vector between the two starting points
	dp := Vector[T]{X: o.Point.X - l.Point.X, Y: o.Point.Y - l.Point.Y, Z: o.Point.Z - l.Point.Z}

	// 1. Check if lines are Skew using the Scalar Triple Product: (v1 x v2) · dp
	cp := v1.Cross(v2)
	// Dot product of cp and dp
	dot := cp.Dot(dp)

	// If the volume of the parallelepiped defined by the three vectors is not 0,
	// the lines are skew and cannot intersect.
	if !AlmostEqual(float64(dot), zero) {
		return nil
	}

	// To find s and t, we can use Cramer's rule on the 2D projection with the largest components
	// This avoids division by zero issues from lines aligned with axes.
	var s, t float64
	denom := float64(v1.X*(-v2.Y)) - float64((-v2.X)*v1.Y)

	if ABS(denom) > epsilon {
		s = (float64(dp.X)*float64(-v2.Y) - float64(-v2.X)*float64(dp.Y)) / denom
		t = (float64(v1.X)*float64(dp.Y) - float64(dp.X)*float64(v1.Y)) / denom
	} else {
		// Try a different plane (X-Z) if X-Y is degenerate
		denom = float64(v1.X*(-v2.Z)) - float64((-v2.X)*v1.Z)
		if ABS(denom) > epsilon {
			s = (float64(dp.X)*float64(-v2.Z) - float64(-v2.X)*float64(dp.Z)) / denom
			t = (float64(v1.X)*float64(dp.Z) - float64(dp.X)*float64(v1.Z)) / denom
		} else {
			// Try Y-Z plane
			denom = float64(v1.Y*(-v2.Z)) - float64((-v2.Y)*v1.Z)
			s = (float64(dp.Y)*float64(-v2.Z) - float64(-v2.Y)*float64(dp.Z)) / denom
			t = (float64(v1.Y)*float64(dp.Z) - float64(dp.Y)*float64(v1.Z)) / denom
		}
	}

	p0 := &Point[float64]{
		X: float64(l.Point.X) + s*float64(l.Direction.X),
		Y: float64(l.Point.Y) + s*float64(l.Direction.Y),
		Z: float64(l.Point.Z) + s*float64(l.Direction.Z),
	}

	p1 := &Point[float64]{
		X: float64(o.Point.X) + t*float64(o.Direction.X),
		Y: float64(o.Point.Y) + t*float64(o.Direction.Y),
		Z: float64(o.Point.Z) + t*float64(o.Direction.Z),
	}

	// if the lines actually intersect, p0 and p1 must be the same
	if p0.Equal(p1) {
		return p0
	}

	return nil
}

func (l *Line[T]) Skew(o *Line[T]) bool {
	if l.Parallel(o) {
		return false
	}

	// Vector between starting points
	dp := Vector[T]{X: o.Point.X - l.Point.X, Y: o.Point.Y - l.Point.Y, Z: o.Point.Z - l.Point.Z}

	// Scalar Triple Product: (v1 x v2) · dp
	// If this is non-zero, the lines are skew.
	cp := l.Direction.Cross(o.Direction)
	tripleProduct := cp.Dot(dp)

	return !AlmostEqual(float64(tripleProduct), zero)
}

func (l *Line[T]) Contains(p *Point[T]) bool {
	// Vector from geom3 start to the point
	vToP := Vector[T]{
		X: p.X - l.Point.X,
		Y: p.Y - l.Point.Y,
		Z: p.Z - l.Point.Z,
	}

	// If the point is on the geom3, the direction and vToP are parallel.
	// Parallel vectors have a cross product of zero.
	cp := l.Direction.Cross(vToP)

	return AlmostEqual(float64(cp.X), zero) &&
		AlmostEqual(float64(cp.Y), zero) &&
		AlmostEqual(float64(cp.Z), zero)
}

// Coincident answers if p1 == p2 + s*v2 ?
func (l *Line[T]) Coincident(o *Line[T]) bool {
	// they must be parallel
	if !l.Parallel(o) {
		return false
	}
	// if they are parallel, they are coincident if one geom3's
	// starting point lies on the other geom3.
	return l.Contains(&o.Point)
}

// Parallel answers if for l1 = Point + s Direction, and l2 = Point + t Direction, whether the two Direction vectors are scalar multiples of each other
func (l *Line[T]) Parallel(o *Line[T]) bool {
	cp := l.Direction.Cross(o.Direction)
	return AlmostEqual(float64(cp.X), zero) && AlmostEqual(float64(cp.Y), zero) && AlmostEqual(float64(cp.Z), zero)
}

func NewLineFromVals[T Number](vals []T) *Line[T] {
	if len(vals) != 6 {
		panic("Line expects 6 values")
	}
	return &Line[T]{Point: Point[T]{X: vals[0], Y: vals[1], Z: vals[2]}, Direction: Vector[T]{X: vals[3], Y: vals[4], Z: vals[5]}}
}

type Vector[T Number] struct {
	X T
	Y T
	Z T
}

func (v Vector[T]) Subtract(o Vector[T]) Vector[T] {
	return v.Transform(-o.X, -o.Y, -o.Z)
}

func (v Vector[T]) Scale(s T) Vector[T] {
	return Vector[T]{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

func (v Vector[T]) Transform(x, y, z T) Vector[T] {
	return Vector[T]{X: v.X + x, Y: v.Y + y, Z: v.Z + z}
}

func (v Vector[T]) Add(o Vector[T]) Vector[T] {
	return v.Transform(o.X, o.Y, o.Z)
}

func (v Vector[T]) Magnitude() float64 {
	return math.Sqrt(float64(v.X*v.X) + float64(v.Y*v.Y) + float64(v.Z*v.Z))
}

func (v Vector[T]) Dot(o Vector[T]) T {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector[T]) Cross(o Vector[T]) Vector[T] {
	return Vector[T]{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

func (v Vector[T]) Normalize() Vector[float64] {
	m := v.Magnitude()
	return Vector[float64]{X: float64(v.X) / m, Y: float64(v.Y) / m, Z: float64(v.Z) / m}
}

func MIN[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func MAX[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func GCD[T Number](a, b T) T {
	a = ABS(a)
	b = ABS(b)
	for b != 0 {
		a, b = b, T(int64(a)%int64(b))
	}
	return a
}

func LCM[T Number](nums ...T) T {
	var res T
	if len(nums) > 1 {
		res = nums[0]
		for i := 1; i < len(nums); i++ {
			res = ABS(res*nums[i]) / GCD(res, nums[i])
		}
	}
	return res
}

func ABS[T Number](a T) T {
	if a < T(0) {
		return a * T(-1)
	}
	return a
}

func AlmostEqual(a, b float64) bool {
	if a == b {
		return true
	}
	diff := ABS(a - b)
	return diff < epsilon*MAX(ABS(a), ABS(b))
}
