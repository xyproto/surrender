package surrender

import (
	"image"
	"image/color"
	"image/draw"
)

// Utility function to compute absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Create new colored image
func NewColoredImage(width, height int, clr color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{clr}, image.Point{}, draw.Src)
	return img
}
