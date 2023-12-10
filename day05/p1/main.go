package main

import (
	"cmp"
	"fmt"
	"github.com/mbordner/aoc2023/common/file"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	reCat = regexp.MustCompile(`^\w+-to-(\w+) map:`)
)

type mapRange struct {
	src int64
	dst int64
	len int64
}

type mapCategory struct {
	name      string
	mapRanges []mapRange
}

func newMapCategory(name string) *mapCategory {
	mc := new(mapCategory)
	mc.name = name
	mc.mapRanges = make([]mapRange, 0, 20)
	return mc
}

func (c *mapCategory) addRange(vals []int64) {
	if len(vals) != 3 {
		log.Fatalln("invalid range values")
	}
	r := mapRange{
		src: vals[1],
		dst: vals[0],
		len: vals[2],
	}
	c.mapRanges = append(c.mapRanges, r)
}

func (c *mapCategory) mapVal(val int64) int64 {
	for _, r := range c.mapRanges {
		if val >= r.src && val < r.src+r.len {
			return r.dst + val - r.src
		} else if val < r.src {
			break
		}
	}
	return val
}

func mapVal(val int64, mappers []*mapCategory) int64 {
	mappedVal := mappers[0].mapVal(val)
	if len(mappers) > 1 {
		mappedVal = mapVal(mappedVal, mappers[1:])
	}
	return mappedVal
}

func main() {
	seeds, mappers := getData("../data.txt")
	locations := make([]int64, len(seeds), len(seeds))
	var max, min int64
	for i, s := range seeds {
		locations[i] = mapVal(s, mappers)
		if i == 0 {
			min = locations[i]
			max = locations[i]
		} else {
			if locations[i] < min {
				min = locations[i]
			} else if locations[i] > max {
				max = locations[i]
			}
		}
	}
	fmt.Println(seeds)
	fmt.Println(locations)
	fmt.Println(min, max)
}

func getData(path string) ([]int64, []*mapCategory) {
	lines, _ := file.GetLines(path)
	seeds := numValues(lines[0][7:])
	cats := make([]*mapCategory, 0, 7)

	i := 2
	for i < len(lines) && len(lines[i]) > 0 {
		match := reCat.FindStringSubmatch(lines[i])
		if match != nil {
			cat := newMapCategory(match[1])
			j := i + 1
			for j < len(lines) && len(lines[j]) > 0 {
				cat.addRange(numValues(lines[j]))
				j++
			}
			i = j + 1
			slices.SortFunc(cat.mapRanges, func(a, b mapRange) int {
				return cmp.Compare(a.src, b.src)
			})
			cats = append(cats, cat)
		} else {
			log.Fatalln("invalid map start")
		}
	}

	return seeds, cats
}

func numValues(svals string) []int64 {
	tokens := strings.Split(svals, " ")
	vals := make([]int64, len(tokens), len(tokens))
	for i, t := range tokens {
		v, _ := strconv.ParseInt(t, 10, 64)
		vals[i] = v
	}
	return vals
}
