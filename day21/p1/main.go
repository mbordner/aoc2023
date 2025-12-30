package main

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
)

func main() {
	_, start, gardens, rocks := getGarden("../data.txt")

	stepsGenerator, steps := NewStepsGenerator(gardens, rocks, start)
	stepsCount := 0
	for stepsCount < 64 {
		steps = stepsGenerator.Generate(steps)
		stepsCount++
		//fmt.Println("-----", stepsCount, "-----")
		//printGrid(grid, steps)
	}

	fmt.Println(steps.Len())
}

func printGrid(grid common.Grid, steps common.PosContainer) {
	g := grid.Clone()
	for p, b := range steps {
		if b {
			g[p.Y][p.X] = DestinationPlot
		}
	}
	g.Print()
}

type StepsGenerator struct {
	gardens common.PosContainer
	rocks   common.PosContainer
}

func NewStepsGenerator(gardens common.PosContainer, rocks common.PosContainer, start common.Pos) (*StepsGenerator, common.PosContainer) {
	sg := &StepsGenerator{gardens: gardens, rocks: rocks}
	pc := make(common.PosContainer)
	pc[start] = true
	return sg, pc
}

func (s *StepsGenerator) Generate(cur common.PosContainer) common.PosContainer {
	dirs := common.Positions{common.Pos{X: 1}, common.Pos{Y: 1}, common.Pos{X: -1}, common.Pos{Y: -1}}
	next := make(common.PosContainer)
	for p, b := range cur {
		if b {
			for _, dir := range dirs {
				np := p.Add(dir)
				if s.gardens.Has(np) {
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
