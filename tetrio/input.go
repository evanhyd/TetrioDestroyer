package tetrio

import (
	"log"

	"git.tcp.direct/kayos/sendkeys"
)

func SendMove(result SearchResult, currentShape int32) {
	period := int32(len(variationTable[currentShape]))
	rotation := (result.Shape - currentShape + period) % period

	keyboard, err := sendkeys.NewKBWrapWithOptions()
	if err != nil {
		log.Fatal(err)
	}

	input := ""
	for r := int32(0); r < rotation; r++ {
		input += "r"
	}
	input += "aaaaaa"
	for c := int32(0); c < result.Column; c++ {
		input += "d"
	}
	input += "w"
	keyboard.Type(input)
}
