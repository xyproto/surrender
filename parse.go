package surrender

import (
	"bufio"
	"image"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"github.com/beevik/etree"
	"golang.org/x/image/colornames"
)

type SvgElement interface {
	Draw(img *image.RGBA, color color.Color)
	Color() color.Color
}

// SvgCircle struct
type SvgCircle struct {
	Cx, Cy, R int
	Fill      color.Color
}

func (c SvgCircle) Color() color.Color {
	return c.Fill
}

// SvgRectangle struct
type SvgRectangle struct {
	X, Y, Width, Height int
	Fill                color.Color
}

func (r SvgRectangle) Color() color.Color {
	return r.Fill
}

// SvgPath struct
type SvgPath struct {
	Commands []PathCommand
	Fill     color.Color
}

func (p SvgPath) Color() color.Color {
	return p.Fill
}

type PathCommand struct {
	Type   string
	Points []image.Point
}

// ParseFile function
func ParseFile(filename string) ([]SvgElement, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		return nil, err
	}

	var elements []SvgElement
	for _, el := range doc.SelectElement("svg").ChildElements() {
		fillColor := getColor(el.SelectAttrValue("fill", "black"))
		switch el.Tag {
		case "circle":
			x, _ := strconv.Atoi(el.SelectAttrValue("cx", "0"))
			y, _ := strconv.Atoi(el.SelectAttrValue("cy", "0"))
			r, _ := strconv.Atoi(el.SelectAttrValue("r", "0"))
			elements = append(elements, SvgCircle{x, y, r, fillColor})

		case "rect":
			x, _ := strconv.Atoi(el.SelectAttrValue("x", "0"))
			y, _ := strconv.Atoi(el.SelectAttrValue("y", "0"))
			w, _ := strconv.Atoi(el.SelectAttrValue("width", "0"))
			h, _ := strconv.Atoi(el.SelectAttrValue("height", "0"))
			elements = append(elements, SvgRectangle{x, y, w, h, fillColor})

		case "path":
			d := el.SelectAttrValue("d", "")
			path, err := parsePath(d)
			if err != nil {
				return nil, err
			}
			path.Fill = fillColor
			elements = append(elements, path)
		}
	}

	return elements, nil
}

// getColor function
func getColor(colorName string) color.Color {
	if c, ok := colornames.Map[colorName]; ok {
		return c
	}
	return color.RGBA{0, 0, 0, 255} // default to black if color name is not recognized
}

// parsePath function
func parsePath(d string) (SvgPath, error) {
	s := bufio.NewScanner(strings.NewReader(d))
	s.Split(bufio.ScanWords)

	commands := make([]PathCommand, 0)
	var cmd PathCommand

	for s.Scan() {
		text := s.Text()
		c := text[0]
		if unicode.IsLetter(rune(c)) {
			if cmd.Type != "" {
				commands = append(commands, cmd)
			}
			cmd = PathCommand{Type: string(c)}
			text = text[1:]
		}

		if text == "" {
			continue
		}

		coords := strings.Split(text, ",")
		for i := 0; i < len(coords); i += 2 {
			x, err := strconv.Atoi(coords[i])
			if err != nil {
				return SvgPath{}, err
			}

			var y int
			if i+1 < len(coords) {
				y, err = strconv.Atoi(coords[i+1])
				if err != nil {
					return SvgPath{}, err
				}
			}

			cmd.Points = append(cmd.Points, image.Point{x, y})
		}
	}

	if err := s.Err(); err != nil {
		return SvgPath{}, err
	}

	if cmd.Type != "" {
		commands = append(commands, cmd)
	}

	return SvgPath{Commands: commands, Fill: color.RGBA{0, 0, 0, 255}}, nil
}
