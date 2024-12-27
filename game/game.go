package game

import (
	"fmt"
	"strings"
)

type GoGame struct {
	BlackPlayedLast bool
	BlackCaptures   int
	WhiteCaptures   int
	board           *board
}

func NewGame(n int) (*GoGame, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid board size (%v x %v)", n, n)
	}
	g := GoGame{}
	b, err := newBoard(n)
	if err != nil {
		return nil, err
	}
	g.board = b
	return &g, nil
}

func (g GoGame) String() string {
	//var nextTurn string
	//if g.BlackPlayedLast {
	//	nextTurn = "W"
	//} else {
	//	nextTurn = "B"
	//}
	var sb strings.Builder
	//sb.WriteString(fmt.Sprintf("Next turn: %v | B captures: %v | W captures: %v | ", nextTurn, g.BlackCaptures, g.WhiteCaptures))
	sb.WriteString(g.board.String())
	return sb.String()
}

func (g *GoGame) Play(x, y int, black bool) error {
	if g.BlackPlayedLast == black {
		switch g.BlackPlayedLast {
		case true:
			return fmt.Errorf("invalid turn. now white must play")
		default:
			return fmt.Errorf("invalid turn. now black must play")
		}
	}
	cap, err := g.board.play(x, y, black)
	if err != nil {
		return err
	}
	if black {
		g.WhiteCaptures += cap
	} else {
		g.BlackCaptures += cap
	}
	g.BlackPlayedLast = !g.BlackPlayedLast
	return nil
}

func (g *GoGame) Close() error {
	g.board = nil
	return nil
}
