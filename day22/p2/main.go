package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
	"github.com/mbordner/aoc2023/common/geom"
)

func main() {
	cuboids := getData("../data.txt")

	sort.Slice(cuboids, func(i, j int) bool {
		return cuboids[i].Less(cuboids[j])
	})

	above := make([][]int, len(cuboids))
	above[0] = []int{}
	for i := 1; i < len(cuboids); i++ {
		above[i] = []int{}
		t := cuboids[i].Clone()
		t.Min.Z = 0
		for j := i - 1; j >= 0; j-- {
			if t.Overlaps(cuboids[j]) {
				above[i] = append(above[i], j)
			}
		}
	}

	down := geom.Vector{Z: -1}

	for i := 0; i < len(cuboids); i++ {
		floor := int64(0)
		if len(above[i]) > 0 {
			for _, a := range above[i] {
				if cuboids[a].Max.Z > floor {
					floor = cuboids[a].Max.Z
				}
			}
		}
		for cuboids[i].Min.Z > floor+1 {
			cuboids[i] = cuboids[i].Transform(down)
		}
	}

	bricks := make([]*Brick, len(cuboids))
	for i := 0; i < len(cuboids); i++ {
		bricks[i] = NewBrick(cuboids[i])
	}

	for i, b := range bricks {
		t := b.Cube.Transform(down)
		for _, a := range above[i] {
			if t.Overlaps(bricks[a].Cube) {
				b.AddSupportedBy(bricks[a])
				bricks[a].AddSupports(b)
			}
		}
	}
	unsafe := make([]*Brick, 0, len(bricks))

	for _, b := range bricks {
		canRemoveSafely := false
		if len(b.Supports) == 0 {
			canRemoveSafely = true
		} else {
			allSupported := true
			for _, s := range b.Supports {
				if len(s.SupportedBy) < 2 {
					allSupported = false
					break
				}
			}
			if allSupported {
				canRemoveSafely = true
			}
		}
		if !canRemoveSafely {
			unsafe = append(unsafe, b)
		}
	}

	sum := 0

	for _, b := range unsafe {
		bs := NewBricks([]*Brick{b})
		findFalling(b, bs)
		sum += len(bs) - 1
	}

	fmt.Println(sum)

}

type Bricks map[string]*Brick

func (bs Bricks) Has(b *Brick) bool {
	if _, e := bs[b.ID]; e {
		return true
	}
	return false
}

func (bs Bricks) Add(b *Brick) {
	bs[b.ID] = b
}

func NewBricks(initial []*Brick) Bricks {
	bricks := make(map[string]*Brick)
	for _, b := range initial {
		bricks[b.ID] = b
	}
	return bricks
}

func findFalling(from *Brick, fallen Bricks) {
	notSupported := make([]*Brick, 0, len(from.Supports))
	for _, o := range from.Supports {
		supported := false
		for _, s := range o.SupportedBy {
			if !fallen.Has(s) {
				supported = true
				break
			}
		}
		if !supported {
			notSupported = append(notSupported, o)
			fallen.Add(o)
		}
	}
	for _, o := range notSupported {
		findFalling(o, fallen)
	}
}

type Brick struct {
	ID          string
	Cube        geom.Cuboid
	Supports    []*Brick
	SupportedBy []*Brick
}

func (b *Brick) AddSupports(o *Brick) {
	b.Supports = append(b.Supports, o)
}

func (b *Brick) AddSupportedBy(o *Brick) {
	b.SupportedBy = append(b.SupportedBy, o)
}

func NewBrick(cube geom.Cuboid) *Brick {
	return &Brick{ID: cube.String(), Cube: cube, SupportedBy: make([]*Brick, 0), Supports: make([]*Brick, 0)}
}

func getData(filename string) geom.Cuboids {
	replacer := strings.NewReplacer("~", ",")
	lines := files.MustGetLines(filename)
	cuboids := make(geom.Cuboids, len(lines))
	for i, line := range lines {
		line = replacer.Replace(line)
		cuboids[i] = geom.NewCuboid(line)
	}
	return cuboids
}
