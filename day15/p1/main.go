package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/file"
	"log"
	"strings"
)

func hash(s string) int {
	bs := []byte(s)
	val := 0
	for _, b := range bs {
		val += int(b)
		val *= 17
		val %= 256
	}
	return val
}

func main() {
	if hash(`HASH`) != 52 {
		log.Fatalln("doesn't work")
	}
	data := getData("../data.txt")
	sum := 0
	for _, s := range data {
		sum += hash(s)
	}
	fmt.Println(sum)
}

func getData(path string) []string {
	lines, _ := file.GetLines(path)
	return strings.Split(lines[0], ",")
}
