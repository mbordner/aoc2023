package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
)

var (
	reLine = regexp.MustCompile(`^(R|D|L|U)\s+(\d+)\s+\(\#([a-f|0-9]{6})\)$`)
)

func main() {
	//fmt.Println(geom.Positions[int]{{X: 1, Y: 1}, {X: 4, Y: 1}, {X: 2, Y: 5}, {X: 1, Y: 1}}.ShoelaceArea())
	points := getPoints("../data.txt")
	fmt.Println(points.PicksTheoremArea())
}

func getPoints(filename string) geom.Positions[int64] {
	lines := files.MustGetLines(filename)
	positions := make(geom.Positions[int64], 0, len(lines)+1)
	positions = append(positions, geom.Pos[int64]{})
	for _, line := range lines {
		var p geom.Pos[int64]
		matches := reLine.FindStringSubmatch(line)
		s, _ := strconv.ParseInt(matches[2], 10, 64)
		switch matches[1] {
		case "R":
			p = geom.Pos[int64]{X: 1}
		case "L":
			p = geom.Pos[int64]{X: -1}
		case "U":
			p = geom.Pos[int64]{Y: -1}
		case "D":
			p = geom.Pos[int64]{Y: 1}
		}
		p = p.Scale(s)
		positions = append(positions, positions[len(positions)-1].Add(p))
	}
	return positions
}
