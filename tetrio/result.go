package tetrio

import "math"

type Result struct {
	Shape  int32
	Column int32
	Eval   float32
}

func EmptyResult() Result {
	return Result{EmptyShape, -1, -math.MaxFloat32}
}

func (result *Result) IsDead() bool {
	return result.Shape < 0
}
