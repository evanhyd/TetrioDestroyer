package tetrio

import (
	"log"
	"time"

	"git.tcp.direct/kayos/sendkeys"
)

func SendMove(rotation, column int32) {
	const kDelay = 500
	keyboard, err := sendkeys.NewKBWrapWithOptions(sendkeys.NoDelay)
	if err != nil {
		log.Fatal(err)
	}

	for r := int32(0); r < rotation; r++ {
		keyboard.Type("r")
		time.Sleep(kDelay * time.Millisecond)
	}
	for a := 0; a < 5; a++ {
		keyboard.Type("a")
		time.Sleep(kDelay * time.Millisecond)
	}
	for c := int32(0); c < column; c++ {
		keyboard.Type("d")
		time.Sleep(kDelay * time.Millisecond)
	}
	keyboard.Type("w")
}
