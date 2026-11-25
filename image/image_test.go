package image

import (
	c "image/color"
	"testing"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/token"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func testFontFace(t *testing.T, size float64) font.Face {
	t.Helper()
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Fatal(err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		t.Fatal(err)
	}
	return face
}

func TestDrawBackgroundMultiline(t *testing.T) {
	assert := assert.New(t)

	const fontSize = 20
	face := testFontFace(t, fontSize)
	img := NewImage(&ImageParam{
		BaseWidth:       6,
		BaseHeight:      2,
		ForegroundColor: c.RGBA{255, 255, 255, 255},
		BackgroundColor: c.RGBA{0, 0, 0, 255},
		FontSize:        fontSize,
		FontFace:        face,
	})

	// Red background across a newline: each line must fill only its own text.
	tokens := token.Tokens{
		{Kind: token.KindColor, ColorType: token.ColorTypeBackground, Color: color.RGBARed},
		{Kind: token.KindText, Text: "   \n   "},
		{Kind: token.KindColor, ColorType: token.ColorTypeReset},
	}
	if err := img.Draw(tokens); err != nil {
		t.Fatal(err)
	}

	var (
		charWidth  = fontSize / 2
		charHeight = int(float64(fontSize) * 1.1)
		red        = c.RGBA{255, 0, 0, 255}
		black      = c.RGBA{0, 0, 0, 255}
	)

	// First line: the three spaces are red.
	assert.Equal(red, img.image.At(1, 1), "first line text background")
	// First line: the fill must stop at the end of the line text.
	assert.Equal(black, img.image.At(4*charWidth+1, 1), "first line fill must not exceed text width")
	// Second line: the three spaces are red too.
	assert.Equal(red, img.image.At(1, charHeight+1), "second line text background")
}
