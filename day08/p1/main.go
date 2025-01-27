package main

import (
	"fmt"
	"regexp"

	"github.com/mbordner/aoc2023/common/files"
)

const (
	startNode = "AAA"
	endNode   = "ZZZ"
)

var (
	reNode = regexp.MustCompile(`^(\w+)\s*=\s*\((\w+)\s*,\s*(\w+)\)$`)
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

func (d *directionsObj) next() string {
	nd := string(d.dir[d.index])
	d.index++
	if d.index == len(d.dir) {
		d.index = 0
	}
	return nd
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

func main() {
	dirs, nodes := getData("../data.txt")
	if dirs != nil && nodes != nil {
		cur := startNode
		steps := 0
		for cur != endNode {
			d := dirs.next()
			next := nodes[cur].next(d)
			steps++
			cur = next
		}
		fmt.Println(steps)
	}
}

func getData(path string) (*directionsObj, nodesObj) {
	lines, _ := files.GetLines(path)
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
