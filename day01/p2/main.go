package main

import (
	"fmt"
	"regexp"

	"github.com/mbordner/aoc2023/common/files"
)

var (
	reDigit = regexp.MustCompile(`(\d|one|two|three|four|five|six|seven|eight|nine)`)
)

func main() {
	lines := getLines("../data.txt")
	sum := 0
	for _, line := range lines {
		val := 0
		matches := reDigit.FindAllStringSubmatch(line, -1)
		val += digit(matches[0][0]) * 10
		for i := len(line) - 1; i >= 0; i-- {
			match := reDigit.FindStringSubmatch(line[i:])
			if match != nil {
				val += digit(match[0])
				break
			}
		}
		sum += val
	}
	fmt.Println(sum)
}

func digit(s string) int {
	if len(s) == 1 {
		return int(s[0] - '0')
	}
	switch s {
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	}
	return 9
}

func getLines(path string) []string {
	lines, _ := files.GetLines(path)
	return lines
}
