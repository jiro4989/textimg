package image

import (
	"bytes"
	"fmt"
	"image"
	c "image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/token"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font/basicfont"
)

const (
	// テストで使うフォントサイズ。
	// NewImage は charWidth = fontSize / 2、charHeight = fontSize * 1.1 で
	// 1文字あたりの大きさを決める。この値のとき 10x22 になる。
	testFontSize   = 20
	testCharWidth  = testFontSize / 2
	testCharHeight = 22

	// 絵文字画像として使うコードポイント。'A' を使うのは、
	// basicfont がグリフを持つため「フォントで描画されたこと」も検査できるから。
	emojiRune = 'A'
)

// テスト用の既定のパラメータを返す。
// フォントには golang.org/x/image に同梱されているビットマップフォントを使う。
// これによりフォントファイルを用意せずに描画処理をテストできる。
func newTestImageParam() *ImageParam {
	return &ImageParam{
		BaseWidth:       10,
		BaseHeight:      1,
		ForegroundColor: c.RGBA{255, 255, 255, 255},
		BackgroundColor: c.RGBA{0, 0, 0, 255},
		FontSize:        testFontSize,
		FontFace:        basicfont.Face7x13,
	}
}

// 指定した範囲のうち、条件を満たすピクセルの数を返す。
// フォントによって字形が変わるため、ピクセルの位置や数を厳密に検査せず、
// 「その色が存在するかどうか」で判定するために使う。
func countPixelsIn(img image.Image, rect image.Rectangle, match func(r, g, b uint32) bool) int {
	var count int
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if match(r, g, b) {
				count++
			}
		}
	}
	return count
}

func countPixels(img image.Image, match func(r, g, b uint32) bool) int {
	return countPixelsIn(img, img.Bounds(), match)
}

// 条件を満たすピクセルを含む一番右の列を返す。存在しなければ -1。
// 背景矩形の右端から、文字が何セルぶんの幅として扱われたかを検査する。
func rightmostColumn(img image.Image, match func(r, g, b uint32) bool) int {
	b := img.Bounds()
	for x := b.Max.X - 1; x >= b.Min.X; x-- {
		if countPixelsIn(img, image.Rect(x, b.Min.Y, x+1, b.Max.Y), match) > 0 {
			return x
		}
	}
	return -1
}

// RGBA() は各値を 0x0000 から 0xffff の範囲で返すため、中間の値を閾値にする。
func isWhite(r, g, b uint32) bool { return r > 0x8000 && g > 0x8000 && b > 0x8000 }
func isRed(r, g, b uint32) bool   { return r > 0x8000 && g < 0x4000 && b < 0x4000 }

// 絵文字画像を模したPNGファイルを生成する。
// 描画されたかどうかを色で判別するために単色で塗りつぶす。
func createEmojiFile(t *testing.T, dir string, r rune, col c.RGBA) {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, col)
		}
	}

	path := filepath.Join(dir, fmt.Sprintf("emoji_u%.4x.png", r))
	fp, err := os.Create(path)
	if err != nil {
		t.Fatalf("絵文字画像ファイルの作成に失敗した: %v", err)
	}
	defer fp.Close()

	if err := png.Encode(fp, img); err != nil {
		t.Fatalf("絵文字画像ファイルの書き込みに失敗した: %v", err)
	}
}

func TestNewImage(t *testing.T) {
	// 期待値は定数から計算せず、そのまま書き下す。
	// 実装と同じ式で期待値を作ると、式が変わっても検出できないため。
	tests := []struct {
		desc       string
		baseWidth  int
		baseHeight int
		fontSize   int
		wantWidth  int
		wantHeight int
	}{
		{desc: "1文字ぶんの画像", baseWidth: 1, baseHeight: 1, fontSize: 20, wantWidth: 10, wantHeight: 22},
		{desc: "横10文字ぶんの画像", baseWidth: 10, baseHeight: 1, fontSize: 20, wantWidth: 100, wantHeight: 22},
		{desc: "縦3行ぶんの画像", baseWidth: 5, baseHeight: 3, fontSize: 20, wantWidth: 50, wantHeight: 66},
		{desc: "フォントサイズが変わると1文字の大きさも変わる", baseWidth: 2, baseHeight: 1, fontSize: 40, wantWidth: 40, wantHeight: 44},
		{desc: "奇数のフォントサイズは切り捨てられる", baseWidth: 2, baseHeight: 1, fontSize: 15, wantWidth: 14, wantHeight: 16},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			p := newTestImageParam()
			p.BaseWidth = tt.baseWidth
			p.BaseHeight = tt.baseHeight
			p.FontSize = tt.fontSize

			img := NewImage(p)
			b := img.image.Bounds()
			assert.Equal(tt.wantWidth, b.Dx())
			assert.Equal(tt.wantHeight, b.Dy())
		})
	}
}

