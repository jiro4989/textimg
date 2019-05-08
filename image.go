package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

type encodeFormat int

const (
	encodeFormatPNG encodeFormat = iota
	encodeFormatJPG
	encodeFormatGIF
)

// writeImage はテキストのEscapeSequenceから色情報などを読み取り、
// wに書き込む。
func writeImage(w io.Writer, encFmt encodeFormat, texts []string, appconf applicationConfig) {
	var (
		charWidth   = appconf.fontsize / 2
		charHeight  = int(float64(appconf.fontsize) * 1.1)
		imageWidth  = maxStringWidth(texts) * charWidth
		imageHeight = len(texts) * charHeight
	)

	if appconf.useAnimation {
		imageHeight /= (len(texts) / appconf.lineCount)
	}

	var (
		img    = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
		face   = readFace(appconf.fontfile, float64(appconf.fontsize))
		imgs   []*image.RGBA
		delays []int
	)

	drawBackgroundAll(img, appconf.background)

	posY := charHeight
	for i, line := range texts {
		posX := 0
		fgCol := appconf.foreground
		bgCol := appconf.background
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			k, prefix, suffix := getPrefix(line)
			switch k {
			case kindEmpty:
				fmt.Fprintln(os.Stderr, "[WARN] input string is empty")
				return
			case kindText:
				text := prefix
				drawBackground(img, posX, posY-charHeight, text, bgCol, charWidth, charHeight)
				// テキストが微妙に見切れるので調整
				drawLabel(img, posX, posY-(charHeight/5), text, fgCol, face)
				posX += runewidth.StringWidth(prefix) * charWidth
			case kindEscapeSequenceColor:
				colors := parseColorEscapeSequence(prefix)
				for _, v := range colors {
					switch v.colorType {
					case colorTypeReset:
						fgCol = appconf.foreground
						bgCol = appconf.background
					case colorTypeReverse:
						fgCol, bgCol = bgCol, fgCol
					case colorTypeForeground:
						fgCol = v.color
					case colorTypeBackground:
						bgCol = v.color
					default:
						// 未実装のcolorTypeでは何もしない
					}
				}
			case kindEscapeSequenceNotColor:
				// 色出力と関係のないエスケープシーケンスの場合は何もしない
			default:
				// 到達しないはず
				panic(fmt.Sprintf("Illegal kind: %v", k))
			}
			// 処理されなかった残りで次の処理対象を上書き
			// 空になるまでループ
			line = suffix
		}
		posY += charHeight

		if appconf.useAnimation {
			if (i+1)%appconf.lineCount == 0 {
				posY = charHeight
				imgs = append(imgs, img)
				delays = append(delays, appconf.delay)
				img = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
				drawBackgroundAll(img, appconf.background)
			}
		}
	}

	var err error
	switch encFmt {
	case encodeFormatPNG:
		err = png.Encode(w, img)
	case encodeFormatJPG:
		err = jpeg.Encode(w, img, nil)
	case encodeFormatGIF:
		if appconf.useAnimation {
			err = gif.EncodeAll(w, &gif.GIF{
				Image: toPalettes(imgs),
				Delay: delays,
			})
		} else {
			err = gif.Encode(w, img, nil)
		}
	default:
		err = errors.New(fmt.Sprintf("%v is not supported.", encFmt))
	}
	if err != nil {
		panic(err)
	}
}

// getEncodeFormat は画像ファイルの拡張子からエンコードフォーマットを取得する。
// 空文字列を指定した場合はPNGを返す。
// 他はPNG, JPG, GIFのみをサポートする。
// それ以外の拡張子のパスが渡された場合はエラーを返す。
func getEncodeFormat(path string) (encodeFormat, error) {
	if path == "" {
		return encodeFormatPNG, nil
	}

	ext := filepath.Ext(strings.ToLower(path))
	switch ext {
	case ".png":
		return encodeFormatPNG, nil
	case ".jpg", ".jpeg":
		return encodeFormatJPG, nil
	case ".gif":
		return encodeFormatGIF, nil
	}

	return -1, errors.New(fmt.Sprintf("[WARN] %s is not supported", ext))
}

// readFace はfontPathのフォントファイルからfaceを返す。
func readFace(fontPath string, fontSize float64) font.Face {
	var fontData []byte

	// ファイルが存在しなければビルトインのフォントをデフォルトとして使う
	_, err := os.Stat(fontPath)
	if err == nil {
		fontData, err = ioutil.ReadFile(fontPath)
		if err != nil {
			panic(err)
		}
	} else {
		msg := fmt.Sprintf("[WARN] %s is not found. please set font path with `-f` option", fontPath)
		fmt.Fprintln(os.Stderr, msg)
		fontData = gomono.TTF
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

func toPalettes(imgs []*image.RGBA) (ret []*image.Paletted) {
	for _, v := range imgs {
		bounds := v.Bounds()
		p := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(p, p.Rect, v, bounds.Min, draw.Over)
		ret = append(ret, p)
	}
	return
}
