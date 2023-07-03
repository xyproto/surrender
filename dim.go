package surrender

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// SVG represents the structure of the SVG file
type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

// GetSVGDimensions reads the specified TinySVG 1.2 file and returns its width and height if declared
func GetSVGDimensions(filename string) (int, int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 512, 512, fmt.Errorf("failed to read file: %w", err)
	}

	var svg SVG
	if err := xml.Unmarshal(data, &svg); err != nil {
		return 512, 512, fmt.Errorf("failed to parse SVG: %w", err)
	}

	width := 512
	height := 512
	if svg.Width != "" {
		if w, err := parseDimension(svg.Width); err == nil {
			width = w
		}
	}
	if svg.Height != "" {
		if h, err := parseDimension(svg.Height); err == nil {
			height = h
		}
	}

	return width, height, nil
}

// parseDimension parses the SVG dimension value and returns an integer
func parseDimension(dim string) (int, error) {
	var value int
	if _, err := fmt.Sscanf(dim, "%d", &value); err != nil {
		return 0, err
	}
	return value, nil
}
