package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/files"
)

type conditionRecord struct {
	condition string
	checks    []int
}

func main() {
	crs := getData("../data.txt")

	sum := int64(0)

	for _, cr := range crs {
		cs := generatePossibleArrangements(cr.condition, cr.checks)
		if len(cs) == 0 {
			log.Fatalln("how?")
		}
		sum += int64(len(cs))
	}

	fmt.Println(sum)
}

func distribute(val, buckets int) [][]int {
	if buckets == 1 {
		return [][]int{{val}}
	}
	vals := make([][]int, 0, 20)
	for x := val; x >= 0; x-- {
		tvals := distribute(val-x, buckets-1)
		for _, vs := range tvals {
			vals = append(vals, append([]int{x}, vs...))
		}
	}
	return vals
}

func generatePossibleArrangements(condition string, checks []int) []string {
	cs := make([]string, 0, 20)

	s1 := array.SumNumbers(checks)
	s2 := len(checks) - 1

	distributeVal := len(condition) - (s1 + s2)
	if distributeVal < 0 {
		log.Fatalln("this can't happen")
	}

	reStr := strings.ReplaceAll(condition, ".", `\.`)
	reStr = strings.ReplaceAll(reStr, "?", `[\.#]`)
	reCondition := regexp.MustCompile(reStr)

	dvs := distribute(distributeVal, s2+2)
	for _, dv := range dvs {
		bs := make([]byte, len(condition))
		p := 0
		for i := 0; i < len(dv); i++ {
			dots := dv[i]
			if i > 0 && i < len(dv)-1 {
				dots++
			}
			for c := 0; c < dots; c, p = c+1, p+1 {
				bs[p] = '.'
			}
			if i < len(checks) {
				pounds := checks[i]
				for c := 0; c < pounds; c, p = c+1, p+1 {
					bs[p] = '#'
				}
			}
		}
		s := string(bs)
		if reCondition.MatchString(s) {
			cs = append(cs, string(bs))
		}
	}

	return cs
}

func getData(path string) []*conditionRecord {
	lines, _ := files.GetLines(path)
	crs := make([]*conditionRecord, len(lines))
	for i, line := range lines {
		cr := conditionRecord{}
		tokens := strings.Split(line, " ")
		cr.condition = tokens[0]
		cr.checks = array.ToNumbers[int](tokens[1], ",")
		crs[i] = &cr
	}
	return crs
}
