package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/file"
	"github.com/mbordner/aoc2023/common/geom"
	"github.com/mbordner/aoc2023/common/graph"
	"github.com/mbordner/aoc2023/common/graph/djikstra"
	"strings"
)

func main() {
	data := getData("../data.txt")
	sky := expandSky(data)
	sky = expandSky(sky)

	g, bb, gm, sn := getGraph(sky) // graph, bounding box, galaxy pos map

	fmt.Println(g.Len(), bb.GetPositionsSize(), len(gm), sn.GetID(), sn.GetProperty("tile").(string))

	sps := make(map[geom.Pos]djikstra.ShortestPaths)

	pairs := array.Pairs(array.Keys(gm))
	for _, pair := range pairs {
		if _, e := sps[pair[0]]; !e {
			n := g.GetNode(pair[0])
			sp := djikstra.GenerateShortestPaths(g, n)
			sps[pair[0]] = sp
		}
	}

	sum := int64(0)

	for _, pair := range pairs {
		p := pair[0]
		q := pair[1]

		nq := g.GetNode(q)
		_, v := sps[p].GetShortestPath(nq)
		sum += int64(v)
	}

	fmt.Println(sum)
}

func getGraph(sky [][]rune) (*graph.Graph, geom.BoundingBox, graph.PosNodeMap[geom.Pos], *graph.Node) {
	bb := geom.BoundingBox{}
	g := graph.NewGraph()
	mg := make(graph.PosNodeMap[geom.Pos])

	var s *graph.Node

	// create nodes
	for j := 0; j < len(sky); j++ {
		for i := 0; i < len(sky[j]); i++ {
			p := geom.Pos{Y: j, X: i}
			if bb.GetPositionsSize() == 0 {
				bb.SetExtents(p.X, p.Y, p.Z, p.X, p.Y, p.Z)
			} else {
				bb.Extend(p)
			}
			t := string(sky[j][i])
			n := g.CreateNode(p)
			n.AddProperty("tile", t)
			if t == "#" {
				mg[p] = n
				if s == nil {
					s = n
				}
			}
		}
	}

	// connect nodes
	for _, p := range bb.GetPositions() {
		n := g.GetNode(p)
		targets := make(map[geom.Direction]geom.Pos) // targets
		targets[geom.North] = p.Transform(0, -1, 0)
		targets[geom.West] = p.Transform(-1, 0, 0)
		targets[geom.East] = p.Transform(1, 0, 0)
		targets[geom.South] = p.Transform(0, 1, 0)
		for td, tp := range targets {
			if bb.Contains(tp) {
				on := g.GetNode(tp)
				e := n.AddEdge(on, 1)
				e.AddProperty("dir", td)
			}
		}
	}

	return g, bb, mg, s
}

func expandSky(data [][]rune) [][]rune {
	data = array.SwapRowCols[rune](data)
	sky := make([][]rune, 0, len(data))
	for _, row := range data {
		sky = append(sky, row)
		if !hasGalaxy(row) {
			sky = append(sky, newRow('.', len(row)))
		}
	}
	return sky
}

func newRow(char rune, length int) []rune {
	r := make([]rune, length)
	for i := range r {
		r[i] = char
	}
	return r
}

func hasGalaxy(row []rune) bool {
	s := string(row)
	return strings.Contains(s, "#")
}

func getData(path string) [][]rune {
	lines, _ := file.GetLines(path)
	data := make([][]rune, len(lines), len(lines))
	for i, line := range lines {
		data[i] = []rune(line)
	}
	return data
}
