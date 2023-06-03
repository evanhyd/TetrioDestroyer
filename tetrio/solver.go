package tetrio

import (
	"fmt"
	"strings"
)

const (
	I0Shape = iota
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
	i, j int
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
	board [][]bool
	rows  int
	cols  int
	score int
}

func NewTetris(rows, cols int) Tetris {
	board := make([][]bool, rows+5)
	for i := range board {
		board[i] = make([]bool, cols)
	}
	return Tetris{board, rows, cols, 0}
}

func CopyTetris(tetris *Tetris) Tetris {
	board := make([][]bool, len(tetris.board))
	for i := range board {
		board[i] = make([]bool, len(tetris.board[i]))
		copy(board[i], tetris.board[i])
	}
	return Tetris{board, tetris.rows, tetris.cols, tetris.score}
}

func (tetris *Tetris) isInBound(i int, j int, shape int) bool {
	return 0 <= i && i+bounds[shape].i < len(tetris.board) && 0 <= j && j+bounds[shape].j < len(tetris.board[0])
}

func (tetris *Tetris) isOnFloor(i int, j int, shape int) bool {
	for _, point := range floors[shape] {
		if i+point.i < 0 || tetris.board[i+point.i][j+point.j] {
			return true
		}
	}
	return false
}

func (tetris *Tetris) isRowFilled(row int) bool {
	for j := 0; j < tetris.cols; j++ {
		if !tetris.board[row][j] {
			return false
		}
	}
	return true
}

func (tetris *Tetris) removeRow(row int) {
	for i := row; i < tetris.rows; i++ {
		copy(tetris.board[i], tetris.board[i+1])
	}
}

func (tetris *Tetris) putShape(i int, j int, shape int) {
	for _, point := range shapes[shape] {
		tetris.board[i+point.i][j+point.j] = true
	}
	for tetris.isRowFilled(i) {
		tetris.removeRow(i)
		tetris.score += 1
	}
}

func (tetris *Tetris) Drop(col int, shape int) bool {
	if !tetris.isInBound(tetris.rows, col, shape) {
		return false
	}

	for i := tetris.rows; i >= 0; i-- {
		if tetris.isOnFloor(i, col, shape) {
			if i < tetris.rows {
				tetris.putShape(i, col, shape)
				return true
			}
			return false
		}
	}
	return false
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
