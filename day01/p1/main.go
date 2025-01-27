package main

import (
	"fmt"

	"github.com/mbordner/aoc2023/common/files"
)

func main() {
	lines := getLines("../data.txt")
	sum := 0
	for _, line := range lines {
		val := 0
		for i := 0; i < len(line); i++ {
			if v := digit(line[i]); v >= 0 {
				val += 10 * v
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if v := digit(line[i]); v >= 0 {
				val += v
				break
			}
		}
		sum += val
	}
	fmt.Println(sum)
}

func digit(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	return -1
}

func getLines(path string) []string {
	lines, _ := files.GetLines(path)
	return lines
}
