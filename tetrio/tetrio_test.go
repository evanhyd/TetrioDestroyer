package tetrio

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func BenchmarkGetTetrioBoardImage(b *testing.B) {
	img, err := getTetrioBoardImage()
	if err != nil {
		b.Error(err)
	}
	file, err := os.Create("tetrioBoard.png")
	if err != nil {
		b.Error(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func BenchmarkGetTetrioShapesImage(b *testing.B) {
	img, err := GetTetrioShapesImage()
	if err != nil {
		b.Error(err)
	}
	file, err := os.Create("tetrioShapes.png")
	if err != nil {
		b.Error(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func BenchmarkGetTetrioBoard(b *testing.B) {
	fmt.Println("Printing tetrio board-------------------")
	board, currentShape := GetTetrioBoard()
	if board != nil {
		tetris := NewTetris()
		tetris.SetBoard(board)
		fmt.Println(&tetris)
		fmt.Println("Shape", currentShape)
	} else {
		fmt.Println("can not find the current shape")
	}
}

func BenchmarkGetTetrioShapes(b *testing.B) {
	fmt.Println("Printing tetrio shapes-------------------")
	shapes := GetTetrioShapes()
	for _, shape := range shapes {
		fmt.Println(nameTable[shape])
	}
}
