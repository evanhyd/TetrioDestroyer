package tetrio

import (
	"math"
	"math/bits"
	"strings"
	"sync"
)

const (
	kBoardRows    = 20
	kBoardColumns = 10
	kMaxDepth     = 8
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

var shapeColumnTable = [kShapeSize][]int32{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{0, 1, 2, 3, 4, 5, 6},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
}

const kTableSize = 1 << (kBoardRows + 3)

var columnHeightTable = [kTableSize]int32{}
var columnHoleTable = [kTableSize]int32{}

func init() {
	for n := uint32(0); n < kTableSize; n++ {
		columnHeightTable[n] = int32(bits.Len32(n))
		columnHoleTable[n] = columnHeightTable[n] - int32(bits.OnesCount32(n))
	}
}

type SearchResult struct {
	Shape  int32
	Column int32
	Score  int32
}

type Tetris struct {
	board   [kBoardColumns]uint32
	history [kMaxDepth][kBoardColumns]uint32
	depth   int32
}

func NewTetris() Tetris {
	return Tetris{}
}

func CopyTetris(tetris *Tetris) Tetris {
	return Tetris{board: tetris.board}
}

func (tetris *Tetris) SetBoard(bitmap [][]uint32) {
	for j := int32(0); j < kBoardColumns; j++ {
		tetris.board[j] = 0
		for i := int32(0); i < kBoardRows; i++ {
			tetris.board[j] |= bitmap[i][j] << i
		}
	}
}

func (tetris *Tetris) SaveState() {
	tetris.history[tetris.depth] = tetris.board
	tetris.depth++
}

func (tetris *Tetris) RollbackState() {
	tetris.depth--
	tetris.board = tetris.history[tetris.depth]
}

func (tetris *Tetris) isRowFilled(row int32) bool {
	count := uint32(0)
	for _, column := range tetris.board {
		count += (column >> row) & 1
	}
	return count == kBoardColumns
}

func (tetris *Tetris) MakeMove(shape int32, column int32) bool {

	//find the column that supports the shape
	collisionRow := int32(-1)
	for _, point := range collisionTable[shape] {
		columnHeight := columnHeightTable[tetris.board[column+point.j]]
		if row := columnHeight - point.i; row > collisionRow {
			collisionRow = columnHeight
		}
	}

	//placement within the board, not dead
	if collisionRow < kBoardRows {

		for _, point := range occupancyTable[shape] {
			tetris.board[column+point.j] |= 1 << (collisionRow + point.i)
		}

		for rowOffset := shapeHeightTable[shape]; rowOffset >= 0; rowOffset-- {
			if row := collisionRow + rowOffset; tetris.isRowFilled(row) {
				for j, column := range tetris.board {
					tetris.board[j] = (column & (1<<row - 1)) | (column >> (row + 1) << row)
				}
			}
		}
		return true
	}
	return false
}

func (tetris *Tetris) FindMove(availableShapes []int32) SearchResult {

	var search func(tetris *Tetris, depth int) int32
	search = func(tetris *Tetris, depth int) int32 {
		if depth == len(availableShapes) {
			return tetris.evaluate()
		}

		mScore := int32(math.MinInt32)
		for _, variation := range variationTable[availableShapes[depth]] {
			for _, column := range shapeColumnTable[variation] {
				tetris.SaveState()
				if tetris.MakeMove(variation, column) {
					if score := search(tetris, depth+1); score > mScore {
						mScore = score
					}
				}
				tetris.RollbackState()
			}
		}
		return mScore
	}

	results := make(chan SearchResult, 4*10)
	waits := sync.WaitGroup{}

	for _, variation := range variationTable[availableShapes[0]] {
		waits.Add(len(shapeColumnTable[variation]))
		for _, column := range shapeColumnTable[variation] {
			go func(tetris Tetris, shape int32, column int32) {
				tetris.SaveState()
				if tetris.MakeMove(shape, column) {
					results <- SearchResult{shape, column, search(&tetris, 1)}
				}
				tetris.RollbackState()
				waits.Done()
			}(CopyTetris(tetris), variation, column)
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
	score := int32(0)

	for j := int32(0); j < kBoardColumns; j++ {
		score -= 20 * columnHoleTable[tetris.board[j]]
		if j < kBoardColumns-1 {
			dHeight := columnHeightTable[tetris.board[j]] - columnHeightTable[tetris.board[j+1]]
			score -= 2 * dHeight * dHeight
		}
	}

	return score
}

func (tetris *Tetris) String() string {
	symbol := [2]byte{'.', 'O'}

	builder := strings.Builder{}
	for i := int32(kBoardRows); i >= 0; i-- {
		for j := int32(0); j < kBoardColumns; j++ {
			builder.WriteByte(symbol[tetris.board[j]>>i&1])
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}
