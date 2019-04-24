package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
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
}

func main() {
	inputStr := readStdin()[0]
	fmt.Println(inputStr)
	outFile := os.Args[1]
	col := getColor(inputStr)

	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	addLabel(img, 20, 30, inputStr, rgbMap[col])

	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func getColor(s string) string {
	if strings.HasPrefix(s, colorBlack) {
		return colorBlack
	}
	if strings.HasPrefix(s, colorRed) {
		return colorRed
	}
	if strings.HasPrefix(s, colorGreen) {
		return colorGreen
	}
	if strings.HasPrefix(s, colorYellow) {
		return colorYellow
	}
	if strings.HasPrefix(s, colorBlue) {
		return colorBlue
	}
	if strings.HasPrefix(s, colorMagenta) {
		return colorMagenta
	}
	if strings.HasPrefix(s, colorSyan) {
		return colorSyan
	}
	return colorWhite
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
