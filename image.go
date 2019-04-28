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

var colorMap = map[string]color.RGBA{
	colorBlack:   color.RGBA{0, 0, 0, 255},
	colorRed:     color.RGBA{255, 0, 0, 255},
	colorGreen:   color.RGBA{0, 255, 0, 255},
	colorYellow:  color.RGBA{255, 255, 0, 255},
	colorBlue:    color.RGBA{0, 0, 255, 255},
	colorMagenta: color.RGBA{255, 0, 255, 255},
	colorCyan:    color.RGBA{0, 255, 255, 255},
	colorWhite:   color.RGBA{255, 255, 255, 255},
}

var face font.Face

func init() {
	// 日本語が使えるフォントのデフォルトとして指定
	fontData, err := ioutil.ReadFile("/usr/share/fonts/truetype/vlgothic/VL-Gothic-Regular.ttf")
	if err != nil {
		panic(err)
	}

	// ft, err := truetype.Parse(gobold.TTF)
	ft, err := truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	opt := truetype.Options{
		Size:              64,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face = truetype.NewFace(ft, &opt)
}

func writeImage(texts []string, w io.Writer) {
	const (
		charWidth  = 32
		charHeight = 64
	)
	var (
		imageWidth  = maxStringWidth(texts) * charWidth
		imageHeight = len(texts) * charHeight
		img         = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	)

	// TODO 入力におうじて変更する
	drawBackground(img, color.RGBA{R: 0, G: 0, B: 0, A: 255})

	posY := charHeight
	for _, line := range texts {
		posX := 0
		// 色コード以外のエスケープコードを削除
		line = removeNotColorEscapeSequences(line)
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			col, matched, suffix := parseText(line)
			drawLabel(img, posX, posY, matched, colorMap[col])
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

func drawLabel(img *image.RGBA, x, y int, label string, col color.RGBA) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		// Face: basicfont.Face7x13,
		Dot: point,
	}
	d.DrawString(label)
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
