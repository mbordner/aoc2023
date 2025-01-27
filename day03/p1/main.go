package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mbordner/aoc2023/common/files"
)

var (
	reDigits = regexp.MustCompile(`\d+`)
	reDigit  = regexp.MustCompile(`\d|\.`)
)

func main() {
	lines, _ := files.GetLines("../data.txt")
	sum := int64(0)
	for j, line := range lines {
		matches := reDigits.FindAllStringIndex(line, -1)
		if matches != nil {
			for _, match := range matches {
				part := false
				num, _ := strconv.ParseInt(line[match[0]:match[1]], 10, 64)
				for i := match[0] - 1; i <= match[1]; i++ {
					if i >= 0 && i < len(line) {
						if j > 0 {
							if !reDigit.MatchString(string(lines[j-1][i])) {
								part = true
								break
							}
						}
						if j < len(lines)-1 {
							if !reDigit.MatchString(string(lines[j+1][i])) {
								part = true
								break
							}
						}
						if i == match[0]-1 || i == match[1] {
							if !reDigit.MatchString(string(lines[j][i])) {
								part = true
								break
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
}
