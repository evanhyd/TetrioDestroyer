package tetrio

type State struct {
	board Board
	score int32
}

type GameState struct {
	State
	history [kMaxDepth]State
	depth   int32
}

func (gm *GameState) save() {
	gm.history[gm.depth] = gm.State
	gm.depth++
}

func (gm *GameState) revert() {
	gm.depth--
	gm.State = gm.history[gm.depth]
}
