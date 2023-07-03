package surrender

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// Draw method for SvgCircle
func (c SvgCircle) Draw(img *image.RGBA, clr color.Color) {
	for y := -c.R; y <= c.R; y++ {
		for x := -c.R; x <= c.R; x++ {
			if x*x+y*y <= c.R*c.R {
				img.Set(c.Cx+x, c.Cy+y, clr)
			}
		}
	}
}

// Draw method for SvgRectangle
func (r SvgRectangle) Draw(img *image.RGBA, clr color.Color) {
	rect := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
	draw.Draw(img, rect, &image.Uniform{clr}, image.Point{}, draw.Src)
}

// Draw method for SvgPath
func (p SvgPath) Draw(img *image.RGBA, clr color.Color) {
	for _, command := range p.Commands {
		switch command.Type {
		case "M", "m", "L", "l":
			for i := 0; i < len(command.Points)-1; i++ {
				DrawLine(img, command.Points[i], command.Points[i+1], clr)
			}
		case "H", "h", "V", "v":
			for i := 0; i < len(command.Points)-1; i++ {
				DrawLine(img, command.Points[i], image.Point{X: command.Points[i+1].X, Y: command.Points[i].Y}, clr)
			}
		case "Z", "z":
			if len(command.Points) > 1 {
				DrawLine(img, command.Points[len(command.Points)-1], command.Points[0], clr)
			}
		}
	}
}

// Draw method for SvgGroup
func (g SvgGroup) Draw(img *image.RGBA, clr color.Color) {
	for _, el := range g.Elements {
		el.Draw(img, el.Color())
	}
}

// Draw method for SvgLine
func (l SvgLine) Draw(img *image.RGBA, clr color.Color) {
	DrawLine(img, image.Point{X: l.X1, Y: l.Y1}, image.Point{X: l.X2, Y: l.Y2}, clr)
}

// DrawLine function to draw a line on an image
func DrawLine(img *image.RGBA, p1, p2 image.Point, clr color.Color) {
	// Bresenham's line algorithm
	dx := abs(p2.X - p1.X)
	dy := abs(p2.Y - p1.Y)
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}
	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(p1.X, p1.Y, clr)
		if p1 == p2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

// Render function takes SVG elements and an image, and renders the elements onto the image
func Render(elements []SvgElement, img *image.RGBA) {
	for _, el := range elements {
		el.Draw(img, el.Color())
	}
}

// SavePNG function to save image as PNG
func SavePNG(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}
