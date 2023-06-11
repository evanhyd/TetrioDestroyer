package tetrio

import "math"

type SearchResult struct {
	Shape  int32
	Column int32
	Eval   float32
}

func NoResult() SearchResult {
	return SearchResult{EmptyShape, -1, -math.MaxFloat32}
}

func (searchResult *SearchResult) IsDead() bool {
	return searchResult.Shape < 0
}
