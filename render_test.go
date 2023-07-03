package surrender

import (
	"image"
	"image/color"
	"testing"
)

func TestRender(t *testing.T) {
	tt := []struct {
		name    string
		svgFile string
		points  []image.Point
		colors  []color.Color
	}{
		{
			name:    "testdata/rainforest_2c_opt.png",
			svgFile: "testdata/rainforest_2c_opt.svg",
			points:  []image.Point{image.Point{0, 0}},
			colors:  []color.Color{color.RGBA{255, 255, 238, 255}},
		},
		// Add more test cases here as needed.
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			elements, err := ParseFile(tc.svgFile)
			if err != nil {
				t.Fatal(err)
			}

			img := image.NewRGBA(image.Rect(0, 0, 302, 240))
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
			//if err := SavePNG(img, tc.name+".png"); err != nil {
			//t.Error(err)
			//}
		})
	}
}
