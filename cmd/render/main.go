package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/xyproto/surrender"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./render input.svg output.png")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read and parse the SVG file
	elements, err := surrender.ParseFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading or parsing file: %v\n", err)
		return
	}

	// Use a black background and a red circle
	bgColor := color.RGBA{0, 0, 0, 255}

	// Render and save the SVG elements to a PNG file
	err = surrender.RenderAndSaveSVG(elements, outputFile, bgColor)
	if err != nil {
		fmt.Printf("Error rendering and saving SVG: %v\n", err)
		return
	}

	fmt.Println("Successfully created PNG file:", outputFile)
}
