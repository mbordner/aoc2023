package main

import (
	"fmt"
	"github.com/mbordner/aoc2023/common/file"
	"log"
	"regexp"
	"strconv"
	"strings"
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

	sum := int64(0)
	for _, g := range games {
		v := int64(1)
		for _, m := range g.maxes {
			v *= m
		}
		sum += v
	}
	fmt.Println(sum)
}

func getGameObjs(path string) []gameObj {
	lines, _ := file.GetLines(path)
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
