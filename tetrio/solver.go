package tetrio

import (
	"fmt"
	"math"
	"strings"
)

const (
	I0Shape int8 = iota
	I90Shape
	I180Shape
	I270Shape

	J0Shape
	J90Shape
	J180Shape
	J270Shape

	L0Shape
	L90Shape
	L180Shape
	L270Shape

	O0Shape
	O90Shape
	O180Shape
	O270Shape

	T0Shape
	T90Shape
	T180Shape
	T270Shape

	S0Shape
	S90Shape
	S180Shape
	S270Shape

	Z0Shape
	Z90Shape
	Z180Shape
	Z270Shape

	kShapeSize
)

type Pair struct {
	i, j int32
}

var bounds = [kShapeSize]Pair{
	{3, 0}, {0, 3}, {3, 0}, {0, 3},
	{2, 1}, {1, 2}, {2, 1}, {1, 2},
	{2, 1}, {1, 2}, {2, 1}, {1, 2},
	{1, 1}, {1, 1}, {1, 1}, {1, 1},
	{1, 2}, {2, 1}, {1, 2}, {2, 1},
	{1, 2}, {2, 1}, {1, 2}, {2, 1},
	{1, 2}, {2, 1}, {1, 2}, {2, 1},
}
var floors = [kShapeSize][]Pair{
	{{-1, 0}},
	{{-1, 0}, {-1, 1}, {-1, 2}, {-1, 3}},
	{{-1, 0}},
	{{-1, 0}, {-1, 1}, {-1, 2}, {-1, 3}},

	{{-1, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}, {-1, 2}},
	{{-1, 0}, {1, 1}},
	{{0, 0}, {0, 1}, {-1, 2}},

	{{-1, 0}, {-1, 1}},
	{{-1, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}, {-1, 2}},

	{{-1, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}},

	{{0, 0}, {-1, 1}, {0, 2}},
	{{0, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}, {-1, 2}},
	{{-1, 0}, {0, 1}},

	{{-1, 0}, {-1, 1}, {0, 2}},
	{{0, 0}, {-1, 1}},
	{{-1, 0}, {-1, 1}, {0, 2}},
	{{0, 0}, {-1, 1}},

	{{0, 0}, {-1, 1}, {-1, 2}},
	{{-1, 0}, {0, 1}},
	{{0, 0}, {-1, 1}, {-1, 2}},
	{{-1, 0}, {0, 1}},
}

