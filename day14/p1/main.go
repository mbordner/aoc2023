package main

import (
	"fmt"

	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
)

const (
	roundedRock = 'O'
	empty       = '.'
	cubeRock    = '#'
)

type platformObj struct {
	rows [][]byte
}

func (p *platformObj) print() {
	for j := 0; j < len(p.rows); j++ {
		fmt.Println(string(p.rows[j]))
	}
}

func (p *platformObj) load(d geom.Direction) int {
	l := 0
	if d == geom.North {
		for i := 0; i < len(p.rows[0]); i++ {
			for j, t := len(p.rows)-1, 1; j >= 0; j, t = j-1, t+1 {
				if p.rows[j][i] == roundedRock {
					l += t
				}
			}
		}
	}
	return l
}

func (p *platformObj) tilt(d geom.Direction) {
	if d == geom.North {
		for i := 0; i < len(p.rows[0]); i++ {
			for j := 0; j < len(p.rows); j++ {
				b := p.rows[j][i]
				if b == roundedRock {
					cj := j
					pj := cj - 1
					for pj >= 0 && p.rows[pj][i] == empty {
						p.rows[pj][i], p.rows[cj][i] = p.rows[cj][i], p.rows[pj][i]
						cj--
						pj--
					}
				}
			}
		}
	}
}

func main() {
	p := getPlatform("../data.txt")

	p.tilt(geom.North)
	p.print()
	fmt.Println(p.load(geom.North))
}

func getPlatform(path string) *platformObj {
	p := new(platformObj)
	lines, _ := files.GetLines(path)
	p.rows = make([][]byte, len(lines))
	for i, line := range lines {
		p.rows[i] = []byte(line)
	}
	return p
}
