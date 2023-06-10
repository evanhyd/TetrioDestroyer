package main

import (
	"fmt"
	"math"
	"math/rand"
	"tetriodestroyer/tetrio"
)

func PlayTest() {
	const (
		kRound     = math.MaxInt
		kDepth     = 5
		kShapeSize = 19
	)
	tetris := tetrio.NewTetris()
	shapes := []int32{}
	for i := 0; i < kDepth; i++ {
		shapes = append(shapes, rand.Int31n(kShapeSize))
	}

	for round := 0; round < kRound; round++ {
		result := tetris.FindMove(shapes[:len(shapes)-1])
		// fmt.Println(&tetris)
		fmt.Printf("Result %+v - %v\n", result, round)
		if result.Score != -1 {
			tetris.DoMove(result.Shape, result.Column)
			shapes = append(shapes[1:], rand.Int31n(kShapeSize))
		} else {
			fmt.Println("GG")
			break
		}
	}
}

func main() {
	PlayTest()
}

// func PlayTetrio() {
// 	tetris := tetrio.NewTetris()
// 	for {
// 		if board, currentShape := tetrio.GetTetrioBoard(); board != nil {
// 			if shapes := tetrio.GetTetrioShapes(); len(shapes) == 5 {
// 				shapes = append([]int32{currentShape}, shapes...)
// 				fmt.Println(shapes)

// 				tetris.SetBoard(board)
// 				column, shape, score := tetris.FindMove(shapes[:len(shapes)-1])

// 				if column != -1 {
// 					rotation := (shape - currentShape + 4) % 4
// 					tetrio.SendMove(rotation, column)
// 					time.Sleep(10 * time.Millisecond)
// 					fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\n", column, shape, score)
// 				} else {
// 					fmt.Println("GG")
// 				}
// 			}
// 		}
// 	}
// }
