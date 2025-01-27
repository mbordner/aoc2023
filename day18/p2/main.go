package main

import (
	"cmp"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/datastructure"
	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
)

var (
	reInstr   = regexp.MustCompile(`(\w)\s+(\d+)\s\(#(\w+)\)`)
	dirMap    = map[string]geom.Direction{`0`: geom.East, `1`: geom.South, `2`: geom.West, `3`: geom.North}
	clockWise = true
)

func pCmp[T geom.IntNumber](a, b geom.Pos[T]) int {
	if a.X == b.X {
		return cmp.Compare(a.Y, b.Y)
	} else if a.Y == b.Y {
		return cmp.Compare(a.X, b.X)
	} else {
		log.Fatalln("row or col should have been the same")
	}
	return 0
}

type RowColPos[T geom.IntNumber] map[T]*datastructure.AnyHeap[geom.Pos[T]]

func (rcp RowColPos[T]) Add(v T, p geom.Pos[T]) {
	if _, e := rcp[v]; !e {
		rcp[v] = datastructure.NewAnyHeap[geom.Pos[T]](pCmp[T])
	}
	rcp[v].Unshift(p)
}

func main() {
	bb, pgls, gls := getData("../test1.txt")
	clockWise = isClockWise(pgls)
	fmt.Println(bb.GetPositionsSize(), len(pgls), len(gls), clockWise)

	rowsp := make(RowColPos[int64])
	colsp := make(RowColPos[int64])

	for p := range pgls {
		colsp.Add(p.X, p)
		rowsp.Add(p.Y, p)
	}

	rows := array.SortedKeys(rowsp)
	cols := array.SortedKeys(colsp)

	blocks := make([][]bool, len(rows)-1)
	for j, r := range rows {
		if j < len(rows)-1 {
			blocks[j] = make([]bool, len(cols)-1)
			inside := false
			for i, c := range cols {
				if i < len(cols)-1 {
					tp := geom.Pos[int64]{Y: r, X: c}
					if lines, e := pgls[tp]; e {
						dir, line := getVerticalGridLine(lines)
						if dir != geom.Unknown {
							isLineBelow := false
							if (dir == geom.North && line.P1 == tp) ||
								(dir == geom.South && line.P0 == tp) {
								isLineBelow = true
							}
							if isLineBelow {
								insideChange := false
								if dir == geom.North {
									insideChange = true
								}
								if !clockWise {
									insideChange = !insideChange
								}
								inside = insideChange
							}
						}
					} else {
						for tj := j - 1; tj >= 0; tj-- {
							ttp := geom.Pos[int64]{Y: rows[tj], X: c}
							if _, e := pgls[ttp]; e {
								inside = blocks[tj][i]
								break
							}
						}
					}
					blocks[j][i] = inside
				}
			}
		}
	}

	lineBoxes := make(LineBoxMap)
	boxes := make([]*geom.GridBox[int64], 0, len(cols)*len(rows))

	for j := 0; j < len(rows)-1; j++ {
		areas := make([]int64, len(cols)-1)
		for i := 0; i < len(cols)-1; i++ {
			if blocks[j][i] {
				nw := geom.Pos[int64]{Y: rows[j], X: cols[i]}
				ne := geom.Pos[int64]{Y: rows[j], X: cols[i+1]}
				sw := geom.Pos[int64]{Y: rows[j+1], X: cols[i]}
				se := geom.Pos[int64]{Y: rows[j+1], X: cols[i+1]}
				box := geom.NewGridBox(nw, ne, se, sw)
				top := geom.GridLine[int64]{P0: nw, P1: ne}
				bottom := geom.GridLine[int64]{P0: sw, P1: se}
				left := geom.GridLine[int64]{P0: nw, P1: sw}
				right := geom.GridLine[int64]{P0: ne, P1: se}
				boxes = append(boxes, &box)
				lineBoxes.Add(top, &box)
				lineBoxes.Add(bottom, &box)
				lineBoxes.Add(left, &box)
				lineBoxes.Add(right, &box)
				areas[i] = box.Area()
			}
		}
		fmt.Println(areas)
	}

	area := int64(0)
	for _, b := range boxes {
		area += b.Area()
	}

	/*
		for l, gbs := range lineBoxes {
			d := len(gbs) - 1
			if d > 0 {
				area -= l.Length() * int64(d)
			}
		}
	*/

	fmt.Println(area)
}

type LineBoxMap map[geom.GridLine[int64]][]*geom.GridBox[int64]

func (lb LineBoxMap) Add(l geom.GridLine[int64], gb *geom.GridBox[int64]) {
	if gbs, e := lb[l]; e {
		lb[l] = append(gbs, gb)
	} else {
		lb[l] = []*geom.GridBox[int64]{gb}
	}
}

func getVerticalGridLine(lines []geom.GridLine[int64]) (geom.Direction, *geom.GridLine[int64]) {
	for _, line := range lines {
		d := line.Direction()
		if d == geom.North || d == geom.South {
			return d, &line
		}
	}
	return geom.Unknown, nil
}

// looking for the line that comes back to start, the direction will tell us if it is
// clockwise direction or not
func isClockWise(pgls geom.PosGridLines[int64]) bool {
	start := geom.Pos[int64]{}
	if lines, e := pgls[start]; e {
		for _, line := range lines {
			if line.P1 == start {
				d := line.Direction()
				switch d {
				case geom.North:
					return true
				case geom.East:
					return true
				case geom.South:
					return false
				case geom.West:
					return false
				}
			}
		}
	} else {
		log.Fatalln("where's the start")
	}
	return false
}

func getData(path string) (geom.BoundingBox[int64], geom.PosGridLines[int64], geom.GridLines[int64]) {
	bb := geom.BoundingBox[int64]{}

	pgls := make(geom.PosGridLines[int64])
	gls := make(geom.GridLines[int64])

	start := geom.Pos[int64]{}
	pp := start

	lines, _ := files.GetLines(path)
	for _, line := range lines {
		matches := reInstr.FindStringSubmatch(line)
		if matches != nil {
			hs := matches[3]
			hs0 := hs[0 : len(hs)-1]
			hs1 := hs[len(hs)-1:]
			d := dirMap[hs1]
			m, _ := strconv.ParseInt(hs0, 16, 64)

			// part 1
			d = map[string]geom.Direction{"U": geom.North, "R": geom.East, "D": geom.South, "L": geom.West}[matches[1]]
			m, _ = strconv.ParseInt(matches[2], 10, 64)

			np := pp.TransformDir(d, m)
			bb.Extend(np)

			gl := geom.GridLine[int64]{P0: pp, P1: np}
			pgls.AddLine(gl)
			gls[gl] = d

			pp = np
		}
	}

	return bb, pgls, gls
}
