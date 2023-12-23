package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/file"
	"log"
)

type dps map[string]map[int]bool

func main() {
	data := getData("../data.txt")
	sum := 0
	for _, p := range data {
		sum += summarizePattern(p)
	}
	fmt.Println(sum)
}

func summarizePattern(p [][]byte) int {
	s, s2 := -1, -1
	dp := make(dps)

	q := array.SwapRowCols(p)

	v := findReflectCol(p, -1, dp)
	if v >= 0 {
		s = v + 1
	} else {
		v = findReflectCol(q, -1, dp)
		if v < 0 {
			log.Fatalln("couldn't find reflection")
		}
		s = (v + 1) * 100
	}

	for j := 0; j < len(p) && s2 < 0; j++ {
		for i := 0; i < len(p[j]) && s2 < 0; i++ {
			oc := p[j][i]
			var nc byte
			if oc == '#' {
				nc = '.'
			} else {
				nc = '#'
			}
			p[j][i] = nc
			q[i][j] = nc

			v2 := findReflectCol(p, s%100-1, dp)
			if v2 >= 0 {
				s2 = v2 + 1
			} else {
				v2 = findReflectCol(q, s/100-1, dp)
				if v2 >= 0 {
					s2 = (v2 + 1) * 100
				}
			}

			p[j][i] = oc
			q[i][j] = oc
		}
	}

	return s2
}

func findReflectCol(p [][]byte, not int, dp dps) int {
outer:
	for i := 0; i < len(p[0]); i++ {
		if i == not {
			continue
		}
		if check(p[0], i, dp) {
			for j := 1; j < len(p); j++ {
				if !check(p[j], i, dp) {
					continue outer
				}
			}
			return i
		}
	}
	return -1
}

func check(bs []byte, c int, dp dps) bool {
	k := string(bs)
	if m, e := dp[k]; e {
		if b, e := m[c]; e {
			return b
		}
	} else {
		dp[k] = make(map[int]bool)
	}

	dp[k][c] = true

	if c == len(bs)-1 {
		dp[k][c] = false
	} else {
		for i, j := c, c+1; i >= 0 && j < len(bs); i, j = i-1, j+1 {
			if bs[i] != bs[j] {
				dp[k][c] = false
				break
			}
		}
	}

	return dp[k][c]
}

func getData(path string) [][][]byte {
	lines, _ := file.GetLines(path)
	data := make([][][]byte, 0)

	pattern := make([][]byte, 0)
	for _, line := range lines {
		if line == "" {
			data = append(data, pattern)
			pattern = make([][]byte, 0)
			continue
		}
		bytes := []byte(line)
		pattern = append(pattern, bytes)
	}

	data = append(data, pattern)

	return data
}
