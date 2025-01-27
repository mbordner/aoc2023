package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
	"github.com/mbordner/aoc2023/common/graph"
)

var (
	reInstr = regexp.MustCompile(`(\w)\s+(\d+)\s\(#(\w+)\)`)
	dirMap  = map[string]geom.Direction{"U": geom.North, "R": geom.East, "D": geom.South, "L": geom.West}

	cornerNW = geom.Direction(geom.South | geom.East)
	cornerNE = geom.Direction(geom.South | geom.West)
	cornerSW = geom.Direction(geom.North | geom.East)
	cornerSE = geom.Direction(geom.North | geom.West)
)

func main() {
	bb, g := getData("../data.txt")
	fmt.Println(bb.GetPositionsSize(), g.Len())
}

func getData(path string) (geom.BoundingBox[int], *graph.Graph) {
	bb := geom.BoundingBox[int]{}
	g := graph.NewGraph()

	sp := geom.Pos[int]{}
	start := g.CreateNode(sp)
	prev := start

	lines, _ := files.GetLines(path)
	for j, line := range lines {
		matches := reInstr.FindStringSubmatch(line)
		if matches != nil {
			d := dirMap[matches[1]]
			od := d.Opposite()
			c, _ := strconv.ParseInt(matches[2], 10, 64)
			col := matches[3]
			if j == 0 {
				start.AddProperty("color", col)
			}
			for i := int64(0); i < c; i++ {
				pp := prev.GetID().(geom.Pos[int])
				np := pp.TransformDirs(d)[0]
				var next *graph.Node
				if np == sp {
					next = start
				} else {
					bb.Extend(np)
					next = g.CreateNode(np)
					next.AddProperty("color", col)
				}
				e := prev.AddEdge(next, 1)
				e.AddProperty("dir", d)
				e = next.AddEdge(prev, 1)
				e.AddProperty("dir", od)
				prev = next
			}
		}
	}

	for j := bb.MinY; j <= bb.MaxY; j++ {
		inside := false
		for i := bb.MinX; i <= bb.MaxX; i++ {
			tp := geom.Pos[int]{Y: j, X: i}
			tn := g.GetNode(tp)
			if tn != nil {
				ed := getEdgeDirs(tn)
				if isCorner(ed) {
					for ci := i + 1; ci <= bb.MaxX; ci++ {
						tcp := geom.Pos[int]{Y: j, X: ci}
						tcn := g.GetNode(tcp)
						if tcn != nil {
							ted := getEdgeDirs(tcn)
							if isCorner(ted) {
								if !closesArea(ed, ted) {
									inside = !inside
								}
								i = ci
								break
							}
						} else {
							log.Fatalln("we should have found a corner")
						}
					}
				} else {
					inside = !inside
				}
			} else {
				if inside {
					g.CreateNode(tp)
				}
			}
		}
	}

	return bb, g
}

func closesArea(c1, c2 geom.Direction) bool {
	if c1 == cornerNW && c2 == cornerSE {
		return false
	}
	if c1 == cornerSW && c2 == cornerNE {
		return false
	}
	return true
}

func getEdgeDirs(n *graph.Node) geom.Direction {
	d := geom.Unknown
	for _, e := range n.GetEdges() {
		ed := e.GetProperty("dir").(geom.Direction)
		d |= ed
	}
	return d
}

func isCorner(d geom.Direction) bool {
	corners := []geom.Direction{cornerNW, cornerNE, cornerSW, cornerSE}
	return array.Contains[geom.Direction](corners, d)
}
