package surrender

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	t.Run("test parsing empty SVG", func(t *testing.T) {
		_, err := ParseFile("testdata/empty.svg")
		assert.NoError(t, err)
	})

	t.Run("test parsing SVG with single circle", func(t *testing.T) {
		elements, err := ParseFile("testdata/circle.svg")
		assert.NoError(t, err)
		assert.Len(t, elements, 1)

		circle, ok := elements[0].(SvgCircle)
		assert.True(t, ok)
		assert.Equal(t, circle.Cx, 50)
		assert.Equal(t, circle.Cy, 50)
		assert.Equal(t, circle.R, 40)
		assert.Equal(t, circle.Fill, color.RGBA{0, 0, 0, 255}) // assuming black color
	})

}
