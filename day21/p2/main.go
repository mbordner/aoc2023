package main

import (
	"fmt"
	"slices"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
)

const (
	Start           = 'S'
	Garden          = '.'
	Rock            = '#'
	DestinationPlot = 'O'
)

func main() {
	grid, start, gardens, rocks := getGarden("../test1.txt")

	stepsGenerator, steps := NewStepsGenerator(grid, gardens, rocks, start)
	stepsCount := 0
	for stepsCount < 5001 {
		steps = stepsGenerator.Generate(steps)
		stepsCount++
		if slices.Contains([]int{6, 10, 50, 100, 500, 1000, 5000}, stepsCount) {
			fmt.Printf("In exactly %d steps, he can reach %d garden plots.\n", stepsCount, len(steps))
		}
	}

	fmt.Println(steps.Len())
}

type StepsGenerator struct {
	grid    common.Grid
	gardens common.PosContainer
	rocks   common.PosContainer
}

func NewStepsGenerator(grid common.Grid, gardens common.PosContainer, rocks common.PosContainer, start common.Pos) (*StepsGenerator, common.PosContainer) {
	sg := &StepsGenerator{grid: grid, gardens: gardens, rocks: rocks}
	pc := make(common.PosContainer)
	pc[start] = true
	return sg, pc
}

func (s *StepsGenerator) Generate(cur common.PosContainer) common.PosContainer {
	dirs := common.Positions{common.Pos{X: 1}, common.Pos{Y: 1}, common.Pos{X: -1}, common.Pos{Y: -1}}
	next := make(common.PosContainer)
	h := len(s.grid)
	w := len(s.grid[0])
	for p, b := range cur {
		if b {
			for _, dir := range dirs {
				np := p.Add(dir)

				// if np is off grid, we need to bring it back to our grid window to check available spots, we'll
				// clone it to cp and transform it back to our grid window
				cp := np.Add(common.Pos{})
				if cp.X < 0 {
					x := cp.X % w
					cp.X = x + w
				}
				if cp.X >= w {
					cp.X = cp.X % w
				}

				if cp.Y < 0 {
					y := cp.Y % h
					cp.Y = y + h
				}
				if cp.Y >= h {
					cp.Y = cp.Y % h
				}

				if s.gardens.Has(cp) {
					next[np] = true
				}
			}
		}
	}
	return next
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
