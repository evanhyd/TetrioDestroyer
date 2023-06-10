package tetrio

import (
	"math"
	"math/bits"
	"strings"
	"sync"
)

const (
	kBoardRows    = 23
	kBoardColumns = 10
)

const (
	EmptyShape int32 = iota - 2
	BlockerShape

	I0Shape
	I90Shape

	J0Shape
	J90Shape
	J180Shape
	J270Shape

	L0Shape
	L90Shape
	L180Shape
	L270Shape

	O0Shape

	T0Shape
	T90Shape
	T180Shape
	T270Shape

	S0Shape
	S90Shape

	Z0Shape
	Z90Shape

	kShapeSize
)

var variationTable = [kShapeSize][]int32{
	{I0Shape, I90Shape},
	{I0Shape, I90Shape},

	{J0Shape, J90Shape, J180Shape, J270Shape},
	{J0Shape, J90Shape, J180Shape, J270Shape},
	{J0Shape, J90Shape, J180Shape, J270Shape},
	{J0Shape, J90Shape, J180Shape, J270Shape},

	{L0Shape, L90Shape, L180Shape, L270Shape},
	{L0Shape, L90Shape, L180Shape, L270Shape},
	{L0Shape, L90Shape, L180Shape, L270Shape},
	{L0Shape, L90Shape, L180Shape, L270Shape},

	{O0Shape},

	{T0Shape, T90Shape, T180Shape, T270Shape},
	{T0Shape, T90Shape, T180Shape, T270Shape},
	{T0Shape, T90Shape, T180Shape, T270Shape},
	{T0Shape, T90Shape, T180Shape, T270Shape},

	{S0Shape, S90Shape},
	{S0Shape, S90Shape},

	{Z0Shape, Z90Shape},
	{Z0Shape, Z90Shape},
}

type Point struct {
	i, j int32
}