// アニメーション有効時に AnimationLineCount が 0 だと、
// imageHeight / (BaseHeight / AnimationLineCount) がゼロ除算になる。
// 現在の挙動を固定しておく。
func TestNewImage_AnimationLineCountZero(t *testing.T) {
	p := newTestImageParam()
	p.BaseHeight = 2
	p.UseAnimation = true
	p.AnimationLineCount = 0

	assert.Panics(t, func() { NewImage(p) })
}

// AnimationLineCount が BaseHeight より大きいと BaseHeight/AnimationLineCount が 0 になり、
// これもゼロ除算になる。
func TestNewImage_AnimationLineCountTooLarge(t *testing.T) {
	p := newTestImageParam()
	p.BaseHeight = 2
	p.UseAnimation = true
	p.AnimationLineCount = 3

	assert.Panics(t, func() { NewImage(p) })
}

func TestImage_Draw(t *testing.T) {
	t.Run("文字を描画すると前景色のピクセルが現れる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}})

		assert.NoError(err)
		assert.Positive(countPixels(img.image, isWhite))
	})

	t.Run("文字が無いときは前景色のピクセルが現れない", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{})

		assert.NoError(err)
		assert.Zero(countPixels(img.image, isWhite))
	})

	t.Run("改行すると2行目に描画される", func(t *testing.T) {
		assert := assert.New(t)
		p := newTestImageParam()
		p.BaseHeight = 2
		img := NewImage(p)

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: "a\nb"}})

		assert.NoError(err)
		b := img.image.Bounds()
		firstLine := image.Rect(b.Min.X, b.Min.Y, b.Max.X, testCharHeight)
		secondLine := image.Rect(b.Min.X, testCharHeight, b.Max.X, b.Max.Y)
		assert.Positive(countPixelsIn(img.image, firstLine, isWhite))
		assert.Positive(countPixelsIn(img.image, secondLine, isWhite))
	})

	t.Run("改行しないときは2行目に描画されない", func(t *testing.T) {
		assert := assert.New(t)
		p := newTestImageParam()
		p.BaseHeight = 2
		img := NewImage(p)

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: "ab"}})

		assert.NoError(err)
		b := img.image.Bounds()
		secondLine := image.Rect(b.Min.X, testCharHeight, b.Max.X, b.Max.Y)
		assert.Zero(countPixelsIn(img.image, secondLine, isWhite))
	})

	t.Run("背景色を変更すると文字の背景が塗られる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{
			{Kind: token.KindColor, ColorType: token.ColorTypeBackground, Color: color.RGBA{R: 255, A: 255}},
			{Kind: token.KindText, Text: "a"},
		})

		assert.NoError(err)
		assert.Positive(countPixels(img.image, isRed))
	})

	t.Run("前景色を変更すると文字が変更後の色で描画される", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{
			{Kind: token.KindColor, ColorType: token.ColorTypeForeground, Color: color.RGBA{R: 255, A: 255}},
			{Kind: token.KindText, Text: "a"},
		})

		assert.NoError(err)
		assert.Positive(countPixels(img.image, isRed))
		assert.Zero(countPixels(img.image, isWhite))
	})

	t.Run("前景色だけをリセットできる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{
			{Kind: token.KindColor, ColorType: token.ColorTypeForeground, Color: color.RGBA{R: 255, A: 255}},
			{Kind: token.KindText, Text: "a"},
			token.NewResetForegroundColor(),
			{Kind: token.KindText, Text: "b"},
		})

		assert.NoError(err)
		b := img.image.Bounds()
		firstChar := image.Rect(b.Min.X, b.Min.Y, testCharWidth, b.Max.Y)
		rest := image.Rect(testCharWidth, b.Min.Y, b.Max.X, b.Max.Y)
		assert.Positive(countPixelsIn(img.image, firstChar, isRed))
		assert.Zero(countPixelsIn(img.image, firstChar, isWhite))
		assert.Zero(countPixelsIn(img.image, rest, isRed))
		assert.Positive(countPixelsIn(img.image, rest, isWhite))
	})

	t.Run("背景色だけをリセットできる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{
			{Kind: token.KindColor, ColorType: token.ColorTypeBackground, Color: color.RGBA{R: 255, A: 255}},
			{Kind: token.KindText, Text: "a"},
			token.NewResetBackgroundColor(),
			{Kind: token.KindText, Text: "b"},
		})

		assert.NoError(err)
		b := img.image.Bounds()
		firstChar := image.Rect(b.Min.X, b.Min.Y, testCharWidth, b.Max.Y)
		rest := image.Rect(testCharWidth, b.Min.Y, b.Max.X, b.Max.Y)
		assert.Positive(countPixelsIn(img.image, firstChar, isRed))
		assert.Zero(countPixelsIn(img.image, rest, isRed))
	})

	t.Run("色をリセットすると前景色と背景色の両方が既定に戻る", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{
			{Kind: token.KindColor, ColorType: token.ColorTypeForeground, Color: color.RGBA{G: 255, A: 255}},
			{Kind: token.KindColor, ColorType: token.ColorTypeBackground, Color: color.RGBA{R: 255, A: 255}},
			{Kind: token.KindText, Text: "a"},
			token.NewResetColor(),
			{Kind: token.KindText, Text: "b"},
		})

		assert.NoError(err)
		b := img.image.Bounds()
		firstChar := image.Rect(b.Min.X, b.Min.Y, testCharWidth, b.Max.Y)
		rest := image.Rect(testCharWidth, b.Min.Y, b.Max.X, b.Max.Y)
		// 1文字目は変更後の背景色(赤)、2文字目は既定の前景色(白)に戻る
		assert.Positive(countPixelsIn(img.image, firstChar, isRed))
		assert.Zero(countPixelsIn(img.image, rest, isRed))
		assert.Positive(countPixelsIn(img.image, rest, isWhite))
	})

	t.Run("色を反転すると前景色と背景色が入れ替わる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())

		err := img.Draw(token.Tokens{token.NewReverseColor(), {Kind: token.KindText, Text: "a"}})

		assert.NoError(err)
		b := img.image.Bounds()
		firstChar := image.Rect(b.Min.X, b.Min.Y, testCharWidth, b.Max.Y)
		// 反転すると背景が前景色(白)で1文字ぶん塗りつぶされる。
		// 字形だけが白い通常の描画と違い、セルのほぼ全体が白くなる。
		white := countPixelsIn(img.image, firstChar, isWhite)
		assert.Greater(white, testCharWidth*testCharHeight/2)
	})
}

