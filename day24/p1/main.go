package main

import (
	"fmt"
	"strings"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom3"
)

func main() {
	//lines, bounds := getData("../test1.txt", 7, 27)
	lines, bounds := getData("../data.txt", 200000000000000, 400000000000000)
	pairs := common.GetPairSets(lines)
	count := 0
	for _, pair := range pairs {
		p := pair[0].IntersectionPoint(pair[1])
		if p != nil {
			if bounds.Contains(p) {
				t1 := (p.X - float64(pair[0].Point.X)) / float64(pair[0].Direction.X)
				t2 := (p.X - float64(pair[1].Point.X)) / float64(pair[1].Direction.X)
				t3 := (p.Y - float64(pair[0].Point.Y)) / float64(pair[0].Direction.Y)
				t4 := (p.Y - float64(pair[1].Point.Y)) / float64(pair[1].Direction.Y)
				if t1 >= 0 && t2 >= 0 && t3 >= 0 && t4 >= 0 {
					fmt.Println("------")
					fmt.Println(pair[0])
					fmt.Println(pair[1])
					fmt.Println("intersect at:", p)
					count++
				}
			}
		}
	}
	fmt.Println(count)
}

func getData(filename string, min, max float64) (geom3.Lines[int64], geom3.Bounds[float64]) {
	replacer := strings.NewReplacer(" ", "", "@", ",")
	flines := files.MustGetLines(filename)
	lines := make(geom3.Lines[int64], 0, len(flines))
	for _, fline := range flines {
		vals := common.IntVals[int64](replacer.Replace(fline))
		vals[2], vals[5] = 0, 0
		lines = append(lines, geom3.NewLineFromVals(vals))
	}
	return lines, geom3.Bounds[float64]{Min: geom3.Point[float64]{X: min, Y: min}, Max: geom3.Point[float64]{X: max, Y: max}}
}
