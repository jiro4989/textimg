package ioimage

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/jiro4989/textimg/v3/escseq"
	"github.com/jiro4989/textimg/v3/log"
	"github.com/mattn/go-runewidth"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func init() {
	// Unicode Neutral で定義されている絵文字(例: 👁)を幅2として扱う
	runewidth.DefaultCondition.StrictEmojiNeutral = false
}

func drawText(img *image.RGBA, x, y int, r rune, fgCol, bgCol escseq.RGBA, face, emojiFace font.Face, emojiDir string, useEmoji bool) {
	path := fmt.Sprintf("%s/emoji_u%.4x.png", emojiDir, r)
	_, err := os.Stat(path)
	if err == nil && !isExceptionallyCodePoint(r) {
		// エラーにならないときは絵文字コードポイントにマッチす
		// る画像ファイルが存在するため絵文字として描画
		if useEmoji {
			// EmojiFontを使うときはTTFから絵文字を描画する
			drawLabel(img, x, y, r, fgCol, emojiFace)
			return
		}
		// EmojiFontを使わないときは画像ファイルから絵文字を
		// 描画する
		drawEmoji(img, x, y, r, path, fgCol, face)
		return
	}
	// 絵文字コードポイントにマッチする画像が存在しないときは
	// 普通のテキストとして描画する
	drawLabel(img, x, y, r, fgCol, face)
}

// drawBackgroundAll はimgにbgを背景色として描画する。
func drawBackgroundAll(img *image.RGBA, bg escseq.RGBA) {
	var (
		bounds = img.Bounds().Max
		width  = bounds.X
		height = bounds.Y
	)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA(bg))
		}
	}
}

// drawLabel はimgにラベルを描画する。
func drawLabel(img *image.RGBA, x, y int, r rune, col escseq.RGBA, face font.Face) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA(col)),
		Face: face,
		Dot:  point,
	}
	d.DrawString(string(r))
}

// 絵文字を画像ファイルから読み取って描画する。
func drawEmoji(img *image.RGBA, x, y int, emojiRune rune, path string, col escseq.RGBA, face font.Face) {
	fp, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return
	}
	defer fp.Close()

	emoji, _, err := image.Decode(fp)
	if err != nil {
		log.Error(err)
		return
	}

	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA(col)),
		Face: face,
		Dot:  point,
	}
	// 画像サイズをフォントサイズに合わせる
	// 0.9でさらに微妙に調整
	size := int(float64(d.Face.Metrics().Ascent.Floor()+d.Face.Metrics().Descent.Floor()) * 0.9)
	rect := image.Rect(0, 0, size, size)
	dst := image.NewRGBA(rect)
	xdraw.ApproxBiLinear.Scale(dst, rect, emoji, emoji.Bounds(), draw.Over, nil)

	p := image.Pt(d.Dot.X.Floor(), d.Dot.Y.Floor()-d.Face.Metrics().Ascent.Floor())
	draw.Draw(img, rect.Add(p), dst, image.Point{}, draw.Over)
}

func drawBackground(img *image.RGBA, posX, posY int, label string, col escseq.RGBA, charWidth, charHeight int) {
	var (
		tw     = runewidth.StringWidth(label)
		width  = tw * charWidth
		height = charHeight
	)
	for x := posX; x < posX+width; x++ {
		for y := posY; y < posY+height; y++ {
			img.Set(x, y, color.RGBA(col))
		}
	}
}
