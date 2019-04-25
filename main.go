package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
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

func main() {
	// 標準入力から文字列を取得
	inputStr := readStdin()[0]
	fmt.Println(inputStr)
	outFile := os.Args[1]

	// 色コード以外のエスケープコードを削除
	inputStr = removeNotColorEscapeSequences(inputStr)

	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	for inputStr != "" {
		// 色文字列の句切れごとに画像に色指定して書き込む
		col, matched, suffix := parseText(inputStr)
		addLabel(img, 20, 30, matched, rgbMap[col])
		inputStr = suffix
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

func readStdin() (ret []string) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		ret = append(ret, line)
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	return
}

func addLabel(img *image.RGBA, x, y int, label string, col color.RGBA) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
