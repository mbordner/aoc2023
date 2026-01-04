package main

import (
	"fmt"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
)

const (
	Forest     = '#'
	Path       = '.'
	SlopeEast  = '>'
	SlopeNorth = '^'
	SlopeWest  = '<'
	SlopeSouth = 'v'
	Step       = 'O'
	Clear      = ' '
)

func main() {
	grid, nodes, start, goal := getGrid("../data.txt")

	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == Forest {
				grid[y][x] = Clear
			}
		}
	}
	for p := range nodes {
		grid[p.Y][p.X] = Step
	}

	grid.Print()

	path := getLongestPath(start, goal, make(common.PosContainer))
	fmt.Println(path.Len())
}

type NodePath []*Node

func (p NodePath) Len() int {
	l := 0
	for i := 0; i < len(p)-1; i++ {
		l += p[i].edges[p[i+1]]
	}
	return l
}

func getLongestPath(n *Node, goal *Node, visited common.PosContainer) NodePath {
	var longest NodePath
	var longestDistance int

	visited[n.pos] = true

	if n == goal {
		longest = NodePath{goal}
	} else {
		for o := range n.edges {
			if !visited.Has(o.pos) {
				path := getLongestPath(o, goal, visited)
				if len(path) > 0 {
					path = append(NodePath{n}, path...)
					pl := path.Len()
					if pl > longestDistance {
						longestDistance = pl
						longest = path
					}
				}
			}
		}
	}

	visited[n.pos] = false

	return longest
}

type Edges map[*Node]int

type Node struct {
	pos     common.Pos
	edges   Edges
	visited common.PosContainer
}

func NewNode(pos common.Pos) *Node {
	return &Node{pos: pos, edges: make(Edges), visited: make(common.PosContainer)}
}

func (n *Node) Pos() common.Pos {
	return n.pos
}

func (n *Node) AddEdge(o *Node, distance int) {
	if d, e := n.edges[o]; !e || distance > d {
		n.edges[o] = distance
	}
}

func (n *Node) AddVisited(p common.Pos) {
	n.visited[p] = true
}

func (n *Node) Visited(p common.Pos) bool {
	return n.visited.Has(p)
}

type Nodes map[common.Pos]*Node

func (ns Nodes) Get(pos common.Pos) *Node {
	if _, e := ns[pos]; !e {
		ns[pos] = NewNode(pos)
		ns[pos].AddVisited(pos)
	}
	return ns[pos]
}

func findNodes(grid common.Grid, nodes Nodes, cur *Node, cp common.Pos, pp common.Pos, distance int, goal *Node) {
	cur.AddVisited(cp)
	if cp == goal.Pos() {
		cur.AddEdge(goal, distance)
		goal.AddEdge(cur, distance)
		goal.AddVisited(pp)
	} else {
		nps := make(common.Positions, 0, 4)
		h, w := len(grid), len(grid[0])
		for _, dir := range common.AdjacentDirs {
			np := cp.Add(dir)
			if np.Y >= 0 && np.Y < h && np.X >= 0 && np.X < w {
				if !cur.Visited(np) && grid[np.Y][np.X] != Forest {
					nps = append(nps, np)
				}
			}
		}
		if len(nps) > 0 {
			if len(nps) == 1 {
				findNodes(grid, nodes, cur, nps[0], cp, distance+1, goal)
			} else {
				nn := nodes.Get(cp)
				nn.AddEdge(cur, distance)
				cur.AddEdge(nn, distance)
				nn.AddVisited(pp)
				for _, np := range nps {
					if !nn.Visited(np) {
						findNodes(grid, nodes, nn, np, cp, 1, goal)
					}
				}
			}
		}
	}
}

func getGrid(filename string) (common.Grid, Nodes, *Node, *Node) {
	grid := common.ConvertGrid(files.MustGetLines(filename))
	var sp, gp common.Pos
	for x := 0; x < len(grid[0]); x++ {
		if grid[0][x] == Path {
			sp.X = x
			break
		}
	}
	gp.Y = len(grid) - 1
	for x := 0; x < len(grid[gp.Y]); x++ {
		if grid[gp.Y][x] == Path {
			gp.X = x
			break
		}
	}

	nodes := make(Nodes)
	start := nodes.Get(sp)
	goal := nodes.Get(gp)
	findNodes(grid, nodes, start, sp, sp, 0, goal)
	return grid, nodes, start, goal
}
