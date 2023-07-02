package surrender

import (
	"bufio"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/beevik/etree"
)

type SvgElement interface {
	Draw(img *image.RGBA, color color.Color)
}

type SvgCircle struct {
	Cx, Cy, R int
}

type SvgRectangle struct {
	X, Y, Width, Height int
}

type PathCommand struct {
	Type   string
	Points []image.Point
}

type SvgPath struct {
	Commands []PathCommand
}

// Parse function
func Parse(filename string) ([]SvgElement, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(file); err != nil {
		return nil, err
	}

	var elements []SvgElement
	for _, el := range doc.SelectElement("svg").ChildElements() {
		switch el.Tag {
		case "circle":
			x, _ := strconv.Atoi(el.SelectAttrValue("cx", "0"))
			y, _ := strconv.Atoi(el.SelectAttrValue("cy", "0"))
			r, _ := strconv.Atoi(el.SelectAttrValue("r", "0"))
			elements = append(elements, SvgCircle{x, y, r})

		case "rect":
			x, _ := strconv.Atoi(el.SelectAttrValue("x", "0"))
			y, _ := strconv.Atoi(el.SelectAttrValue("y", "0"))
			w, _ := strconv.Atoi(el.SelectAttrValue("width", "0"))
			h, _ := strconv.Atoi(el.SelectAttrValue("height", "0"))
			elements = append(elements, SvgRectangle{x, y, w, h})

		case "path":
			d := el.SelectAttrValue("d", "")
			path, err := parsePath(d)
			if err != nil {
				return nil, err
			}
			elements = append(elements, path)
		}
	}
	return elements, nil
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

	return SvgPath{commands}, nil
}
