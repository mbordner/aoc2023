package main

import (
	"cmp"
	"fmt"
	"github.com/mbordner/aoc2023/common/file"
	"regexp"
	"slices"
)

var (
	reNode      = regexp.MustCompile(`^(\w+)\s*=\s*\((\w+)\s*,\s*(\w+)\)$`)
	reStartNode = regexp.MustCompile(`..A`)
	reEndNode   = regexp.MustCompile(`..Z`)
)

type directionsObj struct {
	dir   string
	index int
}

func newDirections(dir string) *directionsObj {
	d := new(directionsObj)
	d.dir = dir
	return d
}

func (d *directionsObj) next() (string, int) {
	index := d.index
	nd := string(d.dir[index])
	d.index++
	if d.index == len(d.dir) {
		d.index = 0
	}
	return nd, index
}

func (d *directionsObj) len() int {
	return len(d.dir)
}

func (d *directionsObj) reset() {
	d.index = 0
}

type lrObj struct {
	l string
	r string
}

func (lr lrObj) next(dir string) string {
	if dir == "R" {
		return lr.r
	}
	return lr.l
}

type nodesObj map[string]lrObj
type ghostsObj []string

func (g ghostsObj) isEnd() bool {
	for _, v := range g {
		if !reEndNode.MatchString(v) {
			return false
		}
	}
	return true
}

type visitedObj struct {
	node     string
	dirIndex int
}

type visitedCounts map[visitedObj][]int

func (vc visitedCounts) add(node string, dirIndex int, pathIndex int) []int {
	vo := visitedObj{
		node:     node,
		dirIndex: dirIndex,
	}
	if _, e := vc[vo]; !e {
		vc[vo] = make([]int, 0, 10)
	}
	vc[vo] = append(vc[vo], pathIndex)
	return vc[vo]
}

type pathStatsObj struct {
	start          string
	goal           string
	loopStartIndex int
	goalIndexes    []int
	loopEndIndex   int
	goalIndex      int
}

func (ps pathStatsObj) String() string {
	return fmt.Sprintf("[%s-%s: %d, %d, %d]", ps.start, ps.goal, ps.loopStartIndex, ps.loopEndIndex, ps.goalIndex)
}

func main() {
	dirs, nodes := getData("../data.txt")
	ghosts := make(ghostsObj, 0, 10)
	for k := range nodes {
		if reStartNode.MatchString(k) {
			ghosts = append(ghosts, k)
		}
	}

	ps := make([]pathStatsObj, len(ghosts), len(ghosts))
	for g := range ghosts {
		ps[g].start = ghosts[g]
	}

	for i := range ghosts {
		dirs.reset()
		pathIndex := 0
		visited := make(visitedCounts)
		var vc []int

		visited.add(ghosts[i], dirs.len()-1, pathIndex)

		for len(vc) < 2 {
			pathIndex++
			d, di := dirs.next()
			ghosts[i] = nodes[ghosts[i]].next(d)
			if reEndNode.MatchString(ghosts[i]) {
				if ps[i].goalIndexes == nil {
					ps[i].goalIndexes = make([]int, 0, 10)
				}
				ps[i].goal = ghosts[i]
				ps[i].goalIndexes = append(ps[i].goalIndexes, pathIndex)
			}
			vc = visited.add(ghosts[i], di, pathIndex)
		}

		ps[i].loopStartIndex = vc[0]
		ps[i].loopEndIndex = vc[1]
	}

	for i := 0; i < len(ps); i++ {
		ps[i].goalIndex = ps[i].goalIndexes[0]
		if len(ps[i].goalIndexes) > 1 {
			ps[i].loopEndIndex = ps[i].loopEndIndex - (ps[i].goalIndexes[1] - ps[i].goalIndexes[0])
		}
	}

	slices.SortFunc(ps, func(a, b pathStatsObj) int {
		return cmp.Compare(a.start, b.start)
	})

	fmt.Println(ps)

	firsts := make([]int64, len(ps), len(ps))
	nexts := make([]int64, len(ps), len(ps))

	for i := 0; i < len(ps); i++ {

		firsts[i] = int64(ps[i].loopStartIndex) +
			int64(ps[i].goalIndex-ps[i].loopStartIndex)
		nexts[i] = (firsts[i] +
			int64(1)*int64(ps[i].loopEndIndex-ps[i].loopStartIndex)) - firsts[i]

	}

	vals := make([]int64, len(ps), len(ps))
	copy(vals, firsts)

	max := int64(0)
	for i := 0; i < len(vals); i++ {
		if max < vals[i] {
			max = vals[i]
		}
	}

	changed := true
	for changed {
		changed = false
		for i := 0; i < len(vals); i++ {
			for vals[i] < max {
				vals[i] += nexts[i]
				if vals[i] > max {
					changed = true
					max = vals[i]
				}
			}
		}
	}
	fmt.Println(vals)
}

func getData(path string) (*directionsObj, nodesObj) {
	lines, _ := file.GetLines(path)
	dirs := newDirections(lines[0])
	nodes := make(nodesObj)
	for i := 2; i < len(lines); i++ {
		matches := reNode.FindStringSubmatch(lines[i])
		if matches != nil {
			nodes[matches[1]] = lrObj{
				l: matches[2],
				r: matches[3],
			}
		}
	}
	return dirs, nodes
}
