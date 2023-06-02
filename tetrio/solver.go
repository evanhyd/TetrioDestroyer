package tetrio

import (
	"strings"
)

type Tetris struct {
	board [][]bool
}

func NewTetris(row, col int) Tetris {
	board := make([][]bool, row)
	for i := range board {
		board[i] = make([]bool, col)
	}
	return Tetris{board}
}

func CopyTetris(tetris *Tetris) Tetris {
	board := make([][]bool, tetris.rows())
	for i := range board {
		board[i] = make([]bool, tetris.cols())
		copy(board[i], tetris.board[i])
	}
	return Tetris{board}
}

func (tetris *Tetris) rows() int {
	return len(tetris.board)
}

func (tetris *Tetris) cols() int {
	return len(tetris.board[0])
}

func (tetris *Tetris) removeRow(row int) {
	lastRow := tetris.rows() - 1
	for i := row; i < lastRow; i++ {
		copy(tetris.board[i], tetris.board[i+1])
	}
	for i := range tetris.board[lastRow] {
		tetris.board[lastRow][i] = false
	}
}

func (tetris *Tetris) checkRowFilled(row int) bool {
	filled := true
	for j := range tetris.board[row] {
		if !tetris.board[row][j] {
			filled = false
			break
		}
	}
	return filled
}

func (tetris *Tetris) DropDot(col int) {
	for i := tetris.rows() - 1; i >= 0; i-- {
		if i-1 < 0 || tetris.board[i-1][col] {
			tetris.board[i][col] = true
		}
	}
}

func (tetris *Tetris) Drop(col int, shape int) {
	for i := tetris.rows() - 1; i >= 0; i-- {
		if FitDimension(tetris, i, col, shape) && OnFloor(tetris, i, col, shape) {
			*tetris = CombineShape(tetris, i, col, shape)
			return
		}
	}
}

func (tetris *Tetris) String() string {
	builder := strings.Builder{}
	for i := tetris.rows() - 1; i >= 0; i-- {
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
	return builder.String()
}
