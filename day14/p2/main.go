package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/file"
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
	} else if d == geom.South {
		for i := 0; i < len(p.rows[0]); i++ {
			for j, t := 0, 1; j < len(p.rows); j, t = j+1, t+1 {
				if p.rows[j][i] == roundedRock {
					l += t
				}
			}
		}
	} else if d == geom.East {
		for j := 0; j < len(p.rows); j++ {
			for i, t := 0, 1; i < len(p.rows[j]); i, t = i+1, t+1 {
				if p.rows[j][i] == roundedRock {
					l += t
				}
			}
		}
	} else if d == geom.West {
		for j := 0; j < len(p.rows); j++ {
			for i, t := len(p.rows[j])-1, 1; i >= 0; i, t = i-1, t+1 {
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
	} else if d == geom.South {
		for i := 0; i < len(p.rows[0]); i++ {
			for j := len(p.rows) - 1; j >= 0; j-- {
				b := p.rows[j][i]
				if b == roundedRock {
					cj := j
					pj := cj + 1
					for pj < len(p.rows) && p.rows[pj][i] == empty {
						p.rows[pj][i], p.rows[cj][i] = p.rows[cj][i], p.rows[pj][i]
						cj++
						pj++
					}
				}
			}
		}
	} else if d == geom.East {
		for j := 0; j < len(p.rows); j++ {
			for i := len(p.rows[j]) - 1; i >= 0; i-- {
				b := p.rows[j][i]
				if b == roundedRock {
					ci := i
					pi := ci + 1
					for pi < len(p.rows[j]) && p.rows[j][pi] == empty {
						p.rows[j][pi], p.rows[j][ci] = p.rows[j][ci], p.rows[j][pi]
						ci++
						pi++
					}
				}
			}
		}
	} else if d == geom.West {
		for j := 0; j < len(p.rows); j++ {
			for i := 0; i < len(p.rows[j]); i++ {
				b := p.rows[j][i]
				if b == roundedRock {
					ci := i
					pi := ci - 1
					for pi >= 0 && p.rows[j][pi] == empty {
						p.rows[j][pi], p.rows[j][ci] = p.rows[j][ci], p.rows[j][pi]
						ci--
						pi--
					}
				}
			}
		}
	}
}

func (p *platformObj) cycle() {
	p.tilt(geom.North)
	p.tilt(geom.West)
	p.tilt(geom.South)
	p.tilt(geom.East)
}

func main() {
	p := getPlatform("../data.txt")

	loads := make([]int, 0, 100)
	cycles := make(map[int][]int)
	skipped := false

	numCycles := uint64(1000000000)
	for i := uint64(0); i < numCycles; i++ {
		p.cycle()
		l := p.load(geom.East)
		loads = append(loads, l)
		li := len(loads) - 1
		if _, e := cycles[l]; !e {
			cycles[l] = make([]int, 0, 100)
		}
		cycles[l] = append(cycles[l], li)
		if !skipped && len(cycles[l]) > 2 {
			lastCycleEnd := cycles[l][len(cycles[l])-1]
			lastCycleStart := cycles[l][len(cycles[l])-2]
			prevCycleEnd := lastCycleStart
			prevCycleStart := cycles[l][len(cycles[l])-3]
			lastCycle := loads[lastCycleStart:lastCycleEnd]
			prevCycle := loads[prevCycleStart:prevCycleEnd]
			if array.Equals(prevCycle, lastCycle) {
				skip := (numCycles - i) / uint64(len(lastCycle))
				i += skip * uint64(len(lastCycle))
				skipped = true
			}
		}
	}

	fmt.Println(p.load(geom.North))

}

func getPlatform(path string) *platformObj {
	p := new(platformObj)
	lines, _ := file.GetLines(path)
	p.rows = make([][]byte, len(lines))
	for i, line := range lines {
		p.rows[i] = []byte(line)
	}
	return p
}
