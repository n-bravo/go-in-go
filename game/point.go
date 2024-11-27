package game

import "fmt"

type pointStateType int

const (
	FREE pointStateType = iota
	BLACK
	WHITE
)

type Point struct {
	X, Y       int
	State      pointStateType
	neighbords []*Point
	board      *board
	chainId    int
}

func (p *Point) Init(b *board, x, y int) {
	p.X, p.Y = x, y
	p.board = b
	p.linkNeighbords()
}

func (p Point) String() string {
	switch p.State {
	case BLACK:
		return "B"
	case WHITE:
		return "W"
	default: //FREE
		return "*"
	}
}

func (p *Point) linkNeighbords() {
	if p.X == 0 && p.Y == 0 { //top left
		p.neighbords = []*Point{
			&(p.board).field[1][0], //bottom
			&(p.board).field[0][1], //right
		}
		return
	}
	if p.X == p.board.size-1 && p.Y == 0 { //bottom left
		p.neighbords = []*Point{
			&(p.board).field[p.board.size-2][0], //top
			&(p.board).field[p.board.size-1][1], //right
		}
		return
	}
	if p.X == 0 && p.Y == p.board.size-1 { //top right
		p.neighbords = []*Point{
			&(p.board).field[0][p.board.size-2], //left
			&(p.board).field[1][p.board.size-1], //bottom
		}
		return
	}
	if p.X == p.board.size-1 && p.Y == p.board.size-1 { //bottom right
		p.neighbords = []*Point{
			&(p.board).field[p.board.size-2][p.board.size-1], //top
			&(p.board).field[p.board.size-1][p.board.size-2], //left
		}
		return
	}
	if p.X == 0 { //top border
		p.neighbords = []*Point{
			&(p.board).field[0][p.Y+1], //right
			&(p.board).field[0][p.Y-1], //left
			&(p.board).field[1][p.Y],   //bottom
		}
		return
	}
	if p.X == p.board.size-1 { //bottom border
		p.neighbords = []*Point{
			&(p.board).field[p.board.size-1][p.Y+1], //right
			&(p.board).field[p.board.size-1][p.Y-1], //left
			&(p.board).field[p.board.size-2][p.Y],   //top
		}
		return
	}
	if p.Y == 0 { //left border
		p.neighbords = []*Point{
			&(p.board).field[p.X+1][0], //bottom
			&(p.board).field[p.X-1][0], //top
			&(p.board).field[p.X][1],   //right
		}
		return
	}
	if p.Y == p.board.size-1 { //right border
		p.neighbords = []*Point{
			&(p.board).field[p.X+1][p.board.size-1], //bottom
			&(p.board).field[p.X-1][p.board.size-1], //top
			&(p.board).field[p.X][p.board.size-2],   //left
		}
		return
	}
	//not in border or corners, so it has 4 neighbors
	p.neighbords = []*Point{
		&(p.board).field[p.X][p.Y+1], //right
		&(p.board).field[p.X][p.Y-1], //left
		&(p.board).field[p.X-1][p.Y], //top
		&(p.board).field[p.X+1][p.Y], //bottom
	}
}

func (p *Point) noSameNeighbor() bool {
	nc := len(p.neighbords)
	for _, n := range p.neighbords {
		if n.State != p.State {
			nc -= 1
		}
	}
	return nc == 0
}

func (p *Point) checkNeighbors() error {
	if p.noSameNeighbor() {
		chId := len(p.board.chains) + 1
		p.chainId = chId
		c, err := NewChain(chId, p)
		if err != nil {
			p.chainId = 0
			return fmt.Errorf("error creating chain: %v", err)
		}
		p.board.chains[chId] = c
		otherPlayerChains := make(map[int]bool)
		for _, n := range p.neighbords {
			if _, seen := otherPlayerChains[n.chainId]; !seen && n.State != FREE && n.State != p.State {
				p.board.chains[n.chainId].liberties -= 1
				otherPlayerChains[n.chainId] = true
			}
		}
	} else {
		playerChains := make(map[int]bool)
		minChainId := p.board.size*p.board.size + 1
		for _, n := range p.neighbords {
			if n.State != FREE {
				_, seen := playerChains[n.chainId]
				if seen {
					continue
				}
				p.board.chains[n.chainId].liberties -= 1
				if n.State == p.State {
					//pick the chain with the min id to join in and merge the other chains of the same player
					minChainId = min(minChainId, n.chainId)
					playerChains[n.chainId] = true
				} else {
					playerChains[n.chainId] = false
				}
			}
		}
		p.chainId = minChainId
		p.board.chains[p.chainId].add(p)
		for c := range playerChains {
			if c != p.chainId && playerChains[c] {
				p.board.chains[p.chainId].merge(p.board.chains[c])
			}
		}
	}
	return nil
}

func (p *Point) play(black bool) error {
	switch p.State {
	case BLACK:
		return fmt.Errorf("point already taken by black")
	case WHITE:
		return fmt.Errorf("point already taken by white")
	default: //FREE
		if black {
			p.State = BLACK
		} else {
			p.State = WHITE
		}
		err := p.checkNeighbors()
		if err != nil {
			return fmt.Errorf("error during neighbors checking: %v", err)
		}
		return nil
	}
}
