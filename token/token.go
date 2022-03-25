package token

import (
	"strconv"
	"strings"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/mattn/go-runewidth"
)

type (
	Kind      int
	ColorType int
	Token     struct {
		Kind      Kind
		ColorType ColorType
		Color     color.RGBA
		Text      string
	}
	Tokens []Token
)

const (
	KindEmpty Kind = iota
	KindText
	KindColor
	KindNotColor

	ColorTypeReset       ColorType = iota // \x1b[0m 指定をリセット
	ColorTypeBold                         // \x1b[1m 太字
	ColorTypeDim                          // \x1b[2m 薄く表示
	ColorTypeItalic                       // \x1b[3m イタリック
	ColorTypeUnderline                    // \x1b[4m アンダーライン
	ColorTypeBlink                        // \x1b[5m ブリンク
	ColorTypeSpeedyBlink                  // \x1b[6m 高速ブリンク
	ColorTypeReverse                      // \x1b[7m 文字色と背景色の反転
	ColorTypeHide                         // \x1b[8m 表示を隠す
	ColorTypeDelete                       // \x1b[9m 取り消し
	ColorTypeForeground
	ColorTypeBackground
)

func NewResetColor() Token {
	return Token{
		Kind:      KindColor,
		ColorType: ColorTypeReset,
	}
}

func NewText(text string) Token {
	return Token{
		Kind: KindText,
		Text: text,
	}
}

func NewStandardColorWithCategory(text string) Token {
	n, _ := strconv.Atoi(text)
	return Token{
		Kind:      KindColor,
		ColorType: colorType(n),
		Color:     color.ANSIMap[n],
	}
}

func NewExtendedColor(text string) Token {
	n, _ := strconv.Atoi(text)
	return Token{
		Kind:      KindColor,
		ColorType: colorType(n),
		Color:     color.RGBA{A: 255},
	}
}

func colorType(n int) ColorType {
	var t ColorType
	switch n / 10 {
	case 3, 9:
		t = ColorTypeForeground
	case 4, 10:
		t = ColorTypeBackground
	}
	return t
}

func (t *Tokens) MaxStringWidth() int {
	var strs []string
	for _, tt := range *t {
		if tt.Kind != KindText {
			continue
		}
		s := tt.Text
		strs = append(strs, s)
	}
	s := strings.Join(strs, "")
	lines := strings.Split(s, "\n")
	var max int
	for _, line := range lines {
		w := runewidth.StringWidth(line)
		if max < w {
			max = w
		}
	}
	return max
}

func (t *Tokens) StringLines() []string {
	var strs []string
	for _, tt := range *t {
		if tt.Kind != KindText {
			continue
		}
		s := tt.Text
		strs = append(strs, s)
	}
	s := strings.Join(strs, "")
	lines := strings.Split(s, "\n")
	return lines
}
