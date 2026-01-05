package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
)

func main() {
	nodes, edges := getData("../data.txt")

	var connectedNodeIDs []string
	for _, node := range nodes {
		connectedNodeIDs = node.GetConnectedNodeIDs()
		break
	}
	if len(connectedNodeIDs) != len(nodes) {
		panic("expected all nodes to be connected into 1 graph at start")
	}

	edgeLords := make(map[string]int)
	for _, edge := range edges {
		edgeLords[edge.ID()] = 0
	}

	pairs := common.GetPairSets(connectedNodeIDs)
	rand.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})

	times := min(10000, len(pairs))

	for i := 0; i < times; i++ {
		id1 := pairs[i][0]
		id2 := pairs[i][1]
		path := getPath(nodes, edges, id1, id2)
		for _, p := range path {
			edgeLords[p]++
		}
	}

	type EC struct {
		id string
		c  int
	}
	ecs := make([]EC, 0, len(edgeLords))
	for id, count := range edgeLords {
		ecs = append(ecs, EC{id: id, c: count})
	}
	sort.Slice(ecs, func(i, j int) bool {
		return ecs[i].c > ecs[j].c
	})

	for _, ec := range ecs[0:3] {
		edges[ec.id].SetConnected(false)
	}

	g1 := edges[ecs[0].id].n1.GetConnectedNodeIDs()
	g2 := edges[ecs[1].id].n2.GetConnectedNodeIDs()
	fmt.Println(len(g1) * len(g2))

}

func getPath(nodes Nodes, edges Edges, id1, id2 string) []string {
	visited := make(common.VisitedState[*Node, bool])
	queue := make(common.Queue[*Node], 0, len(nodes))
	start := nodes.Get(id1)
	goal := nodes.Get(id2)
	visited[start] = true
	prev := make(common.PreviousState[*Node, *Edge])
	queue.Enqueue(start)

	var pathEdgeIds []string

	for !queue.Empty() {
		cur := *(queue.Dequeue())

		if cur == goal {
			actions := prev.GetActions(start, goal)
			pathEdgeIds = make([]string, len(actions))
			for i, action := range actions {
				pathEdgeIds[i] = action.Action().ID()
			}
		} else {
			nextEdges := cur.ConnectedEdges()
			for _, nextEdge := range nextEdges {
				next := nextEdge.To(cur)
				if !visited.Has(next) {
					visited[next] = true
					queue.Enqueue(next)
					prev.Link(next, cur, nextEdge)
				}
			}
		}
	}
	return pathEdgeIds
}

type Edge struct {
	id        string
	n1        *Node
	n2        *Node
	connected bool
}

func (e *Edge) ID() string {
	return e.id
}

func (e *Edge) SetConnected(c bool) {
	e.connected = c
}

func (e *Edge) Connected() bool {
	return e.connected
}

func (e *Edge) otherNode(to *Node) *Node {
	if e.connected {
		if e.n1 == to {
			return e.n2
		}
		return e.n1
	}
	return nil
}

func (e *Edge) To(from *Node) *Node {
	return e.otherNode(from)
}

func (e *Edge) From(to *Node) *Node {
	return e.otherNode(to)
}

type Edges map[string]*Edge

func (edges Edges) Get(n1, n2 *Node) *Edge {
	nodes := []*Node{n1, n2}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].id < nodes[j].id
	})
	name := fmt.Sprintf("%s-%s", nodes[0].id, nodes[1].id)
	if _, exists := edges[name]; !exists {
		edges[name] = &Edge{n1: nodes[0], n2: nodes[1], connected: true, id: name}
	}
	return edges[name]
}

type Nodes map[string]*Node

func (nodes Nodes) Get(id string) *Node {
	if _, exists := nodes[id]; !exists {
		nodes[id] = &Node{id: id, edges: make(map[string]*Edge)}
	}
	return nodes[id]
}

type Node struct {
	id    string
	edges map[string]*Edge
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Connect(e *Edge) {
	n.edges[e.id] = e
}

func (n *Node) ConnectedEdges() []*Edge {
	edges := make([]*Edge, 0, len(n.edges))
	for _, edge := range n.edges {
		if edge.Connected() {
			edges = append(edges, edge)
		}
	}
	return edges
}

func (n *Node) getConnectedNodes(visited common.VisitedState[string, bool]) {
	visited[n.id] = true
	for _, edge := range n.edges {
		if edge.Connected() {
			o := edge.To(n)
			if !visited.Has(o.ID()) {
				o.getConnectedNodes(visited)
			}
		}
	}
}

func (n *Node) GetConnectedNodeIDs() []string {
	visited := make(common.VisitedState[string, bool])
	n.getConnectedNodes(visited)
	connectedNodes := make([]string, 0, len(visited))
	for id := range visited {
		connectedNodes = append(connectedNodes, id)
	}
	return connectedNodes
}

func getData(filename string) (Nodes, Edges) {
	nodes := make(Nodes)
	edges := make(Edges)

	replacer := strings.NewReplacer(": ", " ")

	lines := files.MustGetLines(filename)
	for _, line := range lines {
		fields := strings.Fields(replacer.Replace(line))
		n := nodes.Get(fields[0])
		for _, on := range fields[1:] {
			o := nodes.Get(on)
			e := edges.Get(n, o)
			n.Connect(e)
			o.Connect(e)
		}
	}

	return nodes, edges
}
