package game

import (
	"fmt"
	"strings"
)

type board struct {
	size   int //board is size x size
	field  [][]Point
	chains map[int]*chain
}

func newBoard(s int) (*board, error) {
	board := board{size: s}
	board.chains = make(map[int]*chain)
	board.field = make([][]Point, s)
	for i := 0; i < s; i++ {
		board.field[i] = make([]Point, s)
	}
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			board.field[i][j].Init(&board, i, j)
		}
	}
	return &board, nil
}

func (b board) String() string {
	var sb strings.Builder
	for x := range b.field {
		for y := range b.field[x] {
			sb.WriteString(b.field[x][y].String())
		}
	}
	return sb.String()
}

func (b *board) play(x, y int, black bool) (int, error) { //black is true if the play is from black stones player
	if x < 0 || x > (b.size-1) || y < 0 || y > (b.size-1) {
		return 0, fmt.Errorf("invalid position (%v, %v)", x, y)
	}
	err := b.field[x][y].play(black)
	if err != nil {
		return 0, err
	}
	return b.check(black), nil
}

func (b *board) deleteChain(cid int) (captured int) {
	captured = len(b.chains[cid].points)
	b.chains[cid].free()
	delete(b.chains, cid)
	return
}

func (b *board) check(black bool) int {
	captured := 0
	for cid := range b.chains {
		if b.chains[cid].isBlack != black && b.chains[cid].liberties == 0 {
			n := b.deleteChain(cid)
			captured += n
		}
	}
	return captured
}
