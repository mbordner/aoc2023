package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom3"
	"github.com/mbordner/aoc2023/common/matrices"
)

func main() {

	lines := getData("../data.txt")

	fmt.Println(solve(lines[0], lines[1], lines[2]))
}

// solve determines the rock's starting position (x, y, z) and velocity (vx, vy, vz)
// by setting up a system of linear equations based on three hailstones.
func solve(h0, h1, h2 *geom3.Line[int64]) int64 {
	// we have 6 unknowns: xr, yr, zr (position) and vxr, vyr, vzr (velocity).
	// to solve for 6 variables, we need 6 independent linear equations.
	// each pair of hailstones gives us equations by linearizing the 3D intersection.
	matrix := make([][]float64, 6)
	for i := range matrix {
		matrix[i] = make([]float64, 7)
	}

	// fillRows generates linear equations by eliminating the time variable (t).
	// for a rock to hit hailstone 'i' at time 'ti': Pr + ti*Vr = Pi + ti*Vi.
	// this implies (Pr - Pi) is parallel to (Vi - Vr), meaning their cross product is zero:
	// (Pr - Pi) x (Vr - Vi) = 0.
	// expanding this for multiple hailstones allows us to cancel out the non-linear
	// terms (Pr x Vr) and create a system of linear equations.
	fillRows := func(rowOffset int, a, b *geom3.Line[int64]) {

		// 1. linearized equation for the XY plane:
		// (vay - vby)xr + (vbx - vax)yr + (pby - pay)vxr + (pax - pbx)vyr = (pax*vay - pay*vax) - (pbx*vby - pby*vbx)
		matrix[rowOffset][0] = float64(a.Direction.Y - b.Direction.Y)
		matrix[rowOffset][1] = float64(b.Direction.X - a.Direction.X)
		matrix[rowOffset][3] = float64(b.Point.Y - a.Point.Y)
		matrix[rowOffset][4] = float64(a.Point.X - b.Point.X)
		matrix[rowOffset][6] = float64(a.Point.X*a.Direction.Y-a.Point.Y*a.Direction.X) -
			float64(b.Point.X*b.Direction.Y-b.Point.Y*b.Direction.X)

		// 2. linearized equation for the XZ plane:
		matrix[rowOffset+1][0] = float64(a.Direction.Z - b.Direction.Z)
		matrix[rowOffset+1][2] = float64(b.Direction.X - a.Direction.X)
		matrix[rowOffset+1][3] = float64(b.Point.Z - a.Point.Z)
		matrix[rowOffset+1][5] = float64(a.Point.X - b.Point.X)
		matrix[rowOffset+1][6] = float64(a.Point.X*a.Direction.Z-a.Point.Z*a.Direction.X) -
			float64(b.Point.X*b.Direction.Z-b.Point.Z*b.Direction.X)

		// 3. linearized equation for the YZ plane:
		matrix[rowOffset+2][1] = float64(a.Direction.Z - b.Direction.Z)
		matrix[rowOffset+2][2] = float64(b.Direction.Y - a.Direction.Y)
		matrix[rowOffset+2][4] = float64(b.Point.Z - a.Point.Z)
		matrix[rowOffset+2][5] = float64(a.Point.Y - b.Point.Y)
		matrix[rowOffset+2][6] = float64(a.Point.Y*a.Direction.Z-a.Point.Z*a.Direction.Y) -
			float64(b.Point.Y*b.Direction.Z-b.Point.Z*b.Direction.Y)
	}

	// use two pairs (h0, h1) and (h0, h2) to generate the 6 required equations.
	fillRows(0, h0, h1)
	fillRows(3, h0, h2)

	// transform the matrix into Reduced Row Echelon Form (RREF).
	// after reduction, the 7th column (index 6) contains the solved values for our unknowns.
	matrices.ToFloatReducedEchelonForm(matrix)

	// extract the solved starting positions (xr, yr, zr).
	// we use math.Round to account for potential floating-point precision errors with large inputs.
	x := math.Round(matrix[0][6])
	y := math.Round(matrix[1][6])
	z := math.Round(matrix[2][6])

	// the answer is the sum of the starting X, Y, and Z coordinates.
	return int64(x + y + z)
}

func getData(filename string) geom3.Lines[int64] {
	replacer := strings.NewReplacer(" ", "", "@", ",")
	flines := files.MustGetLines(filename)
	lines := make(geom3.Lines[int64], 0, len(flines))
	for _, fline := range flines {
		vals := common.IntVals[int64](replacer.Replace(fline))
		lines = append(lines, geom3.NewLineFromVals(vals))
	}
	return lines
}
