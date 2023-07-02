package surrender

import (
	"image"
	"image/color"
	"image/draw"
)

// Define image width and height
const (
	Width  = 500
	Height = 500
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

// RenderAndSaveSVG takes SVG elements and color, creates an image, renders the elements onto the image and saves it as PNG
func RenderAndSaveSVG(elements []SvgElement, filename string, bgColor, elementColor color.Color) error {
	img := NewColoredImage(Width, Height, bgColor)
	Render(elements, img, elementColor)

	err := SavePNG(img, filename)
	if err != nil {
		return err
	}
	return nil
}
