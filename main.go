package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
)

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
	drawBackground(img, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	for _, line := range strings.Split(inputStr, "\n") {
		posX := 0
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			col, matched, suffix := parseText(line)
			drawLabel(img, posX, posY, matched, colorMap[col])
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
