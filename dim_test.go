package surrender

import (
	"testing"
)

func TestCircleDimensions(t *testing.T) {
	width, height, err := GetSVGDimensions("testdata/circle.svg")
	if err != nil {
		t.Errorf("Failed to get SVG dimensions: %v", err)
	}
	expectedWidth := 100
	expectedHeight := 100
	if width != expectedWidth {
		t.Errorf("Width mismatch. Expected: %d, Got: %d", expectedWidth, width)
	}
	if height != expectedHeight {
		t.Errorf("Height mismatch. Expected: %d, Got: %d", expectedHeight, height)
	}
}

func TestXyprotoDimensions(t *testing.T) {
	width, height, err := GetSVGDimensions("testdata/xyproto.svg")
	if err != nil {
		t.Errorf("Failed to get SVG dimensions: %v", err)
	}
	expectedWidth := 400
	expectedHeight := 400
	if width != expectedWidth {
		t.Errorf("Width mismatch. Expected: %d, Got: %d", expectedWidth, width)
	}
	if height != expectedHeight {
		t.Errorf("Height mismatch. Expected: %d, Got: %d", expectedHeight, height)
	}
}
