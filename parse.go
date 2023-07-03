package surrender

import (
	"image"
	"image/color"
	"log"
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

// ParseFile will try to parse the given TinySVG 1.2 file into a slice of SvgElements
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
		fillColor := GetColor(el.SelectAttrValue("fill", ""))
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
			strokeColor := GetColor(el.SelectAttrValue("stroke", "black"))
			svgElements = append(svgElements, SvgLine{x1, y1, x2, y2, strokeColor})

		case "path":
			d := el.SelectAttrValue("d", "")
			path, err := ParsePath(d)
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

func GetColor(colorStr string) color.Color {
	// If the string is empty, return nil
	if colorStr == "" {
		return nil
	}

	// If the string is a color name, return the corresponding color
	if c, ok := colornames.Map[colorStr]; ok {
		return c
	}

	// If the string is an RGB hex code on the form "#fff", convert it to RGBA color
	if len(colorStr) == 4 && colorStr[0] == '#' {
		r, _ := strconv.ParseInt(strings.Repeat(string(colorStr[1]), 2), 16, 64)
		g, _ := strconv.ParseInt(strings.Repeat(string(colorStr[2]), 2), 16, 64)
		b, _ := strconv.ParseInt(strings.Repeat(string(colorStr[3]), 2), 16, 64)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}

	// If the string is an RGB hex code on the form "#ffffff", convert it to RGBA color
	if len(colorStr) == 7 && colorStr[0] == '#' {
		r, _ := strconv.ParseInt(colorStr[1:3], 16, 64)
		g, _ := strconv.ParseInt(colorStr[3:5], 16, 64)
		b, _ := strconv.ParseInt(colorStr[5:7], 16, 64)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}

	// If the string is an RGB functional notation on the form "rgb(255, 255, 255)",
	// convert it to RGBA color
	if strings.HasPrefix(colorStr, "rgb(") && strings.HasSuffix(colorStr, ")") {
		rgbStr := strings.TrimPrefix(colorStr, "rgb(")
		rgbStr = strings.TrimSuffix(rgbStr, ")")
		rgbValues := strings.Split(rgbStr, ",")
		if len(rgbValues) == 3 {
			r, _ := strconv.Atoi(strings.TrimSpace(rgbValues[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(rgbValues[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(rgbValues[2]))
			return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		}
	}

	// If the string is not recognized, return black
	return color.RGBA{0, 0, 0, 255}
}

// ParsePath will try to parse a TinySVG 1.2 path attribute string
func ParsePath(d string) (SvgPath, error) {
	var commands []PathCommand
	var currentCmd PathCommand

	// Function to process current command before starting a new one
	processCurrentCmd := func() {
		if currentCmd.Type != "" {
			commands = append(commands, currentCmd)
		}
		currentCmd = PathCommand{}
	}

	for _, cmdStr := range strings.Fields(d) {
		for _, c := range cmdStr {
			// If it's a letter, it's a command type
			if unicode.IsLetter(c) {
				processCurrentCmd()
				currentCmd.Type = string(c)
				continue
			}

			// If it's not a letter, it should be a part of coordinate
			coordPart := string(c)
			switch currentCmd.Type {
			case "h", "H", "v", "V":
				val, err := strconv.Atoi(coordPart)
				if err != nil {
					log.Printf("Parsing error: %v", err)
					return SvgPath{}, err
				}
				currentCmd.Points = append(currentCmd.Points, image.Point{X: val, Y: val})
			default:
				coords := strings.Split(coordPart, ",")
				for i := 0; i < len(coords); i += 2 {
					x, err := strconv.Atoi(coords[i])
					if err != nil {
						log.Printf("Parsing error: %v", err)
						return SvgPath{}, err
					}
					var y int
					if i+1 < len(coords) {
						y, err = strconv.Atoi(coords[i+1])
						if err != nil {
							log.Printf("Parsing error: %v", err)
							return SvgPath{}, err
						}
					}
					currentCmd.Points = append(currentCmd.Points, image.Point{X: x, Y: y})
				}
			}
		}
	}

	// Make sure to process the last command
	processCurrentCmd()

	if cmd := currentCmd; cmd.Type != "" {
		commands = append(commands, cmd)
	}

	return SvgPath{Commands: commands, Fill: color.RGBA{0, 0, 0, 255}}, nil
}
