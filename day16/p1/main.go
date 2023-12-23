package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/file"
	"github.com/mbordner/aoc2023/common/geom"
)

type step struct {
	d geom.Direction
	p geom.Pos
}

type stepCountMap map[step]int
type tileCountMap map[geom.Pos]int
type tilesMap map[geom.Pos]byte

func takeStep(bb *geom.BoundingBox, grid tilesMap, stepCounts stepCountMap, tileCounts tileCountMap, s step) {

	if v, e := stepCounts[s]; e {
		stepCounts[s] = v + 1
	} else {
		stepCounts[s] = 1
	}

	if stepCounts[s] == 2 {
		return
	}

	if v, e := tileCounts[s.p]; e {
		tileCounts[s.p] = v + 1
	} else {
		tileCounts[s.p] = 1
	}

	t := grid[s.p]

	nextSteps := make([]step, 0, 2)
	switch t {
	case '.':
		nextSteps = append(nextSteps, step{d: s.d, p: s.p.TransformDirs(s.d)[0]})
	case '/':
		switch s.d {
		case geom.North:
			nextSteps = append(nextSteps, step{d: geom.East, p: s.p.Transform(1, 0, 0)})
		case geom.East:
			nextSteps = append(nextSteps, step{d: geom.North, p: s.p.Transform(0, -1, 0)})
		case geom.South:
			nextSteps = append(nextSteps, step{d: geom.West, p: s.p.Transform(-1, 0, 0)})
		case geom.West:
			nextSteps = append(nextSteps, step{d: geom.South, p: s.p.Transform(0, 1, 0)})
		}
	case '\\':
		switch s.d {
		case geom.North:
			nextSteps = append(nextSteps, step{d: geom.West, p: s.p.Transform(-1, 0, 0)})
		case geom.East:
			nextSteps = append(nextSteps, step{d: geom.South, p: s.p.Transform(0, 1, 0)})
		case geom.South:
			nextSteps = append(nextSteps, step{d: geom.East, p: s.p.Transform(1, 0, 0)})
		case geom.West:
			nextSteps = append(nextSteps, step{d: geom.North, p: s.p.Transform(0, -1, 0)})
		}
	case '|':
		if s.d == geom.East || s.d == geom.West {
			nextSteps = append(nextSteps, step{d: geom.North, p: s.p.Transform(0, -1, 0)})
			nextSteps = append(nextSteps, step{d: geom.South, p: s.p.Transform(0, 1, 0)})
		} else {
			nextSteps = append(nextSteps, step{d: s.d, p: s.p.TransformDirs(s.d)[0]})
		}
	case '-':
		if s.d == geom.North || s.d == geom.South {
			nextSteps = append(nextSteps, step{d: geom.East, p: s.p.Transform(1, 0, 0)})
			nextSteps = append(nextSteps, step{d: geom.West, p: s.p.Transform(-1, 0, 0)})
		} else {
			nextSteps = append(nextSteps, step{d: s.d, p: s.p.TransformDirs(s.d)[0]})
		}
	}

	for _, ns := range nextSteps {
		if bb.Contains(ns.p) {
			takeStep(bb, grid, stepCounts, tileCounts, ns)
		}
	}
}

func print(bb *geom.BoundingBox, grid tilesMap, counts tileCountMap) {
	pss := bb.GetPositions()
	chars := make([]rune, len(pss))
	for i, p := range pss {
		if _, e := counts[p]; e {
			chars[i] = '#'
		} else {
			chars[i] = rune(grid[p])
		}
	}
	lines := array.Reverse(bb.GetPrintLines(rune('.'), chars, pss))
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	tileCounts := make(tileCountMap)
	stepCounts := make(stepCountMap)
	bb, grid := getData("../data.txt")

	//print(bb, grid, counts)

	startP := geom.Pos{Y: bb.MinY, X: bb.MinX}
	takeStep(bb, grid, stepCounts, tileCounts, step{d: geom.East, p: startP})

	print(bb, grid, tileCounts)
	fmt.Println(len(tileCounts))
}

func getData(path string) (*geom.BoundingBox, tilesMap) {
	lines, _ := file.GetLines(path)

	bb := geom.BoundingBox{}
	grid := make(tilesMap)

	for j, y := len(lines)-1, 0; j >= 0; j, y = j-1, y-1 {
		for i := 0; i < len(lines[j]); i++ {
			p := geom.Pos{Y: y, X: i}
			grid[p] = lines[j][i]
			bb.Extend(p)
		}
	}

	return &bb, grid
}
