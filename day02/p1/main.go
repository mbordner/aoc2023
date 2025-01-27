package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
)

const (
	RED   = "red"
	BLUE  = "blue"
	GREEN = "green"
)

var (
	reID     = regexp.MustCompile(`^Game (\d+): `)
	reReveal = regexp.MustCompile(`(\d+) (red|green|blue)`)
)

type gameObj struct {
	id    int64
	maxes map[string]int64
}

func NewGameObj(id int64) gameObj {
	g := gameObj{id: id}
	g.maxes = make(map[string]int64)
	return g
}

func main() {
	games := getGameObjs("../data.txt")

	possible := make(map[string]int64)
	possible[RED] = 12
	possible[GREEN] = 13
	possible[BLUE] = 14

	sum := int64(0)
	count := 0
	for _, g := range games {
		isPossible := true
		for c, v := range possible {
			if g.maxes[c] > v {
				isPossible = false
				break
			}
		}
		if isPossible {
			count++
			sum += g.id
		}
	}
	fmt.Println(count, sum)
}

func getGameObjs(path string) []gameObj {
	lines, _ := files.GetLines(path)
	games := make([]gameObj, len(lines), len(lines))

	for i, line := range lines {

		idMatch := reID.FindStringSubmatch(line)

		id, _ := strconv.ParseInt(idMatch[1], 10, 64)
		gObj := NewGameObj(id)

		plays := strings.Split(string(line[len(idMatch[0]):]), "; ")
		if len(plays) > 0 {
			for _, reveals := range plays {
				reveal := strings.Split(reveals, ", ")
				for _, r := range reveal {
					matches := reReveal.FindStringSubmatch(r)
					if len(matches) != 3 {
						log.Fatalln("invalid reveal")
					}
					val, _ := strconv.ParseInt(matches[1], 10, 64)
					color := matches[2]
					if val > gObj.maxes[color] {
						gObj.maxes[color] = val
					}
				}
			}
		}

		games[i] = gObj
	}

	return games
}
