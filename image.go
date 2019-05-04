package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// writeImage はテキストのEscapeSequenceから色情報などを読み取り、
// wに書き込む。
func writeImage(w io.Writer, texts []string, appconf applicationConfig) {
	var (
		charWidth   = appconf.fontsize / 2
		charHeight  = int(float64(appconf.fontsize) * 1.1)
		imageWidth  = maxStringWidth(texts) * charWidth
		imageHeight = len(texts) * charHeight
		img         = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
		face        = readFace(appconf.fontfile, float64(appconf.fontsize))
	)

	drawBackgroundAll(img, appconf.background)

	posY := charHeight
	for _, line := range texts {
		posX := 0
		fgCol := colorRGBAWhite
		bgCol := appconf.background
		// 色コード以外のエスケープコードを削除
		line = removeNotColorEscapeSequences(line)
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			col, matched, suffix := parseText(line)
			if strings.HasPrefix(col, "\x1b[38") || strings.HasPrefix(col, "\x1b[48") {
				spl := strings.Split(col, ";")
				// TODO need refactoring
				switch len(spl) {
				case 3:
					// 256 color code
					rep := strings.Replace(spl[2], "m", "", -1)
					colorCode, err := strconv.Atoi(rep)
					if err != nil {
						panic(err)
					}
					c := terminal256ColorMap[colorCode]
					switch spl[0] {
					case "\x1b[38":
						fgCol = c
					case "\x1b[48":
						bgCol = c
					}
				case 5:
					// RGB color code
					// TODO too dirty
					var r, g, b uint64
					var err error
					r, err = strconv.ParseUint(spl[2], 0, 8)
					if err != nil {
						panic(err)
					}
					g, err = strconv.ParseUint(spl[3], 0, 8)
					if err != nil {
						panic(err)
					}
					b, err = strconv.ParseUint(strings.Replace(spl[4], "m", "", -1), 0, 8)
					if err != nil {
						panic(err)
					}
					c := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
					switch spl[0] {
					case "\x1b[38":
						fgCol = c
					case "\x1b[48":
						bgCol = c
					}
				}
			} else if strings.HasPrefix(col, "\x1b[4") {
				// 色が背景色指定の場合
				bgCol = colorEscapeSequenceMap[col]
			} else if strings.HasPrefix(col, "\x1b[3") {
				fgCol = colorEscapeSequenceMap[col]
			}
			switch col {
			case colorEscapeSequenceReset, colorEscapeSequenceResetShort:
				bgCol = appconf.background
				fgCol = colorRGBAWhite
			}
			drawBackground(img, posX, posY-charHeight, matched, bgCol, charWidth, charHeight)
			// テキストが微妙に見切れるので調整
			drawLabel(img, posX, posY-(charHeight/5), matched, fgCol, face)
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

// readFace はfontPathのフォントファイルからfaceを返す。
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

// drawBackgroundAll はimgにbgを背景色として描画する。
func drawBackgroundAll(img *image.RGBA, bg color.RGBA) {
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

// drawLabel はimgにラベルを描画する。
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

func drawBackground(img *image.RGBA, posX, posY int, label string, col color.RGBA, charWidth, charHeight int) {
	var (
		tw     = runewidth.StringWidth(label)
		width  = tw * charWidth
		height = charHeight
	)
	for x := posX; x < posX+width; x++ {
		for y := posY; y < posY+height; y++ {
			img.Set(x, y, col)
		}
	}
}

// maxStringWidth は表示上のテキストの最も幅の長い長さを返却する。
func maxStringWidth(s []string) (max int) {
	for _, v := range s {
		text := classifyString(v).onlyText()
		width := runewidth.StringWidth(text)
		if max < width {
			max = width
		}
	}
	return
}
