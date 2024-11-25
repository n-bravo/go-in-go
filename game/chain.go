package game

import "fmt"

type Chain struct {
	id        int
	isBlack   bool
	board     *Board
	points    []*Point
	liberties int
}

func NewChain(id int, p *Point) (*Chain, error) {
	l := 0
	for _, n := range p.neighbords {
		if n.State == p.State {
			return nil, fmt.Errorf("cannot create a new chain in (%v, %v)", p.X, p.Y)
		}
		if n.State == FREE {
			l += 1
		}
	}
	c := Chain{id: id, isBlack: p.State == BLACK, board: p.board, points: make([]*Point, 1), liberties: l}
	c.points[0] = p
	return &c, nil
}

func (c *Chain) add(p *Point) {
	c.points = append(c.points, p)
	c.updateLiberties()
}

func (c *Chain) updateLiberties() {
	l := 0
	seenFree := make(map[*Point]bool)
	for _, p := range c.points {
		for _, n := range p.neighbords {
			if _, seen := seenFree[n]; !seen && n.State == FREE {
				l += 1
				seenFree[n] = true
			}
		}
	}
	c.liberties = l
}

func (c1 *Chain) merge(c2 *Chain) {
	for _, p2 := range c2.points {
		p2.chainId = c1.id
	}
	c1.points = append(c1.points, c2.points...)
	c1.updateLiberties()
	delete(c1.board.Chains, c2.id)
}
