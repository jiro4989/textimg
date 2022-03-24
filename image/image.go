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
		lineCount              int
		animationLineCount     int
		resizeWidth            int
		resizeHeight           int
		delay                  int
	}
	ImageParam struct {
		BaseWidth          int
		BaseHeight         int
		ForegroundColor    c.RGBA // 文字色
		BackgroundColor    c.RGBA // 背景色
		FontSize           int    // フォントサイズ
		FontFace           font.Face
		EmojiFontFace      font.Face
		EmojiDir           string
		UseEmoji           bool
		AnimationLineCount int
		ResizeWidth        int
		ResizeHeight       int
		Delay              int
	}
)

func NewImage(p *ImageParam) *Image {
	var (
		charWidth   = p.FontSize / 2
		charHeight  = int(float64(p.FontSize) * 1.1)
		imageWidth  = p.BaseWidth * charWidth
		imageHeight = p.BaseHeight * charHeight
	)

	if 0 < p.AnimationLineCount {
		imageHeight /= (p.BaseHeight / p.AnimationLineCount)
	}

	image := newImage(imageWidth, imageHeight)

	return &Image{
		image:                  image,
		foregroundColor:        p.ForegroundColor,
		backgroundColor:        p.BackgroundColor,
		defaultForegroundColor: p.ForegroundColor,
		defaultBackgroundColor: p.BackgroundColor,
		fontSize:               p.FontSize,
		fontFace:               p.FontFace,
		emojiFontFace:          p.EmojiFontFace,
		charWidth:              charWidth,
		charHeight:             charHeight,
		emojiDir:               p.EmojiDir,
		useEmoji:               p.UseEmoji,
		animationLineCount:     p.AnimationLineCount,
		resizeWidth:            p.ResizeWidth,
		resizeHeight:           p.ResizeHeight,
	}
}

func newImage(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

func (i *Image) Draw(tokens token.Tokens) error {
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

	i.scale()

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
	// 特殊な位置調整。なんでこんな計算式にしたのか覚えていない
	var (
		x = i.x
		y = i.y + i.charHeight - (i.charHeight / 5)
	)
	// FIXME: なんか警告出てる
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
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
		if 0 < i.animationLineCount && i.lineCount%i.animationLineCount == 0 {
			i.x = 0
			i.y = 0
			i.animationImages = append(i.animationImages, i.newScaledImage())
			b := i.image.Bounds().Max
			i.image = newImage(b.X, b.Y)
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
			i.image.Set(x, y, c.RGBA(i.backgroundColor))
		}
	}
}

func (i *Image) moveRight(r rune) {
	i.x += runewidth.RuneWidth(r) * i.charWidth
}

func (i *Image) newScaledImage() *image.RGBA {
	if i.resizeWidth == 0 && i.resizeHeight == 0 {
		return i.image
	}

	// 呼び出し側で大きさを調整していること
	rect := i.image.Bounds()
	dst := newImage(i.resizeWidth, i.resizeHeight)
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), i.image, rect, draw.Over, nil)
	return dst
}

func (i *Image) scale() {
	i.image = i.newScaledImage()
}
