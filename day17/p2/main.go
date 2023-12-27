package main

import (
	"cmp"
	"fmt"
	"github.com/mbordner/aoc2023/common/cmath"
	"github.com/mbordner/aoc2023/common/datastructure"
	"github.com/mbordner/aoc2023/common/file"
	"github.com/mbordner/aoc2023/common/geom"
	"log"
)

func main() {
	bb, hl := getData("../data.txt")
	start := geom.Pos{Y: bb.MinY, X: bb.MinX}
	end := geom.Pos{Y: bb.MaxY, X: bb.MaxX}

	loss := heatLoss(bb, hl, end, geom.Unknown, 0, start)
	fmt.Println(loss)
}

type dpks struct {
	d  geom.Direction
	dc int
	p  geom.Pos
}

type DP map[dpks]int

type stepObj struct {
	dpk dpks
	dis int
}

func heatLoss(bb geom.BoundingBox, hl [][]byte, end geom.Pos, d geom.Direction, dc int, p geom.Pos) int {
	dp := make(DP)

	stepCmp := func(a, b *stepObj) int {
		return cmp.Compare(a.dis, b.dis)
	}

	stepHeap := datastructure.NewAnyHeap[*stepObj](stepCmp)

	stepHeap.Unshift(&stepObj{dis: 0, dpk: dpks{p: p, d: d, dc: dc}})

	for stepHeap.Len() > 0 {
		step := stepHeap.Shift()
		if v, e := dp[step.dpk]; e {
			if v <= step.dis {
				continue
			} else {
				log.Fatalln("this shouldn't happen")
			}
		}
		dp[step.dpk] = step.dis
		if step.dpk.p == end {
			continue
		}

		var possibleDirs []geom.Direction
		possibleTurnDirs := step.dpk.d.Not([]geom.Direction{geom.North, geom.East, geom.South, geom.West}, []geom.Direction{step.dpk.d.Opposite()})
		if step.dpk.d == geom.Unknown {
			possibleDirs = possibleTurnDirs
		} else {
			if step.dpk.dc < 4 {
				possibleDirs = []geom.Direction{step.dpk.d}
			} else if step.dpk.dc < 10 {
				possibleDirs = append([]geom.Direction{step.dpk.d}, possibleTurnDirs...)
			} else {
				possibleDirs = possibleTurnDirs
			}
		}

		for _, pd := range possibleDirs {
			pdc := 1
			if pd == step.dpk.d {
				pdc = step.dpk.dc + 1
			}
			pp := step.dpk.p.TransformDirs(pd)[0]
			if bb.Contains(pp) {
				heatLossVal := int(hl[pp.Y][pp.X])
				pdpk := dpks{p: pp, d: pd, dc: pdc}
				stepHeap.Unshift(&stepObj{dis: step.dis + heatLossVal, dpk: pdpk})
			}
		}
	}

	mv := cmath.MaxInt
	for dpk, dis := range dp {
		if dpk.p == end && dpk.dc > 3 && dis < mv {
			mv = dis
		}
	}
	return mv
}

func getData(path string) (geom.BoundingBox, [][]byte) {
	bb := geom.BoundingBox{}
	lines, _ := file.GetLines(path)
	heatLoss := make([][]byte, len(lines))
	for j, line := range lines {
		bytes := make([]byte, len(line))
		for i, b := range line {
			bb.Extend(geom.Pos{Y: j, X: i})
			bytes[i] = byte(b - '0')
		}
		heatLoss[j] = bytes
	}
	return bb, heatLoss
}
