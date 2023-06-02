package main

import (
	"fmt"
	"tetriodestroyer/tetrio"
)

func main() {
	tetris := tetrio.NewTetris(10, 10)
	fmt.Println(&tetris)

	for {
		c, s := 0, 0
		fmt.Scan(&c, &s)
		tetris.Drop(c, s)
		fmt.Println(&tetris)
	}
}
