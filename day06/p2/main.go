package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/array/ints"
	"github.com/mbordner/aoc2023/common/file"
	"strings"
)

func wins(time, distance int64) (distances []int64) {
	distances = make([]int64, 0, time)
	for i := int64(0); i <= time; i++ {
		d := i * (time - i)
		if d > distance {
			distances = append(distances, d)
		}
	}
	return
}

func main() {
	times, distances := getData("../data.txt")
	margin := int64(1)
	for i := range times {
		ds := wins(times[i], distances[i])
		margin *= int64(len(ds))
	}
	fmt.Println(margin)
}

func getData(path string) (times, distances []int64) {
	lines, _ := file.GetLines(path)
	times = ints.NumVals(strings.ReplaceAll(lines[0][strings.Index(lines[0], ":")+1:], " ", ""))
	distances = ints.NumVals(strings.ReplaceAll(lines[1][strings.Index(lines[1], ":")+1:], " ", ""))
	return
}
