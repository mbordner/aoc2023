package main

import (
	"cmp"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
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
	memo      map[int64]int64
}

func newMapCategory(name string) *mapCategory {
	mc := new(mapCategory)
	mc.name = name
	mc.mapRanges = make([]mapRange, 0, 20)
	mc.memo = make(map[int64]int64)
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
	if r, e := c.memo[val]; e {
		return r
	}
	rval := val
	for _, r := range c.mapRanges {
		if val >= r.src && val < r.src+r.len {
			rval = r.dst + val - r.src
			break
		} else if val < r.src {
			break
		}
	}
	c.memo[val] = rval
	return rval
}

func mapVal(val int64, mappers []*mapCategory) int64 {
	mappedVal := mappers[0].mapVal(val)
	if len(mappers) > 1 {
		mappedVal = mapVal(mappedVal, mappers[1:])
	}
	return mappedVal
}

func splitRange(mappers []*mapCategory, seed, length int64) int64 {
	if length == 1 {
		return length
	}

	d := func(s int64) int64 {
		return mapVal(s, mappers) - s
	}

	mid := func(li, hi int64) int64 {
		return (hi-li)/2 + li
	}

	s0 := seed
	d0 := d(s0)

	s1 := seed + length - 1
	d1 := d(s1)

	if d0 != d1 {

		li := s0
		hi := s1

		for {
			var dm int64
			if hi == li+1 {
				dm = d(hi)
				if dm == d0 {
					return hi - s0 + 1
				} else {
					return li - s0 + 1
				}
			}

			mi := mid(li, hi)
			dm = d(mi)

			// if dm == d0, need to go right
			if dm == d0 {
				li = mi
			} else {
				// if dm != d0, need to go left
				hi = mi
			}
		}

	}

	return length

}

func main() {
	seeds, mappers := getData("../data.txt")

	var min int64

	for g := 0; g < len(seeds); g = g + 2 {
		seed := seeds[g]
		fullLength := seeds[g+1]

		loc := mapVal(seed, mappers)
		if g == 0 {
			min = loc
		} else {
			if loc < min {
				min = loc
			}
		}

		for l := splitRange(mappers, seed, fullLength); l != fullLength; l = splitRange(mappers, seed, fullLength) {
			seed = seed + l
			fullLength -= l

			loc = mapVal(seed, mappers)
			if loc < min {
				min = loc
			}
		}
	}

	fmt.Println(min)
}

func getData(path string) ([]int64, []*mapCategory) {
	lines, _ := files.GetLines(path)
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
