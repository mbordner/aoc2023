package main

import (
	"fmt"

	"github.com/mbordner/aoc2023/common/array"
	"github.com/mbordner/aoc2023/common/files"
)

func next(vals []int64) int64 {
	vs := make([][]int64, 0, 10)
	vs = append(vs, array.ReverseNumbers(vals))

	j := 0
	l := len(vals) - 1
	for l > 0 {
		j++
		vs = append(vs, make([]int64, len(vals), len(vals)))
		for i := 0; i < l; i++ {
			vs[j][i] = vs[j-1][i] - vs[j-1][i+1]
		}
		vs[j] = vs[j][0:l]
		if array.AllSameNumbers(vs[j]) && vs[j][0] == int64(0) {
			break
		}
		l--
	}

	sum := int64(0)
	for j := range vs {
		sum += vs[j][0]
	}
	return sum
}

func main() {
	data := getData("../data.txt")
	sum := int64(0)
	for i := range data {
		v := next(data[i])
		sum += v
	}
	fmt.Println(sum)
}

func getData(path string) [][]int64 {
	lines, _ := files.GetLines(path)
	data := make([][]int64, len(lines), len(lines))
	for i := range lines {
		data[i] = array.ToNumbers[int64](lines[i], " ")
	}
	return data
}
