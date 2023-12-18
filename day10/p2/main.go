package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
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

var (
	vertical   = []geom.Direction{geom.North, geom.South}
	horizontal = []geom.Direction{geom.West, geom.East}
	mapTileDir = make(map[string]geom.Direction)
	mapDirTile = make(map[geom.Direction]string)
)

func main() {
	for _, b := range "|-LJ7F" {
		t := string(b)
		d := getDirs(t)
		mapTileDir[t] = d
		mapDirTile[d] = t
	}

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

	edges := start.GetEdges()
	d := geom.Unknown
	for _, e := range edges {
		dir := e.GetProperty("dir").(geom.Direction)
		d |= dir
	}
	start.AddProperty("tile", mapDirTile[d])
	fmt.Println(start.GetProperty("tile").(string))

	sp := djikstra.GenerateShortestPaths(g, start)

	maxDistance := 0
	var mp geom.Pos
	for k, nv := range sp {
		if nv.PreviousNode == nil {
			delete(sp, k)
		}
	}

	loopBB := geom.BoundingBox{}
	ps := start.GetID().(geom.Pos)
	loopBB.SetExtents(ps.X, ps.Y, ps.Z, ps.X, ps.Y, ps.Z)
	loopBB.Extend(ps)
	inner := make(map[geom.Pos]*graph.Node)
	loopPos := make(map[geom.Pos]*graph.Node)
	loopPos[ps] = start

	for k, nv := range sp {
		p := k.(geom.Pos)
		loopPos[p] = g.GetNode(p)
		loopBB.Extend(p)
		if int(nv.Value) > maxDistance {
			maxDistance = int(nv.Value)
			mp = p
		}
	}

	fmt.Println(mp, "in number of moves:", maxDistance)

	for j := loopBB.MinY; j <= loopBB.MaxY; j++ {
		outer := true
		for i := loopBB.MinX; i <= loopBB.MaxX; i++ {
			p := geom.Pos{Y: j, X: i}
			n := g.GetNode(p)
			tile := n.GetProperty("tile").(string)
			if _, ok := loopPos[p]; !ok {
				if !outer {
					inner[p] = n
					n.AddProperty("tile", "I")
				} else {
					n.AddProperty("tile", "O")
				}
			} else {
				if mapTileDir[tile].Is(vertical) {
					outer = !outer
					if mapTileDir[tile].Is(horizontal) {
						np := p
						for i < loopBB.MaxX {
							i++
							np = np.Transform(1, 0, 0)
							if _, ok := loopPos[np]; !ok {
								i--
								break
							}
							nn := g.GetNode(np)
							nt := nn.GetProperty("tile").(string)
							nd := getDirs(nt)
							if nd.Is(vertical) {
								if changingArea(tile, nt) {
									outer = !outer
								}
								break
							}
						}
					}
				}
			}
		}
	}

	pss := loopBB.GetPositions()
	chars := make([]rune, len(pss), len(pss))

	start.AddProperty("tile", "S")

	for i, p := range pss {
		n := g.GetNode(p)
		t := n.GetProperty("tile").(string)
		chars[i] = rune(t[0])
	}

	printLines := array.Reverse(loopBB.GetPrintLines('.', chars, pss))
	for _, line := range printLines {
		fmt.Println(line)
	}

	fmt.Println("nest spots:", len(inner))

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

func getDirs(tile string) geom.Direction {
	switch tile {
	case "|":
		return geom.North | geom.South
	case "-":
		return geom.West | geom.East
	case "L":
		return geom.North | geom.East
	case "J":
		return geom.North | geom.West
	case "7":
		return geom.South | geom.West
	case "F":
		return geom.South | geom.East
	}
	return geom.Unknown
}

func changingArea(c1, c2 string) bool {
	if c1 == "L" {
		if c2 == "7" {
			return false
		} else if c2 == "J" {
			return true
		}
	} else if c1 == "F" {
		if c2 == "7" {
			return true
		} else if c2 == "J" {
			return false
		}
	}
	return false
}
