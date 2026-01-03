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

func main() {
	grid, start, goal := getGrid("../data.txt")
	fmt.Println(len(grid), start, goal)

	path := getLongestPath(grid, start, goal, make(common.PosContainer))

	fmt.Println(len(path) - 2)
}

func getLongestPath(grid common.Grid, pos common.Pos, goal common.Pos, visited common.PosContainer) common.Positions {
	var longestPath common.Positions

	visited[pos] = true

	if pos == goal {
		longestPath = common.Positions{goal}
	} else {
		h, w := len(grid), len(grid[0])

		for _, dir := range common.AdjacentDirs {
			np := pos.Add(dir)
			if np.Y >= 0 && np.Y < h && np.X >= 0 && np.X < w {
				if !visited.Has(np) && grid[np.Y][np.X] != Forest {
					path := getLongestPath(grid, np, goal, visited)
					if len(path) > len(longestPath) {
						longestPath = path
						longestPath = path
					}
				}
			}

		}
	}

	visited[pos] = false
	if len(longestPath) > 0 {
		return append(common.Positions{pos}, longestPath...)
	}

	return longestPath
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
