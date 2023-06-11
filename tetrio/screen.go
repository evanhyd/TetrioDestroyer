package tetrio

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/kbinani/screenshot"
)

var shapeColors = []color.RGBA{
	{50, 180, 130, 255}, //I90
	{83, 59, 206, 255},  //J90
	{190, 120, 70, 255}, //L270
	{180, 153, 50, 255}, //O0
	{165, 62, 155, 255}, //T180
	{131, 180, 50, 255}, //S0
	{194, 63, 70, 255},  //Z0
	{68, 68, 68, 255},   //Blocker
	{0, 0, 0, 255},      //Empty
}

var colorShapeID = []int32{I90Shape, J90Shape, L270Shape, O0Shape, T180Shape, S0Shape, Z0Shape, BlockerShape, EmptyShape}

func getShape(c color.RGBA) int32 {
	colorID := EmptyShape
	mDiff := int32(math.MaxInt32)
	for i := range shapeColors {
		dR := int32(c.R) - int32(shapeColors[i].R)
		dG := int32(c.G) - int32(shapeColors[i].G)
		dB := int32(c.B) - int32(shapeColors[i].B)
		diff := dR*dR + dG*dG + dB*dB
		if diff < mDiff {
			mDiff = diff
			colorID = int32(i)
		}
	}
	return colorShapeID[colorID]
}

func getTetrioBoardImage() (*image.RGBA, error) {
	const (
		kX, kY          = 805, 175
		kWidth, kHeight = 310, 715
	)
	return screenshot.Capture(kX, kY, kWidth, kHeight)
}

func getTetrioShapesImage() (*image.RGBA, error) {
	const (
		kX, kY          = 1150, 300
		kWidth, kHeight = 140, 455
	)
	return screenshot.Capture(kX, kY, kWidth, kHeight)
}

func GetTetrioBoard() ([][]uint32, int32) {
	const (
		kRow    = 23
		kColumn = 10
	)

	img, err := getTetrioBoardImage()
	if err != nil {
		log.Fatal(err)
	}

	blockY := img.Rect.Dy() / kRow
	blockX := img.Rect.Dx() / kColumn

	//get the current shape
	currentShape := EmptyShape
	y := 0

Loop:
	for y = blockY / 2; y < img.Rect.Dy(); y += blockY {
		for x := blockX / 2; x < img.Rect.Dx(); x += blockX {
			if shape := getShape(img.RGBAAt(x, y)); shape != EmptyShape {
				currentShape = shape
				break Loop
			}
		}
	}

	//fail to identify the current shape
	if currentShape < 0 {
		return nil, EmptyShape
	}

	//get the board
	board := make([][]uint32, kRow)
	for i := range board {
		board[i] = make([]uint32, kColumn)
	}

	for y += blockY * int(shapeHeightTable[currentShape]+1); y < img.Rect.Dy(); y += blockY {
		i := kRow - (y / blockY) - 1
		for x := blockX / 2; x < img.Rect.Dx(); x += blockX {
			if getShape(img.RGBAAt(x, y)) != EmptyShape {
				board[i][x/blockX] = 1
			}
		}
	}

	return board, currentShape
}

func GetTetrioShapes() []int32 {
	const (
		kRow    = 5
		kColumn = 5
	)

	img, err := getTetrioShapesImage()
	if err != nil {
		log.Fatal(err)
	}
	blockY := img.Rect.Dy() / kRow
	blockX := img.Rect.Dx() / kColumn

	shapes := []int32{}
	for y := blockY / 2; y < img.Rect.Dy(); y += blockY {
		for x := blockX / 2; x < img.Rect.Dx(); x += blockX {
			if shape := getShape(img.RGBAAt(x, y)); shape != EmptyShape {
				shapes = append(shapes, shape)
				break
			}
		}
	}
	return shapes
}
