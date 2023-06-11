package tetrio

import (
	"fmt"
	"testing"
)

func BenchmarkBounds(t *testing.B) {
	t.ResetTimer()
	for shape := int32(0); shape < kShapeSize; shape++ {
		for i := int32(0); i <= shapeHeightTable[shape]; i++ {
			for j := int32(0); j <= shapeWidthTable[shape]; j++ {
				fmt.Print("O")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func BenchmarkCollision(t *testing.B) {
	fmt.Println("Collision----------------------------")
	board := [5][5]byte{}
	for shape := int32(0); shape < kShapeSize; shape++ {
		for i := range board {
			for j := range board[i] {
				board[i][j] = '.'
			}
		}

		for _, point := range collisionTable[shape] {
			board[point.i][point.j] = 'O'
		}

		for i := range board {
			fmt.Println(string(board[len(board)-i-1][:]))
		}
		fmt.Println()
	}
}

func BenchmarkOccupancy(t *testing.B) {
	fmt.Println("Occupancy----------------------------------")
	board := [5][5]byte{}
	for shape := int32(0); shape < kShapeSize; shape++ {
		for i := range board {
			for j := range board[i] {
				board[i][j] = '.'
			}
		}

		for _, point := range occupancyTable[shape] {
			board[point.i][point.j] = 'O'
		}

		for i := range board {
			fmt.Println(string(board[len(board)-i-1][:]))
		}
		fmt.Println()
	}
}

func BenchmarkSearchStress(t *testing.B) {
	const kDepth = 5
	shapes := []int32{}
	for i := 0; i < kDepth; i++ {
		shapes = append(shapes, T0Shape)
	}

	tetris := NewTetris()

	t.ResetTimer()
	result := tetris.FindMove(shapes)
	t.StopTimer()

	fmt.Printf("Shape: %v\tColumn: %v\tEval: %v\n", result.Shape, result.Column, result.Eval)
	fmt.Println(t.Elapsed())
}

func BenchmarkSearchStressRounds(t *testing.B) {
	const kDepth = 5
	const kRound = 5

	tetris := NewTetris()
	shapes := []int32{}
	for i := 0; i < kDepth; i++ {
		shapes = append(shapes, T0Shape)
	}

	t.ResetTimer()
	for round := 0; round < kRound; round++ {
		tetris.FindMove(shapes)
	}
	t.StopTimer()
	fmt.Println(t.Elapsed())
}
