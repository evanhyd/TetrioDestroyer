package main

import (
	"fmt"
	"tetriodestroyer/tetrio"
)

func main() {
	tetris := tetrio.NewTetris(20, 10)
	fmt.Println(&tetris)

	for {
		c, s := 0, 0
		fmt.Scan(&c, &s)
		if !tetris.Drop(c, s) {
			fmt.Println("You lose")
		}
		fmt.Println(&tetris)
	}
}
