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

func (gm *GameState) saveState() {
	gm.history[gm.depth] = gm.State
	gm.depth++
}

func (gm *GameState) revertState() {
	gm.depth--
	gm.State = gm.history[gm.depth]
}
