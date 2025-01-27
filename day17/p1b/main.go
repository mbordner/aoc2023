package main

import (
	"cmp"
	"fmt"
	"log"

	"github.com/mbordner/aoc2023/common/datastructure"
	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
)

func main() {
	bb, hl := getData("../data.txt")
	start := geom.Pos[int]{Y: bb.MinY, X: bb.MinX}
	end := geom.Pos[int]{Y: bb.MaxY, X: bb.MaxX}

	loss := heatLoss(bb, hl, end, geom.Unknown, 0, start)
	fmt.Println(loss)
}

type dpks struct {
	d  geom.Direction
	dc int
	p  geom.Pos[int]
}

type DP map[dpks]int

type stepObj struct {
	dpk dpks
	dis int
}

func heatLoss(bb geom.BoundingBox[int], hl [][]byte, end geom.Pos[int], d geom.Direction, dc int, p geom.Pos[int]) int {
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

		possibleDirs := step.dpk.d.Not([]geom.Direction{geom.North, geom.East, geom.South, geom.West}, []geom.Direction{step.dpk.d.Opposite()})
		if step.dpk.d != geom.Unknown && step.dpk.dc < 3 {
			possibleDirs = append([]geom.Direction{step.dpk.d}, possibleDirs...)
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

	mv := int(^uint(0) >> 1)
	for dpk, dis := range dp {
		if dpk.p == end && dis < mv {
			mv = dis
		}
	}
	return mv
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
