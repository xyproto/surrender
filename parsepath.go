package surrender

import (
	"errors"
	"image"
	"image/color"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// ParsePath can parse TinySVG 1.2 path attributes
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

	var currentCoords []string

	for _, c := range d {
		if unicode.IsLetter(c) {
			processCurrentCmd()
			currentCmd.Type = string(c)
			continue
		}

		switch c {
		case ' ', ',':
			if len(currentCoords) > 0 {
				coords := strings.Join(currentCoords, "")
				x, y, err := ParseCoordinates(coords)
				if err != nil {
					log.Printf("Parsing error: %v", err)
					return SvgPath{}, err
				}
				currentCmd.Points = append(currentCmd.Points, image.Point{X: x, Y: y})
				currentCoords = nil
			}
		case '-':
			// Handle "-" as a negative sign
			if len(currentCoords) > 0 {
				currentCoords = append(currentCoords, string(c))
			} else {
				currentCoords = []string{string(c)}
			}
		default:
			currentCoords = append(currentCoords, string(c))
		}
	}

	if len(currentCoords) > 0 {
		coords := strings.Join(currentCoords, "")
		x, y, err := ParseCoordinates(coords)
		if err != nil {
			log.Printf("Parsing error: %v", err)
			return SvgPath{}, err
		}
		currentCmd.Points = append(currentCmd.Points, image.Point{X: x, Y: y})
	}

	// Make sure to process the last command
	processCurrentCmd()

	if cmd := currentCmd; cmd.Type != "" {
		commands = append(commands, cmd)
	}

	return SvgPath{Commands: commands, Fill: color.RGBA{R: 0, G: 0, B: 0, A: 255}}, nil
}

// ParseCoordinates tries to parse a TinySVG 1.2 path attribute coordinate
func ParseCoordinates(coords string) (int, int, error) {
	parts := splitCoordinates(coords)
	if len(parts) < 1 || len(parts) > 2 {
		return 0, 0, errors.New("invalid coordinate format")
	}

	x, err := parseCoordinate(parts[0])
	if err != nil {
		return 0, 0, err
	}

	y := 0
	if len(parts) > 1 {
		y, err = parseCoordinate(parts[1])
		if err != nil {
			return 0, 0, err
		}
	}

	return x, y, nil
}

// splitCoordinates splits the coordinate string based on separators (" " and "-")
func splitCoordinates(coords string) []string {
	var parts []string
	var current string

	for _, c := range coords {
		if unicode.IsSpace(c) || c == '-' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		}
		current += string(c)
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// parseCoordinate parses a single coordinate value
func parseCoordinate(coord string) (int, error) {
	if coord == "" {
		return 0, errors.New("empty coordinate")
	}

	negative := false
	if coord[0] == '-' {
		negative = true
		coord = coord[1:]
	}

	value, err := strconv.Atoi(coord)
	if err != nil {
		return 0, err
	}

	if negative {
		value = -value
	}

	return value, nil
}
