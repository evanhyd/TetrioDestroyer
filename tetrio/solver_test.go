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

func BenchmarkSearchStressTest(t *testing.B) {
	const kDepth = 5
	shapes := []int32{}
	for i := 0; i < kDepth; i++ {
		shapes = append(shapes, T0Shape)
	}

	tetris := NewTetris()
	result := tetris.FindMove(shapes)
	fmt.Printf("Shape: %v\nColumn: %v\nScore: %v\n", result.Shape, result.Column, result.Score)
}
