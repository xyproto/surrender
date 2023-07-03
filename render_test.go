package surrender

import (
	"image"
	"image/color"
	"testing"
)

func TestRender(t *testing.T) {
	tt := []struct {
		svgFile string // intput file
		pngFile string // output file (if enabled)
		points  []image.Point
		colors  []color.Color
	}{
		{
			svgFile: "testdata/rainforest_2c_opt.svg",
			pngFile: "testdata/rainforest_2c_opt.png",
			points:  []image.Point{image.Point{0, 0}},
			colors:  []color.Color{color.RGBA{255, 255, 238, 255}},
		},
		// Add more test cases here as needed.
	}

	for _, tc := range tt {
		t.Run(tc.svgFile, func(t *testing.T) {
			elements, err := ParseFile(tc.svgFile)
			if err != nil {
				t.Fatal(err)
			}

			width, height, err := GetSVGDimensions(tc.svgFile)
			if err != nil {
				t.Fatal(err)
			}

			img := image.NewRGBA(image.Rect(0, 0, width, height))
			Render(elements, img)

			for i, point := range tc.points {
				if i >= len(tc.colors) {
					t.Fatalf("Not enough expected colors provided for point %v", point)
				}
				expectedColor := tc.colors[i]
				actualColor := img.At(point.X, point.Y)
				if actualColor != expectedColor {
					t.Errorf("At point %v, expected color %v, but got %v", point, expectedColor, actualColor)
				}
			}
			// If all the pixel tests pass, save the image to a PNG file.
			//if err := SavePNG(img, tc.pngFile); err != nil {
			//t.Error(err)
			//}
		})
	}
}
