package main

import (
	"fmt"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
)

const (
	Forest     = '#'
	Path       = '.'
	SlopeEast  = '>'
	SlopeNorth = '^'
	SlopeWest  = '<'
	SlopeSouth = 'v'
	Step       = 'O'
)

type State struct {
	pos   common.Pos
	steps int
}

func main() {
	grid, start, goal := getGrid("../test1.txt")
	fmt.Println(len(grid), start, goal)
	h, w := len(grid), len(grid[0])
	queue := make(common.Queue[State], 0, w*h)
	visited := make(common.VisitedState[common.Pos, common.Pos])

	queue.Enqueue(State{pos: start, steps: 0})

	slopeTraversable := make(map[byte]common.Pos)
	slopeTraversable[SlopeNorth] = common.DN
	slopeTraversable[SlopeEast] = common.DE
	slopeTraversable[SlopeSouth] = common.DS
	slopeTraversable[SlopeWest] = common.DW

	maxPathLen := 0

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		if cur.pos == goal {

			if cur.steps > maxPathLen {
				maxPathLen = cur.steps
			}

		} else {

			for _, dir := range common.AdjacentDirs {
				np := cur.pos.Add(dir)
				if np.Y >= 0 && np.Y < h && np.X >= 0 && np.X < w {
					if !visited.Has(np) && grid[np.Y][np.X] != Forest {
						traversable := true
						if d, e := slopeTraversable[grid[np.Y][np.X]]; e {
							if d != dir {
								traversable = false
							}
						}
						if traversable {
							visited.Set(np, cur.pos)
							queue.Enqueue(State{pos: np, steps: cur.steps + 1})
						}
					}
				}
			}

		}
	}

	fmt.Println(maxPathLen)
}

func getGrid(filename string) (common.Grid, common.Pos, common.Pos) {
	grid := common.ConvertGrid(files.MustGetLines(filename))
	var start, goal common.Pos
	for x := 0; x < len(grid[0]); x++ {
		if grid[0][x] == Path {
			start.X = x
			break
		}
	}
	goal.Y = len(grid) - 1
	for x := 0; x < len(grid[goal.Y]); x++ {
		if grid[goal.Y][x] == Path {
			goal.X = x
			break
		}
	}
	return grid, start, goal
}
