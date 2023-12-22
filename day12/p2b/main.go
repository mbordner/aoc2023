package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/file"
	"strings"
)

type conditionRecord struct {
	condition string
	checks    []int
}

func main() {
	crs := getData("../data.txt")

	sum := uint64(0)

	for _, cr := range crs {
		val := arrangements(cr.condition, cr.checks)
		sum += val
	}

	fmt.Println(sum)
}

func arrangements(condition string, checks []int) uint64 {

	type dps struct {
		p  int
		ci int
		cl int
	}

	dp := make(map[dps]uint64)
	var f func(p, ci, cl int) uint64

	handleDot := func(p, ci, cl int) uint64 {
		if cl > 0 {
			if cl == checks[ci] {
				if ci+1 <= len(checks) {
					return f(p+1, ci+1, 0)
				}
			}
		} else {
			return f(p+1, ci, 0)
		}
		return 0
	}

	handleHash := func(p, ci, cl int) uint64 {
		if ci < len(checks) && cl < checks[ci] {
			return f(p+1, ci, cl+1)
		}
		return 0
	}

	f = func(p int, ci int, cl int) uint64 {

		dpk := dps{p: p, ci: ci, cl: cl}
		if v, e := dp[dpk]; e {
			return v
		}

		v := uint64(0)

		if p == len(condition) {

			if cl > 0 {
				if cl == checks[ci] {
					cl = 0
					ci++
				}
			}

			if ci == len(checks) {
				v++
			}

		} else {
			b := condition[p]

			switch b {
			case '?':
				v += handleDot(p, ci, cl)
				v += handleHash(p, ci, cl)
			case '.':
				v += handleDot(p, ci, cl)
			case '#':
				v += handleHash(p, ci, cl)
			}
		}

		dp[dpk] = v
		return v
	}

	return f(0, 0, 0)
}

func getData(path string) []*conditionRecord {
	lines, _ := file.GetLines(path)
	crs := make([]*conditionRecord, len(lines))
	for i, line := range lines {
		cr := conditionRecord{}
		tokens := strings.Split(line, " ")
		condition := tokens[0]
		checks := tokens[1]

		for j := 0; j < 4; j++ {
			condition += "?" + tokens[0]
			checks += "," + tokens[1]
		}

		cr.condition = condition
		cr.checks = array.ToNumbers[int](checks, ",")
		crs[i] = &cr
	}
	return crs
}
