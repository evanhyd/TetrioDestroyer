package tetrio

import (
	"math"
	"math/bits"
	"strings"
	"sync"
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
	kBaseShapSize = 7
)

type Pair struct {
	i, j int32
}

var rotations = [kBaseShapSize]int8{2, 4, 4, 1, 4, 2, 2}

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
	{{0, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
	{{0, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},

	{{0, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {2, 1}},
	{{1, 0}, {1, 1}, {0, 2}},

	{{0, 0}, {0, 1}},
	{{0, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},

	{{0, 0}, {0, 1}},
	{{0, 0}, {0, 1}},
	{{0, 0}, {0, 1}},
	{{0, 0}, {0, 1}},

	{{1, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},

	{{0, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},

	{{1, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},
	{{1, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},
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
	board   []uint32
	rows    int32
	cols    int32
	history [][]uint32
}

func NewTetris(rows, cols int32) Tetris {
	board := make([]uint32, cols)
	return Tetris{board, rows, cols, [][]uint32{}}
}

func CopyTetris(tetris *Tetris) Tetris {
	board := make([]uint32, tetris.cols)
	copy(board, tetris.board)
	return Tetris{board, tetris.rows, tetris.cols, [][]uint32{}}
}

func (tetris *Tetris) at(row, col int32) uint32 {
	return tetris.board[col] >> row & 1
}

func (tetris *Tetris) set(row, col int32) {
	tetris.board[col] |= 1 << row
}

func (tetris *Tetris) isInBound(i, j int32, shape int8) bool {
	return 0 <= j && j+bounds[shape].j < int32(len(tetris.board))
}

func (tetris *Tetris) isRowFilled(row int32) bool {
	count := uint32(0)
	for i := range tetris.board {
		count += (tetris.board[i] >> row) & 1
	}
	return count == uint32(tetris.cols)
}

func (tetris *Tetris) doMove(i, j int32, shape int8) {
	state := make([]uint32, tetris.cols)
	copy(state, tetris.board)
	tetris.history = append(tetris.history, state)

	for _, point := range shapes[shape] {
		tetris.set(i+point.i, j+point.j)
	}

	for offset := bounds[shape].i; offset >= 0; offset-- {
		if row := i + offset; tetris.isRowFilled(row) {
			for j := range tetris.board {
				tetris.board[j] = (tetris.board[j] & (1<<row - 1)) | (tetris.board[j] >> (row + 1) << row)
			}
		}
	}
}

func (tetris *Tetris) undoMove() {
	tetris.board = tetris.history[len(tetris.history)-1]
	tetris.history = tetris.history[:len(tetris.history)-1]
}

func (tetris *Tetris) Drop(col int32, shape int8) bool {
	if tetris.isInBound(tetris.rows, col, shape) {
		floor := int32(-1)
		for _, point := range floors[shape] {
			columnHeight := int32(bits.Len32(tetris.board[col+point.j]))
			if row := columnHeight - point.i; row > floor {
				floor = columnHeight
			}
		}

		if floor < tetris.rows {
			tetris.doMove(floor, col, shape)
			return true
		}
	}
	return false
}

func (tetris *Tetris) UnDrop() {
	tetris.undoMove()
}

func (tetris *Tetris) String() string {
	builder := strings.Builder{}
	for i := tetris.rows; i >= 0; i-- {
		for j := int32(0); j < tetris.cols; j++ {
			var symbol byte
			if tetris.at(i, j) == 1 {
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

func (tetris *Tetris) FindMove(shapes []int8) (int32, int8, int32) {

	var search func(tetris *Tetris, depth int) int32
	search = func(tetris *Tetris, depth int) int32 {
		if depth == len(shapes) {
			return tetris.evaluate()
		}

		mScore := int32(math.MinInt32)
		baseShape := shapes[depth] / 4
		for r := int8(0); r < rotations[baseShape]; r++ {
			shapeID := 4*baseShape + r
			for col := int32(0); col < tetris.cols; col++ {
				if tetris.Drop(col, shapeID) {
					if score := search(tetris, depth+1); score > mScore {
						mScore = score
					}
					tetris.UnDrop()
				}
			}
		}
		return mScore
	}

	type SearchResult struct {
		column  int32
		shapeID int8
		score   int32
	}

	results := make(chan SearchResult, 4*10)
	waits := sync.WaitGroup{}
	baseShape := shapes[0] / 4
	for r := int8(0); r < rotations[baseShape]; r++ {
		for col := int32(0); col < tetris.cols; col++ {
			waits.Add(1)
			go func(tetris Tetris, column int32, shapeID int8) {
				defer waits.Done()
				if tetris.Drop(column, shapeID) {
					results <- SearchResult{column, shapeID, search(&tetris, 1)}
					tetris.UnDrop()
				}
			}(CopyTetris(tetris), col, 4*baseShape+r)
		}
	}

	waits.Wait()
	close(results)

	mColumn := int32(-1)
	mShapeID := int8(-1)
	mScore := int32(math.MinInt32)
	for result := range results {
		if result.score > mScore {
			mColumn = result.column
			mShapeID = result.shapeID
			mScore = result.score
		}
	}
	return mColumn, mShapeID, mScore
}

func (tetris *Tetris) evaluate() int32 {

	score := 0

	//zeros
	for _, col := range tetris.board {
		zeros := bits.Len32(col) - bits.OnesCount32(col)
		score -= 20 * zeros
	}

	//curvy surface
	for j := int32(0); j < tetris.cols-1; j++ {
		dH := bits.Len32(tetris.board[j]) - bits.Len32(tetris.board[j+1])
		if dH < 0 {
			dH = -dH
		}
		score -= 2 * dH * dH
	}

	return int32(score)
}
