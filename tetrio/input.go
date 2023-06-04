package tetrio

import (
	"log"

	"git.tcp.direct/kayos/sendkeys"
)

func SendMove(rotation, column int32) {
	keyboard, err := sendkeys.NewKBWrapWithOptions()
	if err != nil {
		log.Fatal(err)
	}

	// const kDelay = 200
	// for r := int32(0); r < rotation; r++ {
	// 	keyboard.Type("r")
	// 	time.Sleep(kDelay * time.Millisecond)
	// }
	// for a := 0; a < 5; a++ {
	// 	keyboard.Type("a")
	// 	time.Sleep(kDelay * time.Millisecond)
	// }
	// for c := int32(0); c < column; c++ {
	// 	keyboard.Type("d")
	// 	time.Sleep(kDelay * time.Millisecond)
	// }
	// keyboard.Type("w")

	input := ""
	for r := int32(0); r < rotation; r++ {
		input += "r"
	}
	input += "aaaaa"
	for c := int32(0); c < column; c++ {
		input += "d"
	}
	input += "w"
	keyboard.Type(input)
}