var collisionTable = [kShapeSize][]Point{
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

	{{1, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},

	{{0, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},

	{{1, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},
}

var occupancyTable = [kShapeSize][4]Point{
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

	{{1, 0}, {1, 1}, {1, 2}, {0, 1}},
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 1}},
	{{0, 0}, {1, 0}, {2, 0}, {1, 1}},

	{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
	{{1, 0}, {2, 0}, {0, 1}, {1, 1}},

	{{1, 0}, {1, 1}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
}

var shapeHeightTable = [kShapeSize]int32{
	3, 0,
	2, 1, 2, 1,
	2, 1, 2, 1,
	1,
	1, 2, 1, 2,
	1, 2,
	1, 2,
}

var shapeWidthTable = [kShapeSize]int32{
	0, 3,
	1, 2, 1, 2,
	1, 2, 1, 2,
	1,
	2, 1, 2, 1,
	2, 1,
	2, 1,
}

var shapeInColumnBoundTable [kShapeSize][kBoardColumns]bool

func init() {
	for shape := int32(0); shape < kShapeSize; shape++ {
		for j := int32(0); j < kBoardColumns; j++ {
			shapeInColumnBoundTable[shape][j] = j+shapeWidthTable[shape] < kBoardColumns
		}
	}
}

type SearchResult struct {
	Shape  int32
	Column int32
	Score  int32
}

type Tetris struct {
	board   []uint32
	history [][]uint32
}

func NewTetris() Tetris {
	board := make([]uint32, kBoardColumns)
	return Tetris{board, [][]uint32{}}
}

func CopyTetris(tetris *Tetris) Tetris {
	board := make([]uint32, kBoardColumns)
	copy(board, tetris.board)
	return Tetris{board, [][]uint32{}}
}

func (tetris *Tetris) SetBoard(bitmap [][]uint32) {
	for j := int32(0); j < kBoardColumns; j++ {
		tetris.board[j] = 0
		for i := int32(0); i < kBoardRows; i++ {
			tetris.board[j] |= bitmap[i][j] << i
		}
	}
}

func (tetris *Tetris) set(row, col int32) {
	tetris.board[col] |= 1 << row
}

func (tetris *Tetris) at(row, col int32) uint32 {
	return tetris.board[col] >> row & 1
}

func (tetris *Tetris) isRowFilled(row int32) bool {
	count := uint32(0)
	for _, column := range tetris.board {
		count += (column >> row) & 1
	}
	return count == kBoardColumns
}

func (tetris *Tetris) DoMove(shape int32, column int32) bool {

	//find the column that supports the shape
	collisionRow := int32(-1)
	for _, point := range collisionTable[shape] {
		columnHeight := int32(bits.Len32(tetris.board[column+point.j]))
		if row := columnHeight - point.i; row > collisionRow {
			collisionRow = columnHeight
		}
	}

	//placement within the board, not dead
	if collisionRow < kBoardRows {
		state := make([]uint32, kBoardColumns)
		copy(state, tetris.board)
		tetris.history = append(tetris.history, state)

		for _, point := range occupancyTable[shape] {
			tetris.set(collisionRow+point.i, column+point.j)
		}

		for rowOffset := shapeHeightTable[shape]; rowOffset >= 0; rowOffset-- {
			if i := collisionRow + rowOffset; tetris.isRowFilled(i) {
				for j := range tetris.board {
					tetris.board[j] = (tetris.board[j] & (1<<i - 1)) | (tetris.board[j] >> (i + 1) << i)
				}
			}
		}
		return true
	}
	return false
}

func (tetris *Tetris) UndoMove() {
	tetris.board = tetris.history[len(tetris.history)-1]
	tetris.history = tetris.history[:len(tetris.history)-1]
}

func (tetris *Tetris) FindMove(availableShapes []int32) SearchResult {

	var search func(tetris *Tetris, depth int) int32
	search = func(tetris *Tetris, depth int) int32 {
		if depth == len(availableShapes) {
			return tetris.evaluate()
		}

		mScore := int32(math.MinInt32)
		for _, variation := range variationTable[availableShapes[depth]] {
			for column := int32(0); column < kBoardColumns; column++ {
				if shapeInColumnBoundTable[variation][column] && tetris.DoMove(variation, column) {
					if score := search(tetris, depth+1); score > mScore {
						mScore = score
					}
					tetris.UndoMove()
				}
			}
		}
		return mScore
	}

	results := make(chan SearchResult, 4*10)
	waits := sync.WaitGroup{}

	for _, variation := range variationTable[availableShapes[0]] {
		for column := int32(0); column < kBoardColumns; column++ {
			if shapeInColumnBoundTable[variation][column] {
				waits.Add(1)
				go func(tetris Tetris, shape int32, column int32) {
					defer waits.Done()
					if tetris.DoMove(shape, column) {
						results <- SearchResult{shape, column, search(&tetris, 1)}
						tetris.UndoMove()
					}
				}(CopyTetris(tetris), variation, column)
			}
		}
	}

	waits.Wait()
	close(results)

	bestResult := SearchResult{int32(-1), EmptyShape, int32(math.MinInt32)}
	for result := range results {
		if result.Score > bestResult.Score {
			bestResult = result
		}
	}
	return bestResult
}

func (tetris *Tetris) evaluate() int32 {
	score := 0

	for j := int32(0); j < kBoardColumns; j++ {
		height := bits.Len32(tetris.board[j])
		zeros := height - bits.OnesCount32(tetris.board[j])
		score -= 20 * zeros

		if j < kBoardColumns-1 {
			dHeight := height - bits.Len32(tetris.board[j+1])
			score -= 2 * dHeight * dHeight
		}
	}

	return int32(score)
}

func (tetris *Tetris) String() string {
	symbol := [2]byte{'.', 'O'}

	builder := strings.Builder{}
	for i := int32(kBoardRows); i >= 0; i-- {
		for j := int32(0); j < kBoardColumns; j++ {
			builder.WriteByte(symbol[tetris.at(i, j)])
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}
