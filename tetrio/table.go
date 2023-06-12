package tetrio

import "math/bits"

type Point struct {
	i, j int32
}

const (
	kBoardRows       = 20
	kBoardColumns    = 10
	kMaxDepth        = 8
	kColumnTableSize = 1 << (kBoardRows + 3) //+3 because I0 Shape stacking on top
)

const (
	EmptyShape int32 = iota - 2
	BlockerShape

	I0Shape
	I90Shape

	J0Shape
	J90Shape
	J180Shape
	J270Shape

	L0Shape
	L90Shape
	L180Shape
	L270Shape

	O0Shape

	T0Shape
	T90Shape
	T180Shape
	T270Shape

	S0Shape
	S90Shape

	Z0Shape
	Z90Shape

	kShapeSize
)

var nameTable = [kShapeSize]string{
	"I0Shape",
	"I90Shape",

	"J0Shape",
	"J90Shape",
	"J180Shape",
	"J270Shape",

	"L0Shape",
	"L90Shape",
	"L180Shape",
	"L270Shape",

	"O0Shape",

	"T0Shape",
	"T90Shape",
	"T180Shape",
	"T270Shape",

	"S0Shape",
	"S90Shape",

	"Z0Shape",
	"Z90Shape",
}

var variationTable = [kShapeSize][]int32{
	{I0Shape, I90Shape},
	{I0Shape, I90Shape},

	{J0Shape, J90Shape, J180Shape, J270Shape},
	{J0Shape, J90Shape, J180Shape, J270Shape},
	{J0Shape, J90Shape, J180Shape, J270Shape},
	{J0Shape, J90Shape, J180Shape, J270Shape},

	{L0Shape, L90Shape, L180Shape, L270Shape},
	{L0Shape, L90Shape, L180Shape, L270Shape},
	{L0Shape, L90Shape, L180Shape, L270Shape},
	{L0Shape, L90Shape, L180Shape, L270Shape},

	{O0Shape},

	{T0Shape, T90Shape, T180Shape, T270Shape},
	{T0Shape, T90Shape, T180Shape, T270Shape},
	{T0Shape, T90Shape, T180Shape, T270Shape},
	{T0Shape, T90Shape, T180Shape, T270Shape},

	{S0Shape, S90Shape},
	{S0Shape, S90Shape},

	{Z0Shape, Z90Shape},
	{Z0Shape, Z90Shape},
}

var collisionTable = [kShapeSize][]Point{
	{{0, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},

	{{0, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {2, 1}},
	{{1, 0}, {1, 1}, {0, 2}},

	{{0, 0}, {0, 1}},
	{{0, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},

	{{0, 0}, {0, 1}},

	{{1, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},

	{{0, 0}, {0, 1}, {1, 2}},
	{{1, 0}, {0, 1}},

	{{1, 0}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 1}},
}

var occupancyTable = [kShapeSize][4]Point{
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},

	{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
	{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
	{{1, 0}, {1, 1}, {1, 2}, {0, 2}},

	{{0, 0}, {0, 1}, {1, 0}, {2, 0}},
	{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {1, 1}, {0, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 2}},

	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},

	{{1, 0}, {1, 1}, {1, 2}, {0, 1}},
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{0, 0}, {0, 1}, {0, 2}, {1, 1}},
	{{0, 0}, {1, 0}, {2, 0}, {1, 1}},

	{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
	{{1, 0}, {2, 0}, {0, 1}, {1, 1}},

	{{1, 0}, {1, 1}, {0, 1}, {0, 2}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
}

var shapeHeightTable = [kShapeSize]int32{
	3, 0,
	2, 1, 2, 1,
	2, 1, 2, 1,
	1,
	1, 2, 1, 2,
	1, 2,
	1, 2,
}

var shapeWidthTable = [kShapeSize]int32{
	0, 3,
	1, 2, 1, 2,
	1, 2, 1, 2,
	1,
	2, 1, 2, 1,
	2, 1,
	2, 1,
}

var shapeColumnTable = [kShapeSize][]int32{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{0, 1, 2, 3, 4, 5, 6},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
}

var columnHeightTable = [kColumnTableSize]int32{}
var columnHoleTable = [kColumnTableSize]int32{}

func init() {
	for n := uint32(0); n < kColumnTableSize; n++ {
		columnHeightTable[n] = int32(bits.Len32(n))
		columnHoleTable[n] = columnHeightTable[n] - int32(bits.OnesCount32(n))
	}
}
