package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameCreation(t *testing.T) {
	assert := assert.New(t)
	var err error
	_, err = NewGame(0)
	assert.Error(err)
	_, err = NewGame(-2)
	assert.Error(err)
	_, err = NewGame(5)
	assert.NoError(err)
}

func TestPointsWellCreated(t *testing.T) {
	assert := assert.New(t)
	g, _ := NewGame(5)
	assert.Equal(3, g.board.field[3][4].X)
	assert.Equal(4, g.board.field[3][4].Y)
	assert.Equal(0, g.board.field[0][0].X)
	assert.Equal(0, g.board.field[2][0].Y)
	assert.Len(g.board.field[0][0].neighbords, 2)
	assert.Len(g.board.field[0][4].neighbords, 2)
	assert.Len(g.board.field[4][0].neighbords, 2)
	assert.Len(g.board.field[4][4].neighbords, 2)
	assert.Len(g.board.field[3][3].neighbords, 4)
	assert.Contains(g.board.field[3][3].neighbords, &g.board.field[3][4])
	assert.Contains(g.board.field[3][3].neighbords, &g.board.field[3][2])
	assert.NotEmpty(g.board.field[0][0].neighbords)
	assert.Contains(g.board.field[0][0].neighbords, &g.board.field[1][0])
	assert.NotContains(g.board.field[0][0].neighbords, &g.board.field[3][3])
}

func TestBoardPlay(t *testing.T) {
	assert := assert.New(t)
	g, _ := NewGame(5)
	var err error
	err = g.Play(-1, 0, true)
	assert.Error(err)
	err = g.Play(0, 5, true)
	assert.Error(err)
	err = g.Play(0, 6, true)
	assert.Error(err)
	err = g.Play(1, 1, true)
	assert.NoError(err)
	err = g.Play(1, 1, false)
	assert.Error(err)
	err = g.Play(1, 2, false)
	assert.NoError(err)
	err = g.Play(1, 2, true)
	assert.Error(err)
	err = g.Play(1, 3, true)
	assert.NoError(err)
	err = g.Play(1, 4, true)
	assert.Error(err)
	err = g.Play(1, 4, false)
	assert.NoError(err)
	err = g.Play(1, 2, false)
	assert.Error(err)
	assert.Equal("******BWBW***************", g.String())
}

func TestChains(t *testing.T) {
	assert := assert.New(t)
	g, _ := NewGame(5)
	g.Play(1, 1, true)
	assert.Equal(1, len(g.board.chains))
	assert.Contains(g.board.chains, 1)
	assert.Equal(1, len(g.board.chains[1].points))
	assert.Equal(4, g.board.chains[1].liberties)
	assert.Same(&g.board.field[1][1], g.board.chains[1].points[0])
	g.Play(1, 2, false)
	assert.Equal(2, len(g.board.chains))
	assert.Contains(g.board.chains, 2)
	assert.Equal(1, len(g.board.chains[2].points))
	assert.Equal(3, g.board.chains[1].liberties)
	assert.Equal(3, g.board.chains[2].liberties)
	assert.Same(&g.board.field[1][2], g.board.chains[2].points[0])
	g.Play(2, 1, true)
	assert.Equal(2, len(g.board.chains))
	assert.Contains(g.board.chains, 1)
	assert.Equal(2, len(g.board.chains[1].points))
	assert.Equal(5, g.board.chains[1].liberties)
	assert.Same(&g.board.field[2][1], g.board.chains[1].points[1])
	g.Play(2, 2, false)
	assert.Equal(2, len(g.board.chains))
	assert.Contains(g.board.chains, 2)
	assert.Equal(2, len(g.board.chains[2].points))
	assert.Equal(4, g.board.chains[2].liberties)
	assert.Same(&g.board.field[2][2], g.board.chains[2].points[1])
	g.Play(3, 2, true)
	assert.Equal(3, len(g.board.chains))
	assert.Contains(g.board.chains, 3)
	assert.Equal(1, len(g.board.chains[3].points))
	assert.Equal(4, g.board.chains[1].liberties)
	assert.Equal(3, g.board.chains[2].liberties)
	assert.Equal(3, g.board.chains[3].liberties)
	assert.Same(&g.board.field[3][2], g.board.chains[3].points[0])
	g.Play(2, 3, false)
	assert.Equal(3, len(g.board.chains))
	assert.Contains(g.board.chains, 2)
	assert.Equal(3, len(g.board.chains[2].points))
	assert.Equal(4, g.board.chains[1].liberties)
	assert.Equal(4, g.board.chains[2].liberties)
	assert.Equal(3, g.board.chains[3].liberties)
	assert.Same(&g.board.field[2][3], g.board.chains[2].points[2])
	g.Play(3, 1, true)
	assert.Equal(2, len(g.board.chains))
	assert.Contains(g.board.chains, 1)
	assert.NotContains(g.board.chains, 3)
	assert.Equal(4, len(g.board.chains[1].points))
	assert.Equal(7, g.board.chains[1].liberties)
	assert.Equal(4, g.board.chains[2].liberties)
	assert.Same(&g.board.field[3][1], g.board.chains[1].points[2])
	g.Play(4, 4, false)
	assert.Equal(3, len(g.board.chains))
	assert.Contains(g.board.chains, 3)
	assert.Equal(1, len(g.board.chains[3].points))
	assert.Equal(7, g.board.chains[1].liberties)
	assert.Equal(4, g.board.chains[2].liberties)
	assert.Equal(2, g.board.chains[3].liberties)
	assert.Same(&g.board.field[4][4], g.board.chains[3].points[0])
}

func TestWhiteCapturesMultipleBlackWithMiddleMove(t *testing.T) {
	//      B B W * B
	//      * W B W *
	//      W B * B W
	//      * W B W *
	//      B * W * B
	//      	|
	//      	| white turn
	//      	V
	//      B B W * B
	//      * W * W *
	//      W * W * W
	//      * W * W *
	//      B * W * B
	assert := assert.New(t)
	g, _ := NewGame(5)
	g.Play(1, 2, true)
	g.Play(0, 2, false)
	g.Play(2, 1, true)
	g.Play(2, 0, false)
	g.Play(3, 2, true)
	g.Play(4, 2, false)
	g.Play(2, 3, true)
	g.Play(2, 4, false)
	g.Play(0, 0, true)
	g.Play(1, 1, false)
	g.Play(4, 0, true)
	g.Play(3, 1, false)
	g.Play(0, 4, true)
	g.Play(1, 3, false)
	g.Play(4, 4, true)
	g.Play(3, 3, false)
	g.Play(0, 1, true)
	err := g.Play(2, 2, false)
	assert.NoError(err)
	assert.Equal("BBW*B*W*W*W*W*W*W*W*B*W*B", g.String())
}

func TestNoSelfCapture(t *testing.T) {
	assert := assert.New(t)
	g, _ := NewGame(5)
	g.Play(0, 1, true)
	g.Play(0, 3, false)
	g.Play(1, 1, true)
	g.Play(1, 3, false)
	g.Play(1, 0, true)
	err := g.Play(0, 0, false)
	assert.Error(err)
	g.Play(1, 4, false)
	err = g.Play(0, 4, true)
	assert.Error(err)
}
