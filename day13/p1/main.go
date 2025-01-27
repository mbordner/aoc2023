package main

import (
	"fmt"
	"log"

	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/files"
)

func main() {
	data := getData("../data.txt")
	sum := 0
	for _, p := range data {
		sum += summarizePattern(p)
	}
	fmt.Println(sum)
}

func summarizePattern(p [][]byte) int {
	s := 0
	v := findReflectCol(p)
	if v >= 0 {
		s += v + 1
	} else {
		q := array.SwapRowCols(p)
		v = findReflectCol(q)
		if v < 0 {
			log.Fatalln("couldn't find reflection")
		}
		s += (v + 1) * 100
	}
	return s
}

func findReflectCol(p [][]byte) int {
outer:
	for i := 0; i < len(p[0]); i++ {
		if check(p[0], i) {
			for j := 1; j < len(p); j++ {
				if !check(p[j], i) {
					continue outer
				}
			}
			return i
		}
	}
	return -1
}

func check(bs []byte, c int) bool {
	if c == len(bs)-1 {
		return false
	}
	for i, j := c, c+1; i >= 0 && j < len(bs); i, j = i-1, j+1 {
		if bs[i] != bs[j] {
			return false
		}
	}

	return true
}

func getData(path string) [][][]byte {
	lines, _ := files.GetLines(path)
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
