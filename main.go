package main

import (
	"fmt"
	"math/rand"
	"tetriodestroyer/tetrio"
)

func RandShape() int8 {
	return int8(rand.Intn(28))
}

func main() {

	tetris := tetrio.NewTetris(20, 10)
	shapes := []int8{RandShape(), RandShape(), RandShape(), RandShape(), RandShape()}
	round := 0
	for {
		column, shape, score := tetris.FindMove(shapes)
		tetris.Drop(column, shape)
		fmt.Println(&tetris)
		shapes = append(shapes[1:], RandShape())
		fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\nRound: %v\n", column, shape, score, round)
		round++
		if shape == -1 {
			fmt.Println("GG")
			break
		}
	}
}
