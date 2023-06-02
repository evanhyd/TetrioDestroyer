package tetrio

import (
	"image"

	"github.com/kbinani/screenshot"
)

func GetTetrioImage() (*image.RGBA, error) {
	const (
		x, y          = 805, 270
		width, height = 500, 620
	)
	return screenshot.Capture(x, y, width, height)
}