// 文字幅は runewidth に従い、半角は1セル、全角は2セルぶんの背景が塗られる。
func TestImage_Draw_RuneWidth(t *testing.T) {
	tests := []struct {
		desc      string
		text      string
		wantCells int
	}{
		{desc: "半角は1セル", text: "a", wantCells: 1},
		{desc: "全角は2セル", text: "あ", wantCells: 2},
		{desc: "半角2文字は2セル", text: "ab", wantCells: 2},
		// image パッケージの init() は StrictEmojiNeutral を false にしているが、
		// この設定は EastAsianWidth が有効なときにしか効かない。既定の条件下では
		// Unicode Neutral の絵文字は幅1のままになる。
		{desc: "既定ではUnicode Neutralの絵文字は1セル", text: "\U0001F441", wantCells: 1}, // 👁
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			img := NewImage(newTestImageParam())

			err := img.Draw(token.Tokens{
				{Kind: token.KindColor, ColorType: token.ColorTypeBackground, Color: color.RGBA{R: 255, A: 255}},
				{Kind: token.KindText, Text: tt.text},
			})

			assert.NoError(err)
			assert.Equal(tt.wantCells*testCharWidth-1, rightmostColumn(img.image, isRed))
		})
	}
}

// EastAsianWidth を有効にすると、init() が設定した StrictEmojiNeutral=false が効き、
// Unicode Neutral の絵文字が幅2になる。init() の意図はこの条件下でのみ達成される。
func TestImage_Draw_RuneWidth_EastAsianWidth(t *testing.T) {
	assert := assert.New(t)

	orig := runewidth.DefaultCondition.EastAsianWidth
	runewidth.DefaultCondition.EastAsianWidth = true
	t.Cleanup(func() { runewidth.DefaultCondition.EastAsianWidth = orig })

	img := NewImage(newTestImageParam())
	err := img.Draw(token.Tokens{
		{Kind: token.KindColor, ColorType: token.ColorTypeBackground, Color: color.RGBA{R: 255, A: 255}},
		{Kind: token.KindText, Text: "\U0001F441"}, // 👁
	})

	assert.NoError(err)
	assert.Equal(2*testCharWidth-1, rightmostColumn(img.image, isRed))
}

