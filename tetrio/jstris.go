package tetrio

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"git.tcp.direct/kayos/sendkeys"
	"github.com/kbinani/screenshot"
)

var shapeColors = []color.RGBA{
	{15, 155, 215, 255},  //I
	{33, 65, 198, 255},   //J
	{227, 91, 2, 255},    //L
	{227, 159, 2, 255},   //O
	{175, 41, 138, 255},  //T
	{89, 177, 1, 255},    //S
	{215, 15, 55, 255},   //Z
	{106, 106, 106, 255}, //Blocker
	{0, 0, 0, 255},       //Empty
}

var colorShapeID = []int32{I90Shape, J90Shape, L270Shape, O0Shape, T180Shape, S0Shape, Z0Shape, BlockerShape, EmptyShape}

func predictShapeByColor(c color.RGBA) int32 {
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
		kX, kY          = 503, 132
		kWidth, kHeight = 300, 600
	)
	return screenshot.Capture(kX, kY, kWidth, kHeight)
}

func GetTetrioShapesImage() (*image.RGBA, error) {
	const (
		kX, kY          = 808, 134
		kWidth, kHeight = 160, 422
	)
	return screenshot.Capture(kX, kY, kWidth, kHeight)
}

func GetTetrioBoard() ([][]uint32, int32) {
	const (
		kRow    = 20
		kColumn = 10
	)

	img, err := getTetrioBoardImage()
	if err != nil {
		log.Fatal(err)
	}
	// file, _ := os.Create("board.png")
	// png.Encode(file, img)
	// file.Close()

	blockY := img.Rect.Dy() / kRow
	blockX := img.Rect.Dx() / kColumn

	//get the current shape
	currentShape := EmptyShape
	y := 0

Loop:
	for y = blockY / 2; y < img.Rect.Dy(); y += blockY {
		for x := blockX / 2; x < img.Rect.Dx(); x += blockX {
			if shape := predictShapeByColor(img.RGBAAt(x, y)); shape != EmptyShape {
				currentShape = shape
				fmt.Println(img.RGBAAt(x, y))
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
			if predictShapeByColor(img.RGBAAt(x, y)) != EmptyShape {
				board[i][x/blockX] = 1
			}
		}
	}

	return board, currentShape
}

func GetTetrioShapes() []int32 {
	const (
		kRow    = 5
		kColumn = 10
	)

	img, err := GetTetrioShapesImage()
	if err != nil {
		log.Fatal(err)
	}
	blockY := img.Rect.Dy() / kRow
	blockX := img.Rect.Dx() / kColumn

	shapes := []int32{}
	for y := blockY / 2; y < img.Rect.Dy(); y += blockY {
		for x := blockX / 2; x < img.Rect.Dx(); x += blockX {
			if shape := predictShapeByColor(img.RGBAAt(x, y)); shape != EmptyShape {
				shapes = append(shapes, shape)
				break
			}
		}
	}
	return shapes
}

func SendMove(result Result, currentShape int32) {
	period := int32(len(variationTable[currentShape]))
	rotation := (result.Shape - currentShape + period) % period

	keyboard, err := sendkeys.NewKBWrapWithOptions()
	if err != nil {
		log.Fatal(err)
	}

	input := ""
	for r := int32(0); r < rotation; r++ {
		input += "e"
	}
	input += "aaaaaaaa"
	for c := int32(0); c < result.Column; c++ {
		input += "d"
	}
	input += "w"
	keyboard.Type(input)
}