var shapes = [kShapeSize][]Pair{
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},

	{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
	{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
	{{1, 0}, {1, 1}, {1, 2}, {0, 2}},

	{{0, 0}, {0, 1}, {1, 0}, {2, 0}},
	{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {1, 1}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 2}},

	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},

	{{1, 0}, {1, 1}, {1, 2}, {0, 1}},
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 1}},
	{{0, 0}, {1, 0}, {2, 0}, {1, 1}},

	{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
	{{1, 0}, {2, 0}, {0, 1}, {1, 1}},
	{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
	{{1, 0}, {2, 0}, {0, 1}, {1, 1}},

	{{1, 0}, {1, 1}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	{{1, 0}, {1, 1}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
}

type Tetris struct {
	board   [][]bool
	rows    int32
	cols    int32
	score   int32
	history []uint32
}

func NewTetris(rows, cols int32) Tetris {
	board := make([][]bool, rows+5)
	for i := range board {
		board[i] = make([]bool, cols)
	}
	return Tetris{board, rows, cols, 0, []uint32{}}
}

/*
*
32 bits in total:

	8        8        8       8

IIIIIIII JJJJJJJJ SSSSSSSS RRRRRRRR
*/
func encodeMove(i, j, shape, rowsEliminated uint32) uint32 {
	return i<<24 | j<<16 | shape<<8 | rowsEliminated
}

func decodeMove(move uint32) (i int32, j int32, shape int8, rowsEliminated int8) {
	return int32(move >> 24 & 0xff), int32(move >> 16 & 0xff), int8(move >> 8), int8(move)
}

func (tetris *Tetris) isInBound(i, j int32, shape int8) bool {
	return 0 <= i && i+bounds[shape].i < int32(len(tetris.board)) && 0 <= j && j+bounds[shape].j < int32(len(tetris.board[0]))
}

func (tetris *Tetris) isOnFloor(i, j int32, shape int8) bool {
	for _, point := range floors[shape] {
		if i+point.i < 0 || tetris.board[i+point.i][j+point.j] {
			return true
		}
	}
	return false
}

func (tetris *Tetris) isRowFilled(row int32) bool {
	for j := int32(0); j < tetris.cols; j++ {
		if !tetris.board[row][j] {
			return false
		}
	}
	return true
}

func (tetris *Tetris) doMove(i, j int32, shape int8) {
	for _, point := range shapes[shape] {
		tetris.board[i+point.i][j+point.j] = true
	}

	eliminatedRows := uint32(0)
	for offset := bounds[shape].i; offset >= 0; offset-- {
		row := i + offset
		if tetris.isRowFilled(row) {
			for ; row < tetris.rows; row++ {
				copy(tetris.board[row], tetris.board[row+1])
			}
			eliminatedRows |= 1 << offset
			tetris.score++
		}
	}
	tetris.history = append(tetris.history, encodeMove(uint32(i), uint32(j), uint32(shape), eliminatedRows))
}

func (tetris *Tetris) undoMove(move uint32) {
	i, j, shape, eliminatedRows := decodeMove(move)

	for offset := int32(0); offset <= bounds[shape].i; offset++ {
		if (eliminatedRows>>offset)&1 == 1 {
			connectedRow := i + offset
			for row := tetris.rows - 1; row >= connectedRow; row-- {
				copy(tetris.board[row+1], tetris.board[row])
			}
			for j := 0; j < int(tetris.cols); j++ {
				tetris.board[connectedRow][j] = true
			}
			tetris.score--
		}
	}

	for _, point := range shapes[shape] {
		tetris.board[i+point.i][j+point.j] = false
	}
}

func (tetris *Tetris) FindMove(shapes []int8) (int32, int8, int32) {
	max := func(a, b int32) int32 {
		if a < b {
			return b
		}
		return a
	}

	var recur func(depth int) int32
	recur = func(depth int) int32 {
		if depth == len(shapes) {
			return tetris.Evaluate()
		}

		mScore := int32(math.MinInt32)
		for rotation := int8(0); rotation < 4; rotation++ {
			shapeID := shapes[depth]/4 + rotation
			for col := int32(0); col < tetris.cols; col++ {
				if tetris.Drop(col, shapeID) {
					mScore = max(mScore, recur(depth+1))
					tetris.UnDrop()
				}
			}
		}
		return mScore
	}

	mColumn := int32(-1)
	mShapeID := int8(-1)
	mScore := int32(math.MinInt32)
	for rotation := int8(0); rotation < 4; rotation++ {
		shapeID := shapes[0]/4 + rotation
		for col := int32(0); col < tetris.cols; col++ {
			if tetris.Drop(col, shapeID) {
				if score := recur(1); score > mScore {
					mColumn = col
					mShapeID = shapeID
					mScore = score
				}
				tetris.UnDrop()
			}
		}
	}
	return mColumn, mShapeID, mScore
}

func (tetris *Tetris) Evaluate() int32 {
	return tetris.score
}

func (tetris *Tetris) Drop(col int32, shape int8) bool {
	if tetris.isInBound(tetris.rows, col, shape) {
		for i := tetris.rows; i >= 0; i-- {
			if tetris.isOnFloor(i, col, shape) {
				if i < tetris.rows {
					tetris.doMove(i, col, shape)
					return true
				}
				break
			}
		}
	}
	return false
}

func (tetris *Tetris) UnDrop() {
	tetris.undoMove(tetris.history[len(tetris.history)-1])
	tetris.history = tetris.history[:len(tetris.history)-1]
}

func (tetris *Tetris) String() string {
	builder := strings.Builder{}
	for i := tetris.rows; i >= 0; i-- {
		for j := range tetris.board[i] {
			var symbol byte
			if tetris.board[i][j] {
				symbol = 'O'
			} else {
				symbol = '.'
			}
			builder.WriteByte(symbol)
		}
		builder.WriteByte('\n')
	}
	builder.WriteString(fmt.Sprintf("Score: %v\n", tetris.score))
	return builder.String()
}
