package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
)

type card struct {
	id           int64
	winners      map[int64]int64
	numbers      map[int64]int64
	matchesValue int64
}

func newCard(id int64) *card {
	c := card{id: id}
	c.winners = make(map[int64]int64)
	c.numbers = make(map[int64]int64)
	c.matchesValue = -1
	return &c
}

func (c *card) score() int64 {
	count := c.matches()
	if count == 0 {
		return 0
	}
	return int64(math.Pow(float64(2), float64(count-1)))
}

func (c *card) matches() int64 {
	if c.matchesValue >= 0 {
		return c.matchesValue
	}
	c.matchesValue = int64(0)
	for n := range c.numbers {
		if _, e := c.winners[n]; e {
			c.matchesValue++
		}
	}
	return c.matchesValue
}

var (
	reID     = regexp.MustCompile(`^Card\s+(\d+):\s*`)
	reDigits = regexp.MustCompile(`\d+`)
)

func count(cards []*card, index int, length int) int64 {
	collectionCount := int64(0)
	for i := index; i < index+length; i++ {
		c := cards[i]
		collectionCount += 1
		matches := c.matches()
		if matches > 0 {
			collectionCount += count(cards, int(c.id), int(matches))
		}
	}
	return collectionCount
}

func main() {
	cards := getCards("../data.txt")
	fmt.Println(count(cards, 0, len(cards)))
}

func getCards(path string) []*card {
	lines, _ := files.GetLines(path)
	cards := make([]*card, len(lines), len(lines))
	for i, line := range lines {
		idMatch := reID.FindStringSubmatch(line)

		id, _ := strconv.ParseInt(idMatch[1], 10, 64)
		card := newCard(id)

		sets := strings.Split(string(line[len(idMatch[0]):]), " | ")
		if len(sets) != 2 {
			log.Fatalln("invalid")
		}

		cards[i] = card
		matches := reDigits.FindAllString(sets[0], -1)
		for _, m := range matches {
			v, _ := strconv.ParseInt(m, 10, 64)
			if _, e := card.winners[v]; e {
				log.Fatalln("dupe")
			}
			card.winners[v] = v
		}
		matches = reDigits.FindAllString(sets[1], -1)
		for _, m := range matches {
			v, _ := strconv.ParseInt(m, 10, 64)
			if _, e := card.numbers[v]; e {
				log.Fatalln("dupe")
			}
			card.numbers[v] = v
		}
	}
	return cards
}
