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

// New structure for SvgGroup
type SvgGroup struct {
	Elements []SvgElement
	Fill     color.Color
}

func (g SvgGroup) Color() color.Color {
	return g.Fill
}

// SvgLine struct
type SvgLine struct {
	X1, Y1, X2, Y2 int
	Stroke         color.Color
}

func (l SvgLine) Color() color.Color {
	return l.Stroke
}

// ParseFile function
func ParseFile(filename string) ([]SvgElement, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		return nil, err
	}

	return parseElements(doc.SelectElement("svg").ChildElements(), color.RGBA{0, 0, 0, 255}) // Pass default color
}

func parseElements(elements []*etree.Element, defaultColor color.Color) ([]SvgElement, error) {
	var svgElements []SvgElement
	for _, el := range elements {
		fillColor := getColor(el.SelectAttrValue("fill", ""))
		if fillColor == nil {
			fillColor = defaultColor
		}

		switch el.Tag {
		case "circle":
			x, _ := strconv.Atoi(el.SelectAttrValue("cx", "0"))
			y, _ := strconv.Atoi(el.SelectAttrValue("cy", "0"))
			r, _ := strconv.Atoi(el.SelectAttrValue("r", "0"))
			svgElements = append(svgElements, SvgCircle{x, y, r, fillColor})

		case "rect":
			x, _ := strconv.Atoi(el.SelectAttrValue("x", "0"))
			y, _ := strconv.Atoi(el.SelectAttrValue("y", "0"))
			w, _ := strconv.Atoi(el.SelectAttrValue("width", "0"))
			h, _ := strconv.Atoi(el.SelectAttrValue("height", "0"))
			svgElements = append(svgElements, SvgRectangle{x, y, w, h, fillColor})

		case "line":
			x1, _ := strconv.Atoi(el.SelectAttrValue("x1", "0"))
			y1, _ := strconv.Atoi(el.SelectAttrValue("y1", "0"))
			x2, _ := strconv.Atoi(el.SelectAttrValue("x2", "0"))
			y2, _ := strconv.Atoi(el.SelectAttrValue("y2", "0"))
			strokeColor := getColor(el.SelectAttrValue("stroke", "black"))
			svgElements = append(svgElements, SvgLine{x1, y1, x2, y2, strokeColor})

		case "path":
			d := el.SelectAttrValue("d", "")
			path, err := parsePath(d)
			if err != nil {
				return nil, err
			}
			path.Fill = fillColor
			svgElements = append(svgElements, path)

		case "g":
			childElements, err := parseElements(el.ChildElements(), fillColor)
			if err != nil {
				return nil, err
			}
			svgElements = append(svgElements, SvgGroup{childElements, fillColor})
		}
	}

	return svgElements, nil
}

// getColor function, now returns a color or nil
func getColor(colorName string) color.Color {
	if colorName == "" {
		return nil
	}
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
