package surrender

import (
	"image/color"
	"os"
	"testing"
)

func TestRenderAndSaveSVG(t *testing.T) {
	circle := SvgCircle{Cx: 250, Cy: 250, R: 50}
	elements := []SvgElement{circle}

	filename := "test.png"

	// Delete the test file if it exists
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	// Use a black background and a red circle
	bgColor := color.RGBA{0, 0, 0, 255}
	circleColor := color.RGBA{255, 0, 0, 255}

	err := RenderAndSaveSVG(elements, filename, bgColor, circleColor)
	if err != nil {
		t.Fatalf("failed to render and save SVG: %s", err)
	}

	// Check if file exists
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		t.Fatalf("file %s was not created", filename)
	}

	// Clean up after test
	os.Remove(filename)
}
