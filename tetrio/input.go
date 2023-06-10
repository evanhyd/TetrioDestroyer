package tetrio

import (
	"log"

	"git.tcp.direct/kayos/sendkeys"
)

func SendMove(rotation, column int32) {

	const (
		kDelay = 10
	)

	keyboard, err := sendkeys.NewKBWrapWithOptions(sendkeys.NoDelay)
	if err != nil {
		log.Fatal(err)
	}

	// for r := int32(0); r < rotation; r++ {
	// 	keyboard.Type("r")
	// 	time.Sleep(kDelay * time.Millisecond)
	// }

	// keyboard.Type("aaaaaa")
	// time.Sleep(kDelay * time.Millisecond)

	// for c := int32(0); c < column; c++ {
	// 	keyboard.Type("d")
	// 	time.Sleep(kDelay * time.Millisecond)
	// }

	// keyboard.Type("w")
	// time.Sleep(kDelay * time.Millisecond)

	input := ""
	for r := int32(0); r < rotation; r++ {
		input += "r"
	}
	input += "aaaaaa"
	for c := int32(0); c < column; c++ {
		input += "d"
	}
	input += "w"
	keyboard.Type(input)
}
