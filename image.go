package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	colorNone    = ""
	colorReset   = "\x1b[0m"
	colorBlack   = "\x1b[30m"
	colorRed     = "\x1b[31m"
	colorGreen   = "\x1b[32m"
	colorYellow  = "\x1b[33m"
	colorBlue    = "\x1b[34m"
	colorMagenta = "\x1b[35m"
	colorCyan    = "\x1b[36m"
	colorWhite   = "\x1b[37m"
)

var (
	colorRGBABlack   = color.RGBA{0, 0, 0, 255}
	colorRGBARed     = color.RGBA{255, 0, 0, 255}
	colorRGBAGreen   = color.RGBA{0, 255, 0, 255}
	colorRGBAYellow  = color.RGBA{255, 255, 0, 255}
	colorRGBABlue    = color.RGBA{0, 0, 255, 255}
	colorRGBAMagenta = color.RGBA{255, 0, 255, 255}
	colorRGBACyan    = color.RGBA{0, 255, 255, 255}
	colorRGBAWhite   = color.RGBA{255, 255, 255, 255}

	colorMap = map[string]color.RGBA{
		colorNone:    colorRGBAWhite,
		colorBlack:   colorRGBABlack,
		colorRed:     colorRGBARed,
		colorGreen:   colorRGBAGreen,
		colorYellow:  colorRGBAYellow,
		colorBlue:    colorRGBABlue,
		colorMagenta: colorRGBAMagenta,
		colorCyan:    colorRGBACyan,
		colorWhite:   colorRGBAWhite,
	}

	colorStringMap = map[string]color.RGBA{
		"black":   colorRGBABlack,
		"red":     colorRGBARed,
		"green":   colorRGBAGreen,
		"yellow":  colorRGBAYellow,
		"blue":    colorRGBABlue,
		"magenta": colorRGBAMagenta,
		"cyan":    colorRGBACyan,
		"white":   colorRGBAWhite,
	}
)

func writeImage(texts []string, w io.Writer, appconf applicationConfig) {
	var (
		charWidth   = appconf.fontsize / 2
		charHeight  = appconf.fontsize
		imageWidth  = maxStringWidth(texts) * charWidth
		imageHeight = len(texts) * charHeight
		img         = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
		face        = readFace(appconf.fontfile, float64(appconf.fontsize))
	)

	drawBackground(img, appconf.background)

	posY := charHeight
	for _, line := range texts {
		posX := 0
		// 色コード以外のエスケープコードを削除
		line = removeNotColorEscapeSequences(line)
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			col, matched, suffix := parseText(line)
			drawLabel(img, posX, posY, matched, colorMap[col], face)
			// 処理されなかった残りで次の処理対象を上書き
			// 空になるまでループ
			line = suffix
			posX += runewidth.StringWidth(matched) * charWidth
		}
		posY += charHeight
	}

	if err := png.Encode(w, img); err != nil {
		panic(err)
	}
}

func readFace(fontPath string, fontSize float64) font.Face {
	fontData, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	ft, err := truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	opt := truetype.Options{
		Size:              fontSize,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face := truetype.NewFace(ft, &opt)
	return face
}

func drawBackground(img *image.RGBA, bg color.RGBA) {
	var (
		bounds = img.Bounds().Max
		width  = bounds.X
		height = bounds.Y
	)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, bg)
		}
	}
}

func drawLabel(img *image.RGBA, x, y int, label string, col color.RGBA, face font.Face) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}
