package main

// Package main will parse a data input grid from a file, and calculate the number
// of destination points an elf would end on if they walked 26501365 steps from the start position (S).
// the elf can walk north, south, east or west towards any garden plot (.) that is not a rock (#).
// the grid repeats infinitely in all directions.
// Also, S is considered a garden plot, i.e. can be a destination location.

import (
	"fmt"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
)

const (
	Start           = 'S'
	Garden          = '.'
	Rock            = '#'
	DestinationPlot = 'O'
	StepsToTake     = 26501365
)

func main() {
	grid, start, gardens, _ := getGarden("../data.txt")

	h := len(grid)
	w := len(grid[0])

	// assumptions:
	// our data input grid is square, and starting position is in the center
	// the border row and columns of the grid, and the row and column where S is located has no rocks
	// these empty rows and columns will allow the steps to generally fill out to the entire width and height from start for
	// our desired steps count.
	// because of these empty rows and columns from start, and on the borders, and the grid is generally sparse, the geometry will generally
	// grow to be a diamond (as if the grid was completely open).
	// we assume and will verify that this number of steps needed is a multiple of our width (or height since it's square) with an additional
	// half of the width + 1.
	// this half of the width + 1, e.g. for our example 131/2+1 = 65 is also the position of our start r,c in the 131x131 sample grid
	// because of the nature of sample data grid, if we find all destinations from start from 65 steps out, the destinations extents
	// will have the same w*h as our grid.  this is because the steps can expand without disruption N, S, E and W from S.
	// also, to find the next consecutive result shapes that will be a multiples of our w (or h), we need to add multiples of w to 65.
	// for example, if we create the following series from the following steps to find values:
	// { (0 * w) + 65, (1 * w) + 65, (2 * w) + 65, (3 * w) + 65.... }  the values will be:
	// { 3835, 34125, 94603, 185269, ... }  notice that if these values are {y0, y1, y2, y3 }
	// { y1-y0, y2-y1, y3-y2 } == { 30290, 60478, 90666 }
	// and also notice that
	// { (y2-y1)-(y1-y0), (y3-y2)-(y2-y1) } == { 30188, 30188 }
	// this 2nd differences list of values shows that the 2nd differences are constant.
	// what this means for sure, is that our size of steps are not growing linearly, but in way the area
	// of the diamond is growing linearly.
	//
	// In other words the series will fit to some quadratic function such that:
	// f(n) = an^2 + bn + c,  where n would be the interval of grid widths, and f(n) being the number of destination plots
	// we can find values to solve this with some points (notice we are switching x coordinate to be full grid width, starting at 0)
	// e.g. if we are looking for 1 full width, which we know would be 65 steps:  65-65/131 = 0
	// (0,y0) = (0,3835) [ 1 full width ((65)-65)/131 ] = 0
	// (1,y1) = (1,34125) [ 2 full widths ((131+65)-65)/131 ] = 1
	// (2,y2) = (2,94603)
	// (3,y3) = (3,185269)
	//
	// let's set up the system for the first 3 points:
	// y0 = a(0)^2 + b(0) + c    ->   c = y0
	// y1 = a(1)^2 + b(1) + c    ->  y1 = a + b + c
	// y2 = a(2)^2 + b(2) + c    ->  y2 = 4a + 2b + c
	//
	// so...     a +  b = y1 - y0
	//          4a + 2b = y2 - y0
	// subtracting twice the first equation from the second ([4a + 2b] - [2a + 2b])
	//  2a = (y2 - y0) - 2(y1 - y0)
	//  2a = y2 -2y1 + y0
	//   a =  (y2 - 2y1 + y0) / 2
	//  and..
	//   b = y1 - y0 - a
	//   b = y1 - y0 - (y2 - 2y1 + y0) / 2
	//
	// once we have a, b and c, we can plug in 202300 full grid widths [which is (26501365-65)/131]
	// f(202300) = a*202300^2 + b*202300 + c

	// grid is square, and distance from start to either edge would be the same if w&h are odd
	if h != w && w%2 != 0 {
		panic("the grid must be square of some odd number")
	}

	s := w / 2 // we know that since w and h are both odd that w%2 == 1 and h%2 == 1
	if s != start.X || s != start.Y {
		panic("expecting start position to be directly in the center, with 64 chars to the edge in either N/S/E/W direction")
	}

	for x := 0; x < w; x++ {
		if grid[0][x] == Rock || grid[h-1][x] == Rock || grid[s][x] == Rock {
			panic("assumed top and bottom border row to be open, as well as row where S is")
		}
	}

	for y := 0; y < h; y++ {
		if grid[y][0] == Rock || grid[y][s] == Rock || grid[y][w-1] == Rock {
			panic("assumed left and right border column to be open, as well as the column where S is")
		}
	}

	dest := getStepDestinations(grid, gardens, start, s)
	y0 := dest.Len()
	g := getDestinationsAsGrid(grid, dest)

	if len(g)%h != 0 || len(g[0])%h != 0 {
		panic("expecting result grid extents to be a multiple of w*h")
	}

	dest = getStepDestinations(grid, gardens, start, s+w)
	y1 := dest.Len()
	g = getDestinationsAsGrid(grid, dest)
	if len(g)%h != 0 || len(g[0])%h != 0 {
		panic("expecting result grid extents to be a multiple of w*h")
	}

	dest = getStepDestinations(grid, gardens, start, s+w+w)
	y2 := dest.Len()
	g = getDestinationsAsGrid(grid, dest)
	if len(g)%h != 0 || len(g[0])%h != 0 {
		panic("expecting result grid extents to be a multiple of w*h")
	}

	dest = getStepDestinations(grid, gardens, start, s+w+w+w)
	y3 := dest.Len()
	g = getDestinationsAsGrid(grid, dest)
	if len(g)%h != 0 || len(g[0])%h != 0 {
		panic("expecting result grid extents to be a multiple of w*h")
	}

	y := []int{y0, y1, y2, y3}

	firstDifferences := make([]int, 0, 3)
	for i := 1; i < len(y); i++ {
		firstDifferences = append(firstDifferences, y[i]-y[i-1])
	}
	secondDifferences := make([]int, 0, 2)
	for i := 1; i < len(firstDifferences); i++ {
		secondDifferences = append(secondDifferences, firstDifferences[i]-firstDifferences[i-1])
	}

	if secondDifferences[0] != secondDifferences[1] {
		panic("expecting area rate of change to be constant")
	}

	c := y[0]
	a := (y[2] - 2*y[1] + c) / 2
	b := y[1] - c - a

	t := uint64(3)
	tAns := uint64(a)*(t*t) + uint64(b)*t + uint64(c)

	// test this with our y3
	if tAns != uint64(y[3]) {
		panic("we tested this with our results from BFS for steps (s+w+w+w) using the getStepDestinations function, but the result didn't match")
	}

	n := uint64(202300) // number [(26501365-65)/131] of full widths to get our answer
	ans := uint64(a)*(n*n) + uint64(b)*n + uint64(c)

	fmt.Println(ans)
}

