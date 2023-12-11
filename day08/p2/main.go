package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/file"
	"regexp"
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
type ghostsObj []string

func (g ghostsObj) isEnd() bool {
	for _, v := range g {
		if !reEndNode.MatchString(v) {
			return false
		}
	}
	return true
}

func main() {
	dirs, nodes := getData("../data.txt")
	if dirs != nil && nodes != nil {
		ghosts := make(ghostsObj, 0, 10)
		for k := range nodes {
			if reStartNode.MatchString(k) {
				ghosts = append(ghosts, k)
			}
		}

		steps := int64(0)

		for ghosts.isEnd() == false {
			d := dirs.next()
			for i := range ghosts {
				ghosts[i] = nodes[ghosts[i]].next(d)
			}
			steps++
		}

		fmt.Println(steps)
	}
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
