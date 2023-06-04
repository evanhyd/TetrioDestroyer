package tetrio

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkBounds(t *testing.B) {
	t.ResetTimer()
	for shape := int32(0); shape < kShapeSize; shape++ {
		for i := int32(0); i <= bounds[shape].i; i++ {
			for j := int32(0); j <= bounds[shape].j; j++ {
				fmt.Print("O")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func BenchmarkFloors(t *testing.B) {
	t.ResetTimer()
	for shape := int32(0); shape < kShapeSize; shape++ {
		board := [5][5]byte{}
		for i := range board {
			for j := range board[i] {
				board[i][j] = '.'
			}
		}
		for _, point := range floors[shape] {
			board[point.i][point.j] = 'O'
		}
		for i := range board {
			for j := range board[i] {
				fmt.Print(string(board[len(board)-i-1][j]))
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func BenchmarkShapes(t *testing.B) {
	fmt.Println("BenchmarkShape----------------------------------")
	for baseShape := int32(0); baseShape < kBaseShapSize; baseShape++ {
		for r := int32(0); r < rotations[baseShape]; r++ {
			shape := 4*baseShape + r
			board := [5][5]byte{}
			for i := range board {
				for j := range board[i] {
					board[i][j] = '.'
				}
			}
			for _, point := range shapes[shape] {
				board[point.i][point.j] = 'O'
			}
			for i := range board {
				for j := range board[i] {
					fmt.Print(string(board[len(board)-i-1][j]))
				}
				fmt.Println()
			}
			fmt.Println()
		}
	}
}

func BenchmarkSearchStressTest(t *testing.B) {
	tetris := NewTetris(20, 10)
	column, shapeID, score := tetris.FindMove([]int32{T0Shape, T0Shape, T0Shape, T0Shape, T0Shape})
	fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\n", column, shapeID, score)
}

func BenchmarkPlayTest(b *testing.B) {
	tetris := NewTetris(20, 10)
	shapes := []int32{rand.Int31n(kShapeSize), rand.Int31n(kShapeSize), rand.Int31n(kShapeSize), rand.Int31n(kShapeSize), rand.Int31n(kShapeSize)}

	for round := 0; round < 5; round++ {
		column, shape, score := tetris.FindMove(shapes)
		tetris.Drop(column, shape)
		fmt.Println(&tetris)
		shapes = append(shapes[1:], rand.Int31n(rand.Int31n(kShapeSize)))
		fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\nRound: %v\n", column, shape, score, round)
		if shape == -1 {
			fmt.Println("GG")
			break
		}
	}
}
