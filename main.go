package main

import (
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"os"
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

func PlayTetrio() {
	const kDepth = 5
	tetris := tetrio.NewTetris()
	for {
		if board, currentShape := tetrio.GetTetrioBoard(); board != nil {
			if shapes := tetrio.GetTetrioShapes(); len(shapes) == 5 {

				for _, shape := range shapes {
					if shape < 0 {
						img, _ := tetrio.GetTetrioShapesImage()
						file, _ := os.Create("tetrioShapes.png")
						png.Encode(file, img)
						file.Close()
					}
				}

				shapes = append([]int32{currentShape}, shapes...)
				fmt.Println(shapes)

				tetris.SetBoard(board)
				result := tetris.FindMove(shapes[:kDepth])

				if !result.IsDead() {
					tetrio.SendMove(result, currentShape)
					fmt.Printf("%+v\n", result)
					time.Sleep(10 * time.Millisecond)
				}
			}
		}
	}
}
