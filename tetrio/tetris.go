package tetrio

import (
	"math"
	"strings"
	"sync"
)

type Board = [kBoardColumns]uint32

type Tetris struct {
	GameState
	EvaluationStrategy
}

func NewTetris() Tetris {
	// return Tetris{EvaluationStrategy: NewEvaluationStrategy(-17.943802, -45.524887, -100.4868, -39.00439)}
	return Tetris{EvaluationStrategy: NewEvaluationStrategy(0.10710144, -25.082375, -91.383804, 70.28813)}
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

func (tetris *Tetris) SetBoard(bitmap [][]uint32) {
	for j := int32(0); j < kBoardColumns; j++ {
		tetris.board[j] = 0
		for i := int32(0); i < kBoardRows; i++ {
			tetris.board[j] |= bitmap[i][j] << i
		}
	}
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
				tetris.score++
			}
		}
		return true
	}
	return false
}

func (tetris *Tetris) FindMove(availableShapes []int32) Result {

	var search func(tetris *Tetris, depth int) float32
	search = func(tetris *Tetris, depth int) float32 {
		if depth == len(availableShapes) {
			return tetris.Evaluate(&tetris.State)
		}

		mScore := float32(-math.MaxFloat32)
		for _, variation := range variationTable[availableShapes[depth]] {
			for _, column := range shapeColumnTable[variation] {
				tetris.save()
				if tetris.MakeMove(variation, column) {
					if score := search(tetris, depth+1); score > mScore {
						mScore = score
					}
				}
				tetris.revert()
			}
		}
		return mScore
	}

	results := make(chan Result, 4*10)
	waits := sync.WaitGroup{}

	for _, variation := range variationTable[availableShapes[0]] {
		waits.Add(len(shapeColumnTable[variation]))
		for _, column := range shapeColumnTable[variation] {
			go func(tetris Tetris, shape int32, column int32) {
				tetris.save()
				if tetris.MakeMove(shape, column) {
					results <- Result{shape, column, search(&tetris, 1)}
				}
				tetris.revert()
				waits.Done()
			}(*tetris, variation, column)
		}
	}

	waits.Wait()
	close(results)

	//reset the line cleared
	tetris.score = 0
	bestResult := EmptyResult()
	for result := range results {
		if result.Eval > bestResult.Eval {
			bestResult = result
		}
	}
	return bestResult
}
