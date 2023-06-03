package tetrio

import (
	"fmt"
	"testing"
)

func BenchmarkDimensions(t *testing.B) {
	t.ResetTimer()
	for shape := 0; shape < kShapeSize; shape++ {
		for i := 0; i <= bounds[shape].i; i++ {
			for j := 0; j <= bounds[shape].j; j++ {
				fmt.Print("O")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func BenchmarkFloor(t *testing.B) {
	t.ResetTimer()
	for shape := 0; shape < kShapeSize; shape++ {
		board := [5][5]byte{}
		for i := range board {
			for j := range board[i] {
				board[i][j] = '.'
			}
		}
		for _, point := range floors[shape] {
			board[point.i+1][point.j] = 'O'
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
	t.ResetTimer()
	for shape := 0; shape < kShapeSize; shape++ {
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
