package main

import (
	"fmt"
	"tetriodestroyer/tetrio"
)

func main() {
	tetris := tetrio.NewTetris(20, 10)
	fmt.Println(&tetris)

	column, shapeID, score := tetris.FindMove([]int8{tetrio.I0Shape, tetrio.S0Shape, tetrio.J0Shape, tetrio.L0Shape, tetrio.T0Shape})
	fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\n", column, shapeID, score)
}
