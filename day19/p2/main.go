package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/ranges"
)

var (
	reWorkflowLine      = regexp.MustCompile(`^(\w+)\{(.*)\}$`)
	reRuleWithCondition = regexp.MustCompile(`(x|m|a|s)\s*(>|<)\s*(\d+)\s*:\s*(\w+)`)
	reAcceptRejectRule  = regexp.MustCompile(`(A|R)`)
	reGotoRule          = regexp.MustCompile(`(\w+)`)
	rePartLine          = regexp.MustCompile(`^\{x=(\d+)\s*,\s*m=(\d+)\s*,\s*a=(\d+)\s*,\s*s=(\d+)\}$`)
)

const (
	MinVal = 1
	MaxVal = 4000
)

func main() {
	workflows := getData("../data.txt")

	nodes := make(Nodes)

	buildTreeFrom("in", workflows, nodes)

	accepted := nodes.Get("A")
	sum := uint64(0)
	for _, edge := range accepted.From {
		sum += edge.Ratings.Len()
		fmt.Printf("%s -> %s (%d) : %s\n", edge.From.ID, edge.To.ID, edge.Ratings.Len(), edge.Ratings)
	}
	fmt.Println(sum)
}

func buildTreeFrom(id string, workflows Workflows, nodes Nodes) {
	if slices.Contains([]string{"A", "R"}, id) {
		// terminal
	} else {
		cur := nodes.Get(id)
		ratings := cur.GetRatingRanges()
		for _, rule := range workflows[id].Rules {
			if !reRuleWithCondition.MatchString(string(rule)) && (reAcceptRejectRule.MatchString(string(rule)) || reGotoRule.MatchString(string(rule))) {
				next := nodes.Get(string(rule))
				edge := NewEdge(ratings, cur, next)
				next.AddFrom(edge)
				cur.AddTo(edge)
			} else if reRuleWithCondition.MatchString(string(rule)) {

				var nextRatings *XMAS
				matches := reRuleWithCondition.FindStringSubmatch(string(rule))
				val, _ := strconv.Atoi(matches[3])
				var l1, r1, l2, r2 int
				if matches[2] == "<" {
					l1 = val
					r1 = MaxVal
					l2 = MinVal
					r2 = l1 - 1
				} else {
					l1 = MinVal
					r1 = val
					l2 = val + 1
					r2 = MaxVal
				}

				nextRatings = ratings.Remove(matches[1], l1, r1)
				ratings = ratings.Remove(matches[1], l2, r2)

				next := nodes.Get(matches[4])
				edge := NewEdge(nextRatings, cur, next)
				next.AddFrom(edge)
				cur.AddTo(edge)

			}
		}

		for _, next := range cur.To {
			buildTreeFrom(next.To.ID, workflows, nodes)
		}
	}
}

type Nodes map[string]*Node

func (nodes Nodes) Get(id string) *Node {
	if n, e := nodes[id]; e {
		return n
	}
	nodes[id] = NewNode(id)
	return nodes[id]
}

type Edge struct {
	Ratings *XMAS
	To      *Node
	From    *Node
}

func NewEdge(ratings *XMAS, from, to *Node) *Edge {
	return &Edge{Ratings: ratings, From: from, To: to}
}

type Node struct {
	ID   string
	To   []*Edge
	From []*Edge
}

func (n *Node) GetRatingRanges() *XMAS {
	if len(n.From) == 0 {
		return NewXMAS()
	}
	if len(n.From) > 1 {
		panic("too many edges from nodes")
	}
	return n.From[0].Ratings
}

func (n *Node) AddTo(e *Edge) {
	n.To = append(n.To, e)
}

func (n *Node) AddFrom(e *Edge) {
	n.From = append(n.From, e)
}

func NewNode(id string) *Node {
	return &Node{ID: id, To: make([]*Edge, 0, 4), From: make([]*Edge, 0, 4)}
}

type XMAS struct {
	X *ranges.Collection[int]
	M *ranges.Collection[int]
	A *ranges.Collection[int]
	S *ranges.Collection[int]
}

func (xmas *XMAS) Len() uint64 {
	sum := uint64(1)
	sum *= uint64(xmas.X.Len())
	sum *= uint64(xmas.M.Len())
	sum *= uint64(xmas.A.Len())
	sum *= uint64(xmas.S.Len())
	return sum
}

func (xmas *XMAS) String() string {
	return fmt.Sprintf("XMAS{X:%v, M:%v, A:%v, S:%v}", xmas.X, xmas.M, xmas.A, xmas.S)
}

func (xmas *XMAS) Clone() *XMAS {
	oxmas := &XMAS{}
	oxmas.X = xmas.X.Clone()
	oxmas.M = xmas.M.Clone()
	oxmas.A = xmas.A.Clone()
	oxmas.S = xmas.S.Clone()
	return oxmas
}

func (xmas *XMAS) Remove(c string, l, r int) *XMAS {
	oxmas := xmas.Clone()
	switch c {
	case "x":
		_, _ = oxmas.X.Remove(l, r)
	case "m":
		_, _ = oxmas.M.Remove(l, r)
	case "a":
		_, _ = oxmas.A.Remove(l, r)
	case "s":
		_, _ = oxmas.S.Remove(l, r)
	}
	return oxmas
}

func NewXMAS() *XMAS {
	xmas := &XMAS{}
	xmas.X = &ranges.Collection[int]{}
	xmas.X.Add(MinVal, MaxVal)
	xmas.M = &ranges.Collection[int]{}
	xmas.M.Add(MinVal, MaxVal)
	xmas.A = &ranges.Collection[int]{}
	xmas.A.Add(MinVal, MaxVal)
	xmas.S = &ranges.Collection[int]{}
	xmas.S.Add(MinVal, MaxVal)
	return xmas
}

type Workflows map[string]*Workflow

type Rule string

type Rules []Rule

type Workflow struct {
	ID    string
	Rules Rules
}

func getData(filename string) Workflows {
	lines := files.MustGetLines(filename)

	workflows := make(Workflows)
	for l, line := range lines {
		if line == "" {
			lines = lines[l+1:]
			break
		} else {
			if reWorkflowLine.MatchString(line) {
				matches := reWorkflowLine.FindStringSubmatch(line)
				workflow := Workflow{ID: matches[1]}
				rules := strings.Split(matches[2], ",")
				workflow.Rules = make(Rules, len(rules))

				for r, rule := range rules {
					workflow.Rules[r] = Rule(rule)
				}

				workflows[workflow.ID] = &workflow
			} else {
				panic("invalid line: " + line)
			}
		}
	}

	return workflows
}
