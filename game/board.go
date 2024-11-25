package game

import (
	"fmt"
	"strings"
)

type Board struct {
	Size            int //board is size x size
	Field           [][]Point
	Chains          map[int]*Chain
	BlackPlayedLast bool
}

func NewBoard(s int) (*Board, error) {
	if s <= 0 {
		return nil, fmt.Errorf("invalid board size (%v x %v)", s, s)
	}
	board := Board{Size: s}
	board.Chains = make(map[int]*Chain)
	board.Field = make([][]Point, s)
	for i := 0; i < s; i++ {
		board.Field[i] = make([]Point, s)
	}
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			board.Field[i][j].Init(&board, i, j)
		}
	}
	return &board, nil
}

func (b Board) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Size: %v, Chains: %v\n", b.Size, len(b.Chains)))
	for x := range b.Field {
		for y := range b.Field[x] {
			sb.WriteString(b.Field[x][y].String())
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *Board) Play(x, y int, black bool) error { //black is true if the play is from black stones player
	if x < 0 || x > (b.Size-1) || y < 0 || y > (b.Size-1) {
		return fmt.Errorf("invalid position (%v, %v)", x, y)
	}
	if b.BlackPlayedLast == black {
		switch b.BlackPlayedLast {
		case true:
			return fmt.Errorf("invalid turn. now white must play")
		default:
			return fmt.Errorf("invalid turn. now black must play")
		}
	}
	err := b.Field[x][y].play(black)
	if err != nil {
		return err
	}
	b.BlackPlayedLast = !b.BlackPlayedLast
	b.check(x, y)
	return nil
}

func (b *Board) check(x, y int) { //TODO: deberia chequear los chains que no tienen libertad y capturarlos

}
