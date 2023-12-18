package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/file"
	"github.com/mbordner/aoc2023/common/geom"
	"github.com/mbordner/aoc2023/common/graph"
	"github.com/mbordner/aoc2023/common/graph/djikstra"
)

/*
Pipes:
| is a vertical pipe connecting north and south.
- is a horizontal pipe connecting east and west.
L is a 90-degree bend connecting north and east.
J is a 90-degree bend connecting north and west.
7 is a 90-degree bend connecting south and west.
F is a 90-degree bend connecting south and east.

. is ground
S is start
*/

func main() {

	lines, _ := file.GetLines("../data.txt")
	g := graph.NewGraph()
	bb := geom.BoundingBox{}
	var start *graph.Node

	// create nodes
	for j := 0; j < len(lines); j++ {
		for i, b := range lines[j] {
			p := geom.Pos{Y: j, X: i}
			n := g.CreateNode(p)
			tile := string(b)
			if tile == "S" {
				start = n
			} else if tile == "." {
				n.SetTraversable(false)
			}
			n.AddProperty("tile", tile)
			bb.Extend(p)
		}
	}

	// connect nodes
	for j := 0; j < len(lines); j++ {
		for i, b := range lines[j] {
			tile := string(b)
			p := geom.Pos{Y: j, X: i}
			np := g.GetNode(p)

			cp := getConnectingPos(tile, p)

			for d, op := range cp {
				if bb.Contains(op) {
					on := g.GetNode(op)
					ot := on.GetProperty("tile").(string)
					ocp := getConnectingPos(ot, op)
					if _, ok := ocp[d.Opposite()]; ok {
						e := np.AddEdge(on, 1)
						e.AddProperty("dir", d)
					}

					if ot == "S" {
						e := start.AddEdge(np, 1)
						e.AddProperty("dir", d.Opposite())
					}
				}
			}
		}
	}

	fmt.Println(start.GetProperty("tile").(string))

	sp := djikstra.GenerateShortestPaths(g, start)

	maxDistance := 0
	var mp geom.Pos
	for k, nv := range sp {
		if nv.PreviousNode != nil {
			if int(nv.Value) > maxDistance {
				maxDistance = int(nv.Value)
				mp = k.(geom.Pos)
			}
		}
	}

	fmt.Println(mp, "in number of moves:", maxDistance)
}

func getConnectingPos(tile string, p geom.Pos) map[geom.Direction]geom.Pos {
	cp := make(map[geom.Direction]geom.Pos) // pos to connect
	switch tile {
	case "|":
		cp[geom.North] = p.Transform(0, -1, 0)
		cp[geom.South] = p.Transform(0, 1, 0)
	case "-":
		cp[geom.West] = p.Transform(-1, 0, 0)
		cp[geom.East] = p.Transform(1, 0, 0)
	case "L":
		cp[geom.North] = p.Transform(0, -1, 0)
		cp[geom.East] = p.Transform(1, 0, 0)
	case "J":
		cp[geom.North] = p.Transform(0, -1, 0)
		cp[geom.West] = p.Transform(-1, 0, 0)
	case "7":
		cp[geom.South] = p.Transform(0, 1, 0)
		cp[geom.West] = p.Transform(-1, 0, 0)
	case "F":
		cp[geom.South] = p.Transform(0, 1, 0)
		cp[geom.East] = p.Transform(1, 0, 0)
	}
	return cp
}
