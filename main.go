package main

import (
	"fmt"
	"math"
	"math/rand"
	"tetriodestroyer/tetrio"
	"time"
)

func main() {
	// training := tetrio.NewTraining("weights.txt", 1000, 0.10, 0.20, 10000)
	// training.Train()
	// PlayTest()
	PlayTetrio()
}

func PlayTest() {
	const (
		kRound = math.MaxInt
		kDepth = 4
	)

	tetris := tetrio.NewTetris()
	baseShapes := []int32{tetrio.I0Shape, tetrio.J0Shape, tetrio.L0Shape, tetrio.J0Shape, tetrio.O0Shape, tetrio.T0Shape, tetrio.S0Shape, tetrio.Z0Shape}
	shapes := []int32{}
	for i := 0; i < kDepth; i++ {
		shapes = append(shapes, baseShapes[rand.Intn(len(baseShapes))])
	}

	for round := 0; round < kRound; round++ {
		result := tetris.FindMove(shapes)
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

func PlayTetrio() {
	time.Sleep(1 * time.Second)
	const kDepth = 4
	tetris := tetrio.NewTetris()
	for {
		if board, currentShape := tetrio.GetTetrioBoard(); board != nil {
			shapes := append([]int32{currentShape}, tetrio.GetTetrioShapes()...)

			if len(shapes) >= kDepth {
				// img, _ := tetrio.GetTetrioShapesImage()
				// file, _ := os.Create("tetrioShapes.png")
				// png.Encode(file, img)
				// file.Close()
				tetris.SetBoard(board)
				result := tetris.FindMove(shapes[:kDepth])

				if !result.IsDead() {
					tetrio.SendMove(result, currentShape)
				}
			}
		}
	}
}
