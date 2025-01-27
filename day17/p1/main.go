package main

import (
	"fmt"

	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
)

func main() {
	bb, hl := getData("../test1.txt")
	start := geom.Pos[int]{Y: bb.MinY, X: bb.MinX}
	end := geom.Pos[int]{Y: bb.MaxY, X: bb.MaxX}

	loss := heatLoss(bb, hl, end, geom.Unknown, 0, start, []geom.Pos[int]{})
	fmt.Println(loss)
}

type dpks struct {
	d  geom.Direction
	dc int
	p  geom.Pos[int]
}

type DP map[dpks]int
type VC map[geom.Pos[int]]int

var (
	dp = make(DP)
)

func heatLoss(bb geom.BoundingBox[int], hl [][]byte, end geom.Pos[int], d geom.Direction, dc int, p geom.Pos[int], path []geom.Pos[int]) int {
	dpk := dpks{d: d, dc: dc, p: p}
	if v, e := dp[dpk]; e {
		return v
	}

	if len(path) == len(hl)*len(hl[0]) {
		return -1
	}

	thisHeadLoss := hl[p.Y][p.X]
	if p == end {
		return int(thisHeadLoss)
	}

	minLoss := -1

	possibleDirs := d.Not([]geom.Direction{geom.North, geom.East, geom.South, geom.West}, []geom.Direction{d.Opposite()})
	if d != geom.Unknown && dc < 3 {
		possibleDirs = append([]geom.Direction{d}, possibleDirs...)
	}

	np := append(path, p)

	for _, pd := range possibleDirs {
		pdc := 1
		if pd == d {
			pdc = dc + 1
		}
		pp := p.TransformDirs(pd)[0]
		if bb.Contains(pp) && !array.Contains(np, pp) {
			phl := heatLoss(bb, hl, end, pd, pdc, pp, np)
			if phl >= 0 {
				if minLoss == -1 {
					minLoss = phl
				} else if phl < minLoss {
					minLoss = phl
				}
			}
		}
	}

	if minLoss == -1 {
		//dp[dpk] = minLoss
		return minLoss
	} else {
		if v, e := dp[dpk]; e {
			if v > minLoss+int(thisHeadLoss) {
				dp[dpk] = minLoss + int(thisHeadLoss)
			}
		} else {
			dp[dpk] = minLoss + int(thisHeadLoss)
		}
	}

	return dp[dpk]
}

func getData(path string) (geom.BoundingBox[int], [][]byte) {
	bb := geom.BoundingBox[int]{}
	lines, _ := files.GetLines(path)
	heatLoss := make([][]byte, len(lines))
	for j, line := range lines {
		bytes := make([]byte, len(line))
		for i, b := range line {
			bb.Extend(geom.Pos[int]{Y: j, X: i})
			bytes[i] = byte(b - '0')
		}
		heatLoss[j] = bytes
	}
	return bb, heatLoss
}
