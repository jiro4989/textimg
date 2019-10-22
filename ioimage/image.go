package ioimage

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/jiro4989/textimg/escseq"
	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

type (
	WriteConfig struct {
		Foreground    escseq.RGBA // 文字色
		Background    escseq.RGBA // 背景色
		FontFace      font.Face   // フォントファイル
		EmojiFontFace font.Face   // 絵文字用のフォントファイル
		EmojiDir      string      // 絵文字画像ファイルの存在するディレクトリ
		UseEmojiFont  bool        // 絵文字TTFを使う
		FontSize      int         // フォントサイズ
		UseAnimation  bool        // アニメーションGIFを生成する
		Delay         int         // アニメーションのディレイ時間
		LineCount     int         // 入力データのうち何行を1フレーム画像に使うか
	}
)

var (
	// 絵文字描画の際に、普通に描画してほしいけれど絵文字としても定義されている
	// 文字のコードポイント
	exRunes = []rune{
		0x0023, // #
		0x002A, // *
		0x0030, // 0
		0x0031, // 1
		0x0032, // 2
		0x0033, // 3
		0x0034, // 4
		0x0035, // 5
		0x0036, // 6
		0x0037, // 7
		0x0038, // 8
		0x0039, // 9
		0x00A9, // ©
		0x00AE, // ®️
	}
)

// Write はテキストのEscapeSequenceから色情報などを読み取り、
// wに書き込む。
func Write(w io.Writer, imgExt string, texts []string, conf WriteConfig) error {
	var (
		charWidth   = conf.FontSize / 2
		charHeight  = int(float64(conf.FontSize) * 1.1)
		imageWidth  = escseq.StringWidth(texts) * charWidth
		imageHeight = len(texts) * charHeight
	)

	if conf.UseAnimation {
		imageHeight /= (len(texts) / conf.LineCount)
	}

	var (
		img       = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
		face      = conf.FontFace
		emojiFace = conf.EmojiFontFace
		imgs      []*image.RGBA
		delays    []int
	)

	drawBackgroundAll(img, conf.Background)

	posY := charHeight
	for i, line := range texts {
		posX := 0
		fgCol := conf.Foreground
		bgCol := conf.Background
		for line != "" {
			// 色文字列の句切れごとに画像に色指定して書き込む
			k, prefix, suffix := escseq.Prefix(line)
			switch k {
			case escseq.KindEmpty:
				err := errors.New("input string is empty")
				return err
			case escseq.KindText:
				text := prefix
				drawBackground(img, posX, posY-charHeight, text, bgCol, charWidth, charHeight)
				// drawLabel(img, posX, posY-(charHeight/5), text, fgCol, face)
				for _, r := range []rune(text) {
					drawText(img, posX, posY-(charHeight/5), r, fgCol, bgCol, face, emojiFace, conf.EmojiDir, conf.UseEmojiFont)
					posX += runewidth.RuneWidth(r) * charWidth
				}
			case escseq.KindColor:
				colors := escseq.ParseColor(prefix)
				for _, v := range colors {
					switch v.ColorType {
					case escseq.ColorTypeReset:
						fgCol = conf.Foreground
						bgCol = conf.Background
					case escseq.ColorTypeReverse:
						fgCol, bgCol = bgCol, fgCol
					case escseq.ColorTypeForeground:
						fgCol = v.Color
					case escseq.ColorTypeBackground:
						bgCol = v.Color
					default:
						// 未実装のcolorTypeでは何もしない
					}
				}
			case escseq.KindNotColor:
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

		if conf.UseAnimation {
			if (i+1)%conf.LineCount == 0 {
				posY = charHeight
				imgs = append(imgs, img)
				delays = append(delays, conf.Delay)
				img = image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
				drawBackgroundAll(img, conf.Background)
			}
		}
	}

	var err error
	switch imgExt {
	case ".png":
		err = png.Encode(w, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(w, img, nil)
	case ".gif":
		if conf.UseAnimation {
			err = gif.EncodeAll(w, &gif.GIF{
				Image: toPalettes(imgs),
				Delay: delays,
			})
		} else {
			err = gif.Encode(w, img, nil)
		}
	default:
		// root.goで拡張子の判定をしているため、このブロックには到達しないはず
		err = errors.New(fmt.Sprintf("%s is not supported extension.", imgExt))
	}
	if err != nil {
		return err
	}
	return nil
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

// r が例外的なコードポイントに存在するかを判定する。
// http://unicode.org/Public/emoji/4.0/emoji-data.txt
//
// ここでtrueを返す文字は、絵文字データ的には絵文字ではあるものの、
// シェル芸bot環境ではテキストとして表示したいので例外的に除外するために指定して
// いる。
func isExceptionallyCodePoint(r rune) bool {
	for _, ex := range exRunes {
		if r == ex {
			return true
		}
	}

	return false
}
