package main

import (
	"fmt"
	"tetriodestroyer/tetrio"
)

func main() {
	tetris := tetrio.NewTetris(23, 10)
	for {
		if board, currentShape := tetrio.GetTetrioBoard(); board != nil {
			if shapes := tetrio.GetTetrioShapes(); len(shapes) == 5 {
				shapes = append([]int32{currentShape}, shapes...)
				fmt.Println(shapes)

				tetris.SetBoard(board)
				column, shape, score := tetris.FindMove(shapes[:len(shapes)-1])

				if column != -1 {
					rotation := (shape - currentShape + 4) % 4
					tetrio.SendMove(rotation, column)
					fmt.Printf("Column: %v\nShapeID: %v\nScore: %v\n", column, shape, score)
				} else {
					fmt.Println("GG")
				}
			} else {
				// log.Println("can not find shapes")
			}
		} else {
			// log.Println("can find current shape")
		}
	}
}
