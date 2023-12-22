package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/cmath"
	"github.com/mbordner/aoc2023/common/file"
	"log"
	"regexp"
	"strings"
)

type conditionRecord struct {
	condition string
	checks    []int
	vals      []uint64
}

// 1067550648907 too low
func main() {
	crs := getData("../data.txt")

	sum := uint64(0)

	for _, cr := range crs {
		cs := generatePossibleArrangements(cr.condition, cr.checks)
		val := uint64(len(cs))
		cr.vals = append(cr.vals, val)

		condition2 := cr.condition + "?" + cr.condition
		checks2 := append(array.Clone(cr.checks), cr.checks...)

		cs = generatePossibleArrangements(condition2, checks2)
		val = uint64(len(cs))
		cr.vals = append(cr.vals, val)

		mul := cr.vals[1] / cr.vals[0]
		for i := 0; i < 3; i++ {
			val *= mul
		}

		sum += val
	}

	fmt.Println(sum)
}

// n identical objects into k distinct bins,     (n+k-1)! / ((k-1)! * n!) arrangements
// see stars and bars, how many bars needed is k-1, add that to n, and do combination
func distributeCombinations(val, buckets int) int64 {
	n := val
	k := buckets
	return cmath.Factorial(int64(n+k-1)) / (cmath.Factorial(int64(k-1)) * cmath.Factorial(int64(n)))
}

type dKey struct {
	n int
	k int
}

var (
	dDP = make(map[dKey][][]int)
)

func distribute(val, buckets int) [][]int {
	if buckets == 1 {
		return [][]int{{val}}
	}
	k := dKey{n: val, k: buckets}
	if a, e := dDP[k]; e {
		return a
	}
	vals := make([][]int, 0, 20)
	for x := val; x >= 0; x-- {
		tvals := distribute(val-x, buckets-1)
		for _, vs := range tvals {
			vals = append(vals, append([]int{x}, vs...))
		}
	}
	dDP[k] = vals
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
	lines, _ := file.GetLines(path)
	crs := make([]*conditionRecord, len(lines))
	for i, line := range lines {
		cr := conditionRecord{}
		cr.vals = make([]uint64, 0, 5)
		tokens := strings.Split(line, " ")
		condition := tokens[0]
		checks := tokens[1]
		cr.condition = condition
		cr.checks = array.ToNumbers[int](checks, ",")
		crs[i] = &cr
	}
	return crs
}
