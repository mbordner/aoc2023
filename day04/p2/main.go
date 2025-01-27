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
	id      int64
	winners map[int64]int64
	numbers map[int64]int64
}

func newCard(id int64) card {
	c := card{id: id}
	c.winners = make(map[int64]int64)
	c.numbers = make(map[int64]int64)
	return c
}

func (c card) score() int64 {
	count := c.matches()
	if count == 0 {
		return 0
	}
	return int64(math.Pow(float64(2), float64(count-1)))
}

func (c card) matches() int64 {
	count := int64(0)
	for n := range c.numbers {
		if _, e := c.winners[n]; e {
			count++
		}
	}
	return count
}

var (
	reID     = regexp.MustCompile(`^Card\s+(\d+):\s*`)
	reDigits = regexp.MustCompile(`\d+`)
)

func main() {
	cards := getCards("../data.txt")
	matches := make([]int64, len(cards), len(cards))
	for i, c := range cards {
		matches[i] = c.matches()
	}

	collection := make([]int64, 0, len(cards))

	for _, c := range cards {
		collection = append(collection, c.id)
	}

	for i := 0; i < len(collection); i++ {
		id := collection[i]
		for j := int64(0); j < matches[id-1]; j++ {
			collection = append(collection, cards[id+j].id)
		}
	}

	fmt.Println(len(collection))
}

func getCards(path string) []card {
	lines, _ := files.GetLines(path)
	cards := make([]card, len(lines), len(lines))
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
