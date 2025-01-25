package game

import (
	"fmt"
	"strconv"
	"strings"
)

type board struct {
	size       int //board is size x size
	field      [][]Point
	prevField  string
	chains     map[int]*chain
	prevChains string
}

func newBoard(s int) (*board, error) {
	board := board{size: s}
	board.chains = make(map[int]*chain)
	board.field = make([][]Point, s)
	board.prevField = strings.Repeat("*", s*s)
	board.prevChains = ""
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

func (b *board) checkpoint() {
	b.prevField = b.String()
	var sb strings.Builder
	i := 0
	for _, c := range b.chains {
		sb.WriteString(c.encode())
		if i < (len(b.chains) - 1) {
			sb.WriteString("%")
		}
		i += 1
	}
	b.prevChains = sb.String()
}

func (b *board) rollBack() error {
	//restore field
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			prevState := b.prevField[b.size*i+j]
			switch prevState {
			case 'B':
				b.field[i][j].State = BLACK
			case 'W':
				b.field[i][j].State = WHITE
			default:
				b.field[i][j].State = FREE
			}
		}
	}
	//restore chains
	b.chains = make(map[int]*chain)
	chainStrSlc := strings.Split(b.prevChains, "%")
	for _, chainStr := range chainStrSlc {
		data := strings.Split(chainStr, "-")
		id, err := strconv.Atoi(data[1])
		if err != nil {
			return err
		}
		liberties, err := strconv.Atoi(data[2])
		if err != nil {
			return err
		}
		c := chain{id: id, isBlack: data[0] == "1", liberties: liberties, points: make([]*Point, 0)}
		pxy := strings.Split(data[3], "|")
		for _, p := range pxy {
			xy := strings.Split(p, ",")
			x, err := strconv.Atoi(xy[0])
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(xy[1])
			if err != nil {
				return err
			}
			c.points = append(c.points, &b.field[x][y])
			b.field[x][y].chainId = c.id
		}
		b.chains[c.id] = &c
	}
	return nil
}

func (b *board) play(x, y int, black bool) (int, error) { //black is true if the play is from black stones player
	if x < 0 || x > (b.size-1) || y < 0 || y > (b.size-1) {
		return 0, fmt.Errorf("invalid position (%v, %v)", x, y)
	}
	err := b.field[x][y].play(black)
	if err != nil {
		return 0, err
	}
	captured, err := b.check(black, x, y)
	if err != nil {
		b.field[x][y].free()
		return 0, err
	}
	b.checkpoint()
	return captured, nil
}

func (b *board) deleteChain(cid int) (captured int) {
	captured = len(b.chains[cid].points)
	b.chains[cid].free()
	delete(b.chains, cid)
	return
}

func (b *board) check(black bool, x int, y int) (int, error) {
	captured := 0
	for cid := range b.chains {
		if b.chains[cid].isBlack != black && b.chains[cid].liberties == 0 {
			n := b.deleteChain(cid)
			captured += n
		}
	}
	if b.chains[b.field[x][y].chainId].liberties == 0 {
		//self-captured, rollback and throw error
		err := b.rollBack()
		if err != nil {
			panic(err)
		}
		return 0, fmt.Errorf("error self-capture forbidden")
	}
	return captured, nil
}
