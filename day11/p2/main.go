package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/file"
	"github.com/mbordner/aoc2023/common/geom"
	"github.com/mbordner/aoc2023/common/graph"
	"github.com/mbordner/aoc2023/common/graph/djikstra"
	"log"
)

const (
	multiplier = 1000000
)

func expandSky(sky [][]rune) {
	for j := 0; j < len(sky); j++ {
		galaxy := false
		for i := 0; i < len(sky[j]); i++ {
			if sky[j][i] == '#' {
				galaxy = true
				break
			}
		}
		if !galaxy {
			for i := 0; i < len(sky[j]); i++ {
				if sky[j][i] == '.' {
					sky[j][i] = '1'
				} else {
					sky[j][i] = '2'
				}
			}
		}
	}

	for i := 0; i < len(sky[0]); i++ {
		galaxy := false
		for j := 0; j < len(sky); j++ {
			if sky[j][i] == '#' {
				galaxy = true
				break
			}
		}
		if !galaxy {
			for j := 0; j < len(sky); j++ {
				if sky[j][i] == '.' {
					sky[j][i] = '1'
				} else {
					sky[j][i] = '2'
				}
			}
		}
	}
}

func main() {
	sky := getData("../data.txt")
	expandSky(sky)

	g, bb, gm := getGraph(sky) // graph, bounding box, galaxy pos map
	if bb.GetPositionsSize() == 0 {
		log.Fatal("no nodes?")
	}

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

func getGraph(sky [][]rune) (*graph.Graph, geom.BoundingBox, graph.PosNodeMap[geom.Pos]) {
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
		cost := 1
		t := n.GetProperty("tile").(string)
		if t == "1" {
			cost = multiplier
		} else if t == "2" {
			cost = multiplier * 2
		}
		targets := make(map[geom.Direction]geom.Pos) // targets
		targets[geom.North] = p.Transform(0, -1, 0)
		targets[geom.West] = p.Transform(-1, 0, 0)
		targets[geom.East] = p.Transform(1, 0, 0)
		targets[geom.South] = p.Transform(0, 1, 0)
		for td, tp := range targets {
			if bb.Contains(tp) {
				on := g.GetNode(tp)
				e := n.AddEdge(on, float64(cost))
				e.AddProperty("dir", td)
			}
		}
	}

	return g, bb, mg
}

func getData(path string) [][]rune {
	lines, _ := file.GetLines(path)
	data := make([][]rune, len(lines), len(lines))
	for i, line := range lines {
		data[i] = []rune(line)
	}
	return data
}
