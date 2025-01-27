package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
)

type typeHand int
type cardCounts map[string]int
type handBid struct {
	hand string
	bid  int64
}

func (hb *handBid) String() string {
	return fmt.Sprintf("[%s %d]", hb.hand, hb.bid)
}

const (
	handHighCard    typeHand = iota // 5 unique 1,1,1,1,1
	handOnePair                     // 4 unique 2,1,1,1
	handTwoPair                     // 3 unique 2,2,1 (4)
	handThreeOfKind                 // 3 unique 3,1,1 (3)
	handFullHouse                   // 2 unique 3,2 (6)
	handFourOfKind                  // 2 unique 4,1 (4)
	handFiveOfKind                  // 1 unique 1
)

var (
	cardRanking = "23456789TJQKA"
)

func countCards(hand string) cardCounts {
	cc := make(cardCounts)
	for _, c := range hand {
		s := string(c)
		if v, e := cc[s]; e {
			cc[s] = v + 1
		} else {
			cc[s] = 1
		}
	}
	return cc
}

func cardCountsProduct(cc cardCounts) int {
	p := 1
	for _, v := range cc {
		p *= v
	}
	return p
}

func handType(hand string) typeHand {
	cc := countCards(hand)
	switch len(cc) {
	case 1:
		return handFiveOfKind
	case 2:
		if cardCountsProduct(cc) == 4 {
			return handFourOfKind
		} else {
			return handFullHouse
		}
	case 3:
		if cardCountsProduct(cc) == 3 {
			return handThreeOfKind
		} else {
			return handTwoPair
		}
	case 4:
		return handOnePair
	}
	return handHighCard
}

func compareHand(a, b string) int {
	v := handType(a)
	w := handType(b)
	if v > w {
		return 1
	} else if v < w {
		return -1
	} else {
		for i := 0; i < len(a); i++ {
			x := string(a[i])
			y := string(b[i])
			c := compareCard(x, y)
			if c > 0 {
				return 1
			} else if c < 0 {
				return -1
			}
		}
	}
	return 0
}

func compareCard(a, b string) int {
	v := strings.Index(cardRanking, a)
	w := strings.Index(cardRanking, b)
	if v > w {
		return 1
	} else if v < w {
		return -1
	}
	return 0
}

func main() {
	handBids := getData("../data.txt")
	slices.SortFunc(handBids, func(a, b *handBid) int {
		return compareHand(a.hand, b.hand)
	})
	winnings := int64(0)
	for i := 0; i < len(handBids); i++ {
		r := int64(i) + 1
		winnings += r * handBids[i].bid
	}
	fmt.Println(winnings)
}

func getData(path string) []*handBid {
	lines, _ := files.GetLines(path)
	handBids := make([]*handBid, len(lines), len(lines))
	for i, line := range lines {
		tokens := strings.Split(line, " ")
		bid, _ := strconv.ParseInt(tokens[1], 10, 64)
		hb := handBid{
			hand: tokens[0],
			bid:  bid,
		}
		handBids[i] = &hb
	}
	return handBids
}
