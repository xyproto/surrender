package surrender

import (
	"image"
	"testing"
)

// TestRender function
func TestRender(t *testing.T) {
	tt := []struct {
		name    string
		file    string
		outfile string
		imgSize image.Point
	}{
		{name: "Render Rainforest", file: "testdata/rainforest_2c_opt.svg", outfile: "testdata/rainforest_2c_opt.png", imgSize: image.Point{X: 302, Y: 240}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the SVG file
			elements, err := ParseFile(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			// Create a new RGBA image
			img := image.NewRGBA(image.Rect(0, 0, tc.imgSize.X, tc.imgSize.Y))

			// Render the SVG elements onto the image
			Render(elements, img)

			// Save the image to a file
			if err := SavePNG(img, tc.outfile); err != nil {
				t.Fatal(err)
			}

			// Output a message to let the tester know what file to check
			t.Logf("Saved rendered image to %s. Please verify it looks correct.", tc.outfile)
		})
	}
}