func TestImage_Draw_Emoji(t *testing.T) {
	t.Run("絵文字画像が存在するときは画像を描画する", func(t *testing.T) {
		assert := assert.New(t)
		dir := t.TempDir()
		createEmojiFile(t, dir, emojiRune, c.RGBA{255, 0, 0, 255})

		p := newTestImageParam()
		p.EmojiDir = dir
		img := NewImage(p)

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: string(rune(emojiRune))}})

		assert.NoError(err)
		// 画像が1文字目の位置に描かれ、フォントのグリフは描かれていない
		assert.Positive(countPixels(img.image, isRed))
		assert.Zero(countPixels(img.image, isWhite))
		b := img.image.Bounds()
		assert.Less(rightmostColumn(img.image, isRed), b.Max.X/2)
	})

	t.Run("絵文字フォントを使うときは画像ではなくフォントで描画する", func(t *testing.T) {
		assert := assert.New(t)
		dir := t.TempDir()
		createEmojiFile(t, dir, emojiRune, c.RGBA{255, 0, 0, 255})

		p := newTestImageParam()
		p.EmojiDir = dir
		p.UseEmoji = true
		p.EmojiFontFace = basicfont.Face7x13
		img := NewImage(p)

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: string(rune(emojiRune))}})

		assert.NoError(err)
		// 画像は使われず、絵文字フォントのグリフが描かれる
		assert.Zero(countPixels(img.image, isRed))
		assert.Positive(countPixels(img.image, isWhite))
	})

	t.Run("絵文字画像が存在しないときはフォントで描画する", func(t *testing.T) {
		assert := assert.New(t)
		p := newTestImageParam()
		p.EmojiDir = t.TempDir()
		img := NewImage(p)

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: string(rune(emojiRune))}})

		assert.NoError(err)
		assert.Zero(countPixels(img.image, isRed))
		assert.Positive(countPixels(img.image, isWhite))
	})

	t.Run("異常系: 絵文字画像が壊れているときはエラーを返す", func(t *testing.T) {
		assert := assert.New(t)
		dir := t.TempDir()
		path := filepath.Join(dir, fmt.Sprintf("emoji_u%.4x.png", emojiRune))
		if err := os.WriteFile(path, []byte("これはPNGではない"), 0600); err != nil {
			t.Fatalf("ファイルの作成に失敗した: %v", err)
		}

		p := newTestImageParam()
		p.EmojiDir = dir
		img := NewImage(p)

		err := img.Draw(token.Tokens{{Kind: token.KindText, Text: string(rune(emojiRune))}})

		assert.Error(err)
	})
}

func TestImage_Draw_Resize(t *testing.T) {
	assert := assert.New(t)
	p := newTestImageParam()
	p.ResizeWidth = 50
	p.ResizeHeight = 11
	img := NewImage(p)

	err := img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}})

	assert.NoError(err)
	b := img.image.Bounds()
	assert.Equal(50, b.Dx())
	assert.Equal(11, b.Dy())
	// 縮小しても描画内容が残っていること。空画像で通らないようにする。
	assert.Positive(countPixels(img.image, isWhite))
}

