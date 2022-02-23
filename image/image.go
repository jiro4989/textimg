package image

import (
	"image"
	c "image/color"
	"image/draw"
	"os"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/token"
	"github.com/mattn/go-runewidth"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type (
	Image struct {
		image                  *image.RGBA
		animationImages        []*image.RGBA
		x                      int
		y                      int
		foregroundColor        c.RGBA // 文字色
		backgroundColor        c.RGBA // 背景色
		defaultForegroundColor c.RGBA // 文字色
		defaultBackgroundColor c.RGBA // 背景色
		fontSize               int    // フォントサイズ
		fontFace               font.Face
		emojiFontFace          font.Face
		charWidth              int
		charHeight             int
		emojiDir               string
		useEmoji               bool
		useAnimation           bool
		lineCount              int
		animationLineCount     int
	}
	WriteParam struct {
		Foreground    c.RGBA    // 文字色
		Background    c.RGBA    // 背景色
		FontFace      font.Face // フォントファイル
		EmojiFontFace font.Face // 絵文字用のフォントファイル
		EmojiDir      string    // 絵文字画像ファイルの存在するディレクトリ
		UseEmojiFont  bool      // 絵文字TTFを使う
		FontSize      int       // フォントサイズ
		UseAnimation  bool      // アニメーションGIFを生成する
		Delay         int       // アニメーションのディレイ時間
		LineCount     int       // 入力データのうち何行を1フレーム画像に使うか
		ToSlackIcon   bool      // Slackのアイコンサイズにする
		ResizeWidth   int       // 画像の横幅
		ResizeHeight  int       // 画像の縦幅
	}
)

func (i *Image) Write(tokens token.Tokens, conf WriteParam) error {
	i.drawBackgroundAll()

	for _, t := range tokens {
		switch t.Kind {
		case token.KindColor:
			i.updateColor(t.ColorType, t.Color)
		case token.KindText:
			for _, r := range t.Text {
				if err := i.draw(r); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 背景色をデフォルト色で塗りつぶす。
func (i *Image) drawBackgroundAll() {
	var (
		bounds = i.image.Bounds().Max
		width  = bounds.X
		height = bounds.Y
	)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			i.image.Set(x, y, c.RGBA(i.defaultBackgroundColor))
		}
	}
}

func (i *Image) updateColor(t token.ColorType, col color.RGBA) {
	switch t {
	case token.ColorTypeReset:
		i.foregroundColor = i.defaultForegroundColor
		i.backgroundColor = i.defaultBackgroundColor
	case token.ColorTypeReverse:
		i.foregroundColor, i.backgroundColor = i.backgroundColor, i.foregroundColor
	case token.ColorTypeForeground:
		i.foregroundColor = c.RGBA(col)
	case token.ColorTypeBackground:
		i.backgroundColor = c.RGBA(col)
	}
}

func (i *Image) newDrawer(f font.Face) *font.Drawer {
	// FIXME: なんか警告出てる
	point := fixed.Point26_6{fixed.Int26_6(i.x * 64), fixed.Int26_6(i.y * 64)}
	d := &font.Drawer{
		Dst:  i.image,
		Src:  image.NewUniform(c.RGBA(i.foregroundColor)),
		Face: f,
		Dot:  point,
	}
	return d
}

func (i *Image) draw(r rune) error {
	if r == '\n' {
		i.x = 0
		i.y += i.charHeight
		i.lineCount++
		if i.lineCount%i.animationLineCount == 0 {
			i.x = 0
			i.y = 0
			i.animationImages = append(i.animationImages, i.image)
			b := i.image.Bounds().Max
			i.image = image.NewRGBA(image.Rect(0, 0, b.X, b.Y))
			i.drawBackgroundAll()
		}
		return nil
	}

	i.drawBackground(r)
	if ok, emojiPath := isEmoji(r, i.emojiDir); ok {
		if i.useEmoji {
			i.drawRune(r, i.emojiFontFace)
			return nil
		}
		return i.drawEmoji(r, emojiPath)
	}
	i.drawRune(r, i.fontFace)
	return nil
}

// rune文字を画像に書き込む。
// 書き込み終えると座標を更新する。
func (i *Image) drawRune(r rune, f font.Face) {
	d := i.newDrawer(f)
	d.DrawString(string(r))
	i.moveRight(r)
}

func (i *Image) drawEmoji(r rune, path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	emoji, _, err := image.Decode(fp)
	if err != nil {
		return err
	}

	d := i.newDrawer(i.fontFace)
	// 画像サイズをフォントサイズに合わせる
	// 0.9でさらに微妙に調整
	size := int(float64(d.Face.Metrics().Ascent.Floor()+d.Face.Metrics().Descent.Floor()) * 0.9)
	rect := image.Rect(0, 0, size, size)
	dst := image.NewRGBA(rect)
	xdraw.ApproxBiLinear.Scale(dst, rect, emoji, emoji.Bounds(), draw.Over, nil)

	p := image.Pt(d.Dot.X.Floor(), d.Dot.Y.Floor()-d.Face.Metrics().Ascent.Floor())
	draw.Draw(i.image, rect.Add(p), dst, image.Point{}, draw.Over)
	i.moveRight(r)
	return nil
}

func (i *Image) drawBackground(r rune) {
	var (
		tw     = runewidth.RuneWidth(r)
		width  = tw * i.charWidth
		height = i.charHeight
		posX   = i.x
		posY   = i.y
	)
	for x := posX; x < posX+width; x++ {
		for y := posY; y < posY+height; y++ {
			i.image.Set(x, y, c.RGBA(i.foregroundColor))
		}
	}
}

func (i *Image) moveRight(r rune) {
	i.x += runewidth.RuneWidth(r) * i.charWidth
}
