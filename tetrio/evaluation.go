package tetrio

import (
	"fmt"
)

type EvaluationStrategy struct {
	weights [4]float32
}

func NewEvaluationStrategy(aggregateHeightWeight, bumpinessWeight, holesWeight, scoreWeight float32) EvaluationStrategy {
	return EvaluationStrategy{weights: [4]float32{aggregateHeightWeight, bumpinessWeight, holesWeight, scoreWeight}}
}

func (strategy *EvaluationStrategy) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", strategy.weights[0], strategy.weights[1], strategy.weights[2], strategy.weights[3])
}

func (strategy *EvaluationStrategy) Evaluate(state *State) float32 {
	aggregateHeight := int32(0)
	bumpiness := int32(0)
	holes := int32(0)

	for _, column := range state.board {
		aggregateHeight += columnHeightTable[column]
		holes += columnHoleTable[column]
	}

	for i := 0; i < len(state.board)-1; i++ {
		heightDiff := columnHeightTable[state.board[i]] - columnHeightTable[state.board[i+1]]
		if heightDiff < 0 {
			heightDiff = -heightDiff
		}
		bumpiness += heightDiff
	}

	return strategy.weights[0]*float32(aggregateHeight) +
		strategy.weights[1]*float32(bumpiness) +
		strategy.weights[2]*float32(holes) +
		strategy.weights[3]*float32(state.score)
}

// func (strategy *EvaluationStrategy) Normalize() {
// 	magnitude := float32(math.Sqrt(float64(strategy.weights[0]*strategy.weights[0] + strategy.weights[1]*strategy.weights[1] + strategy.weights[2]*strategy.weights[2] + strategy.weights[3]*strategy.weights[3])))
// 	strategy.weights[0] /= magnitude
// 	strategy.weights[1] /= magnitude
// 	strategy.weights[2] /= magnitude
// 	strategy.weights[3] /= magnitude
// }
