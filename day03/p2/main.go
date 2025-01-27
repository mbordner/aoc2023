package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mbordner/aoc2023/common/files"
)

type pos struct {
	j int
	i int
}

var (
	reDigits = regexp.MustCompile(`\d+`)
	reDigit  = regexp.MustCompile(`\d|\.`)
	stars    = make(map[pos]map[pos]int64)
)

func checkPart(lines []string, num pos, val int64, symbol pos) bool {
	if !reDigit.MatchString(string(lines[symbol.j][symbol.i])) {
		if lines[symbol.j][symbol.i] == '*' {
			if _, ok := stars[symbol]; !ok {
				stars[symbol] = make(map[pos]int64)
			}
			stars[symbol][num] = val
		}
		return true
	}
	return false
}

func main() {
	lines, _ := files.GetLines("../data.txt")
	sum := int64(0)
	for j, line := range lines {
		matches := reDigits.FindAllStringIndex(line, -1)
		if matches != nil {
			for _, match := range matches {
				part := false
				num, _ := strconv.ParseInt(line[match[0]:match[1]], 10, 64)
				numPos := pos{i: match[0], j: j}
				for i := match[0] - 1; i <= match[1]; i++ {
					if i >= 0 && i < len(line) {
						if j > 0 {
							if checkPart(lines, numPos, num, pos{i: i, j: j - 1}) {
								part = true
							}
						}
						if j < len(lines)-1 {
							if checkPart(lines, numPos, num, pos{i: i, j: j + 1}) {
								part = true
							}
						}
						if i == match[0]-1 || i == match[1] {
							if checkPart(lines, numPos, num, pos{i: i, j: j}) {
								part = true
							}
						}
					}
				}
				if part {
					sum += num
				}
			}
		}
	}
	fmt.Println(sum)

	sum = 0

	for _, vs := range stars {
		if len(vs) == 2 {
			val := int64(1)
			for _, v := range vs {
				val *= v
			}
			sum += val
		}
	}

	fmt.Println(sum)
}
