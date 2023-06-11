package main

import (
	"fmt"
	"math"
	"math/rand"
	"tetriodestroyer/tetrio"
	"time"
)

func PlayTest() {
	const (
		kRound = math.MaxInt
		kDepth = 5
	)

	tetris := tetrio.NewTetris()
	baseShapes := []int32{tetrio.I0Shape, tetrio.J0Shape, tetrio.L0Shape, tetrio.J0Shape, tetrio.O0Shape, tetrio.T0Shape, tetrio.S0Shape, tetrio.Z0Shape}
	shapes := []int32{}
	for i := 0; i < kDepth; i++ {
		shapes = append(shapes, baseShapes[rand.Intn(len(baseShapes))])
	}

	for round := 0; round < kRound; round++ {
		result := tetris.FindMove(shapes)
		// fmt.Println(&tetris)
		fmt.Printf("Result %+v - %v\n", result, round)
		if !result.IsDead() {
			tetris.MakeMove(result.Shape, result.Column)
			shapes = append(shapes[1:], baseShapes[rand.Intn(len(baseShapes))])
		} else {
			fmt.Println("GG")
			break
		}
	}
}

func main() {
	// training := tetrio.NewTraining("weights.txt", 100, 0.05, 0.02, 50000)
	// training.Train()
	// PlayTest()
	PlayTetrio()
}

func PlayTetrio() {
	const kDepth = 5
	tetris := tetrio.NewTetris()
	for {
		if board, currentShape := tetrio.GetTetrioBoard(); board != nil {
			if shapes := tetrio.GetTetrioShapes(); len(shapes) == 5 {
				shapes = append([]int32{currentShape}, shapes...)
				fmt.Println(shapes)

				tetris.SetBoard(board)
				result := tetris.FindMove(shapes[:kDepth])

				if !result.IsDead() {
					tetrio.SendMove(result, currentShape)
					fmt.Printf("%+v\n", result)
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	}
}