func getDestinationsAsGrid(grid common.Grid, destinations common.PosContainer) common.Grid {
	minP, maxP := destinations.Extents()
	gh := maxP.Y - minP.Y + 1
	gw := maxP.X - minP.X + 1
	g := make(common.Grid, gh)
	for y := minP.Y; y <= maxP.Y; y++ {
		gy := y - minP.Y
		g[gy] = make([]byte, gw)
		for x := minP.X; x <= maxP.X; x++ {
			gx := x - minP.X
			if destinations.Has(common.Pos{X: x, Y: y}) {
				g[gy][gx] = DestinationPlot
			} else {
				gp := transformToGridWindow(grid, common.Pos{X: x, Y: y})
				g[gy][gx] = grid[gp.Y][gp.X]
			}
		}
	}
	return g
}

type QueueState struct {
	p     common.Pos
	steps int
}

func transformToGridWindow(grid common.Grid, p common.Pos) common.Pos {
	h := len(grid)
	w := len(grid[0])
	gp := p.Add(common.Pos{})
	if gp.X < 0 {
		x := gp.X % w
		gp.X = x + w
	}
	if gp.X >= w {
		gp.X = gp.X % w
	}
	if gp.Y < 0 {
		y := gp.Y % h
		gp.Y = y + h
	}
	if gp.Y >= h {
		gp.Y = gp.Y % h
	}
	return gp
}

func getStepDestinations(grid common.Grid, gardens common.PosContainer, start common.Pos, numSteps int) common.PosContainer {
	destinations := make(common.PosContainer)
	visited := make(common.VisitedState[common.Pos, bool])

	queue := make(common.Queue[QueueState], 0, 100)
	queue.Enqueue(QueueState{p: start, steps: numSteps})
	visited[start] = true

	dirs := common.Positions{common.Pos{X: 1}, common.Pos{Y: 1}, common.Pos{X: -1}, common.Pos{Y: -1}}

	for !queue.Empty() {
		cur := *(queue.Dequeue())

		if cur.steps%2 == 0 {
			destinations[cur.p] = true
		}

		if cur.steps > 0 {

			for _, dir := range dirs {
				np := cur.p.Add(dir)

				if !visited.Has(np) {

					cp := transformToGridWindow(grid, np)

					if gardens.Has(cp) {
						visited[np] = true
						queue.Enqueue(QueueState{p: np, steps: cur.steps - 1})
					}

				}

			}

		}
	}

	return destinations
}

func getGarden(filename string) (common.Grid, common.Pos, common.PosContainer, common.PosContainer) {
	grid := common.ConvertGrid(files.MustGetLines(filename))
	var start common.Pos
	gardenPlots := make(common.PosContainer)
	rockPlots := make(common.PosContainer)
	for y := range grid {
		for x := range grid[y] {
			p := common.Pos{Y: y, X: x}
			if grid[y][x] == Start {
				start = p
			}
			if grid[y][x] == Rock {
				rockPlots[p] = true
			} else {
				gardenPlots[p] = true
			}
		}
	}
	return grid, start, gardenPlots, rockPlots
}
