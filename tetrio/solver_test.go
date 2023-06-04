package tetrio

import (
	"fmt"
	"testing"
)

func BenchmarkDimensions(t *testing.B) {
	t.ResetTimer()
	for shape := int8(0); shape < kShapeSize; shape++ {
		for i := int32(0); i <= bounds[shape].i; i++ {
			for j := int32(0); j <= bounds[shape].j; j++ {
				fmt.Print("O")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func BenchmarkFloor(t *testing.B) {
	t.ResetTimer()
	for shape := int8(0); shape < kShapeSize; shape++ {
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

func BenchmarkShape(t *testing.B) {
	fmt.Println("BenchmarkShape----------------------------------")
	for baseShape := int8(0); baseShape < kBaseShapSize; baseShape++ {
		for r := int8(0); r < rotations[baseShape]; r++ {
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

	t.ResetTimer()
	column, shapeID, score := tetris.FindMove([]int8{T0Shape, T0Shape, T0Shape, T0Shape, T0Shape})
	t.StopTimer()

	fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\n", column, shapeID, score)
}

func BenchmarkSearchCorrectnessTest(t *testing.B) {
	fmt.Println("BenchmarkSearchCorrectnessTest----------------------------------")
	tetris := NewTetris(20, 10)
	column, shapeID, score := tetris.FindMove([]int8{L90Shape})
	fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\n", column, shapeID, score)
}
