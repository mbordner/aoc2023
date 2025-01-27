package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
)

type lensObj struct {
	focal int
}

type boxObj struct {
	slots    []int
	labelMap map[string]int
}

func newBox() *boxObj {
	b := new(boxObj)

	b.slots = make([]int, 0, 9)
	b.labelMap = make(map[string]int)

	return b
}

func (b *boxObj) add(label string, focal int) {
	if label == "" {
		log.Fatal("empty label?")
	}
	if i, e := b.labelMap[label]; e {
		b.slots[i] = focal
	} else {
		b.labelMap[label] = len(b.slots)
		b.slots = append(b.slots, focal)
	}
}

func (b *boxObj) rm(label string) {
	if i, e := b.labelMap[label]; e {
		b.slots = append(b.slots[0:i], b.slots[i+1:]...)
		delete(b.labelMap, label)
		for k, v := range b.labelMap {
			if v >= i {
				b.labelMap[k] = v - 1
			}
		}
	}
}

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
	data := getData("../data.txt")

	boxes := make([]*boxObj, 256)
	for i := range boxes {
		boxes[i] = newBox()
	}

	for _, s := range data {
		if strings.Contains(s, "=") {
			tokens := strings.Split(s, "=")
			label := tokens[0]
			h := hash(label)
			focal, err := strconv.ParseInt(tokens[1], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			boxes[h].add(label, int(focal))
		} else if strings.Contains(s, "-") {
			tokens := strings.Split(s, "-")
			label := tokens[0]
			h := hash(label)
			boxes[h].rm(label)
		} else {
			log.Fatal("no operation?")
		}
	}

	fp := 0
	for b := range boxes {
		bfp := 0
		for l := range boxes[b].slots {
			bfp += (b + 1) * (l + 1) * boxes[b].slots[l]
		}
		fp += bfp
	}

	fmt.Println(fp)
}

func getData(path string) []string {
	lines, _ := files.GetLines(path)
	return strings.Split(lines[0], ",")
}