func TestImage_Draw_Animation(t *testing.T) {
	assert := assert.New(t)
	p := newTestImageParam()
	p.BaseHeight = 2
	p.UseAnimation = true
	p.AnimationLineCount = 1
	p.Delay = 20
	img := NewImage(p)

	// 1行目だけに描画し、フレームの切り出し位置を検査できるようにする
	err := img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}})

	assert.NoError(err)
	// 2行の画像を1行ずつ切り出すため、2枚のフレームが生成される
	assert.Len(img.animationImages, 2)
	for _, frame := range img.animationImages {
		assert.Equal(testCharHeight, frame.Bounds().Dy())
	}
	// 1枚目には描画内容があり、2枚目は空
	assert.Positive(countPixels(img.animationImages[0], isWhite))
	assert.Zero(countPixels(img.animationImages[1], isWhite))

	var buf bytes.Buffer
	assert.NoError(img.Encode(&buf, ".gif"))
	g, err := gif.DecodeAll(&buf)
	assert.NoError(err)
	assert.Len(g.Image, 2)
	assert.Equal([]int{20, 20}, g.Delay)
}

func TestImage_Draw_AnimationWithResize(t *testing.T) {
	assert := assert.New(t)
	p := newTestImageParam()
	p.BaseHeight = 2
	p.UseAnimation = true
	p.AnimationLineCount = 1
	p.ResizeWidth = 50
	p.ResizeHeight = 11
	img := NewImage(p)

	err := img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}})

	assert.NoError(err)
	// 各フレームもリサイズ後の大きさになる
	assert.Len(img.animationImages, 2)
	for _, frame := range img.animationImages {
		assert.Equal(50, frame.Bounds().Dx())
		assert.Equal(11, frame.Bounds().Dy())
	}
	assert.Positive(countPixels(img.animationImages[0], isWhite))
}

func TestImage_Encode(t *testing.T) {
	t.Run("正常系: pngはデコードでき、大きさと内容が保たれる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())
		assert.NoError(img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}}))

		var buf bytes.Buffer
		assert.NoError(img.Encode(&buf, ".png"))

		decoded, err := png.Decode(&buf)
		assert.NoError(err)
		assert.Equal(100, decoded.Bounds().Dx())
		assert.Equal(testCharHeight, decoded.Bounds().Dy())
		assert.Positive(countPixels(decoded, isWhite))
	})

	t.Run("正常系: gifはデコードでき、1フレームになる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())
		assert.NoError(img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}}))

		var buf bytes.Buffer
		assert.NoError(img.Encode(&buf, ".gif"))

		decoded, err := gif.DecodeAll(&buf)
		assert.NoError(err)
		assert.Len(decoded.Image, 1)
	})

	t.Run("正常系: jpgはデコードでき、大きさが保たれる", func(t *testing.T) {
		assert := assert.New(t)
		img := NewImage(newTestImageParam())
		assert.NoError(img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}}))

		var buf bytes.Buffer
		assert.NoError(img.Encode(&buf, ".jpg"))

		// jpegは非可逆なので色は検査せず、デコードできることと大きさだけを見る
		decoded, err := jpeg.Decode(&buf)
		assert.NoError(err)
		assert.Equal(100, decoded.Bounds().Dx())
		assert.Equal(testCharHeight, decoded.Bounds().Dy())
	})

	tests := []struct {
		desc    string
		ext     string
		wantErr bool
	}{
		{desc: "正常系: jpeg", ext: ".jpeg"},
		{desc: "異常系: 対応していない拡張子", ext: ".bmp", wantErr: true},
		{desc: "異常系: 拡張子が空", ext: "", wantErr: true},
		{desc: "異常系: ドットが無い", ext: "png", wantErr: true},
		{desc: "異常系: 大文字", ext: ".PNG", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			img := NewImage(newTestImageParam())
			assert.NoError(img.Draw(token.Tokens{{Kind: token.KindText, Text: "hello"}}))

			var buf bytes.Buffer
			err := img.Encode(&buf, tt.ext)

			if tt.wantErr {
				assert.Error(err)
				assert.Zero(buf.Len())
				return
			}
			assert.NoError(err)
			assert.Positive(buf.Len())
		})
	}
}
