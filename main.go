package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

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
	colorSyan    = "\x1b[36m"
	colorWhite   = "\x1b[37m"
)

var colors = []string{
	colorReset,
	colorBlack,
	colorRed,
	colorGreen,
	colorYellow,
	colorBlue,
	colorMagenta,
	colorSyan,
	colorWhite,
}

var rgbMap = map[string]color.RGBA{
	colorRed:   color.RGBA{255, 0, 0, 255},
	colorGreen: color.RGBA{0, 255, 0, 255},
	colorBlue:  color.RGBA{0, 0, 255, 255},
}

var face font.Face

func init() {
	// 日本語が使えるフォントのデフォルトとして指定
	fontData, err := ioutil.ReadFile("/usr/share/fonts/truetype/vlgothic/VL-Gothic-Regular.ttf")
	// fontData, err := ioutil.ReadFile("/usr/share/fonts/truetype/fonts-japanese-gothic.ttf")
	// fontData, err := ioutil.ReadFile("/usr/share/fonts/truetype/noto/NotoSansJavanese-Regular.ttf")
	// fontData, err := ioutil.ReadFile("/usr/share/fonts/truetype/ubuntu/UbuntuMono-R.ttf")
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

func main() {
	// 標準入力から文字列を取得
	inputStr := strings.Join(readStdin(), "\n")
	fmt.Println(inputStr)
	outFile := os.Args[1]

	// 色コード以外のエスケープコードを削除
	inputStr = removeNotColorEscapeSequences(inputStr)

	const charWidth = 32
	const charHeight = 64

	posY := charHeight
	imageWidth := maxStringWidth(inputStr) * charWidth
	imageHeight := len(strings.Split(inputStr, "\n")) * charHeight
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	for _, line := range strings.Split(inputStr, "\n") {
		posX := 0
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			col, matched, suffix := parseText(line)
			addLabel(img, posX, posY, matched, rgbMap[col])
			line = suffix
			posX += runewidth.StringWidth(matched) * charWidth
		}
		posY += charHeight
	}

	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func addLabel(img *image.RGBA, x, y int, label string, col color.RGBA) {
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

func maxStringWidth(s string) (max int) {
	list := strings.Split(s, "\n")
	for _, v := range list {
		text := classifyString(v).OnlyText()
		width := runewidth.StringWidth(text)
		if max < width {
			max = width
		}
	}
	return
}
