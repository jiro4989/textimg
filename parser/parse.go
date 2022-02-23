package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/mattn/go-runewidth"
)

type (
	Kind      int // エスケープシーケンスの種類
	ColorType int // 文字色か背景色か
	Color     struct {
		Kind      Kind
		ColorType ColorType
		Color     color.RGBA
		Text      string
	}
	Colors []Color
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

var (
	// 出力の仕方や色を制御するエスケープシーケンスにマッチする
	reColorEscapeSequences = regexp.MustCompile(`^\x1b\[[\d;]*m`)
	// 出力の仕方や色を制御するエスケープシーケンス以外のエスケープシーケンスに
	// マッチする
	reNotColorEscapeSequences = regexp.MustCompile(`^\x1b\[\d*[A-HfSTJK]`)

	// \x1b[Nm とかの N に紐づく色種別
	// 例: \x1b[0m
	AttributeMap = map[int]ColorType{
		0: ColorTypeReset,
		1: ColorTypeBold,
		2: ColorTypeDim,
		3: ColorTypeItalic,
		4: ColorTypeUnderline,
		5: ColorTypeBlink,
		6: ColorTypeSpeedyBlink,
		7: ColorTypeReverse,
		8: ColorTypeHide,
		9: ColorTypeDelete,
	}
)

func newResetColor() Color {
	return Color{
		Kind:      KindNotColor,
		ColorType: ColorTypeReset,
	}
}

func newText(text string) Color {
	return Color{
		Kind: KindText,
		Text: text,
	}
}

func newStandardColorWithCategory(text string) Color {
	n, _ := strconv.Atoi(text)
	var t ColorType
	switch n / 10 {
	case 3, 9:
		t = ColorTypeForeground
	case 4, 10:
		t = ColorTypeBackground
	}
	return Color{
		Kind:      KindColor,
		ColorType: t,
		Color:     color.ANSIMap[n],
	}
}

func newExtendedColor256(text string) Color {
	n, _ := strconv.Atoi(text)
	return Color{
		Kind:  KindColor,
		Color: color.Map256[n],
	}
}

func newExtendedColorRGB() Color {
	return Color{
		Kind: KindColor,
	}
}

// ParseColor は色のエスケープシーケンスを解析してRGBAに変換する。
func ParseColor(s string) (colors Colors) {
	s = strings.Replace(s, "\x1b[", "", -1)
	s = strings.Replace(s, "m", "", -1)
	spl := strings.Split(s, ";")

	colorReset := Color{
		ColorType: ColorTypeReset,
		Color:     color.RGBA{},
	}

	for i := 0; i < len(spl); i++ {
		v := spl[i]
		if v == "" {
			colors = append(colors, colorReset)
			continue
		}

		n, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}

		if 0 <= n && n <= 9 {
			c := Color{
				ColorType: AttributeMap[n],
			}
			colors = append(colors, c)
			continue
		}

		if (30 <= n && n <= 37) || (90 <= n && n <= 97) {
			c := Color{
				ColorType: ColorTypeForeground,
				Color:     color.ANSIMap[n],
			}
			colors = append(colors, c)
			continue
		}

		if (40 <= n && n <= 47) || (100 <= n && n <= 107) {
			c := Color{
				ColorType: ColorTypeBackground,
				Color:     color.ANSIMap[n],
			}
			colors = append(colors, c)
			continue
		}

		if n == 38 || n == 48 {
			v := spl[i+1]
			n2, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			if n2 != 2 && n2 != 5 {
				panic(fmt.Errorf("%v is illegal format.", s))
			}

			var ct ColorType
			if n == 38 {
				ct = ColorTypeForeground
			} else {
				ct = ColorTypeBackground
			}

			if n2 == 2 {
				// RGB指定
				var (
					r   uint64
					g   uint64
					b   uint64
					rs  = spl[i+2]
					gs  = spl[i+3]
					bs  = spl[i+4]
					err error
				)
				r, err = strconv.ParseUint(rs, 10, 8)
				if err != nil {
					panic(err)
				}
				g, err = strconv.ParseUint(gs, 10, 8)
				if err != nil {
					panic(err)
				}
				b, err = strconv.ParseUint(bs, 10, 8)
				if err != nil {
					panic(err)
				}
				c := Color{
					ColorType: ct,
					Color:     color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255},
				}
				colors = append(colors, c)
				i = i + 4
				continue
			}

			// 255指定
			c255, err := strconv.Atoi(spl[i+2])
			if err != nil {
				panic(err)
			}
			c := Color{
				ColorType: ct,
				Color:     color.Map256[c255],
			}
			colors = append(colors, c)
			i = i + 2
		}
	}
	return
}

// Prefix は文字列の先頭の要素を種類とともに返却と残り部分とともに返す。
// エスケープシーケンス系とマッチしたらマッチした部分を返す。
// エスケープシーケンス以外の場合は次のエスケープシーケンスが出現するまでを返す。
// エスケープシーケンスが出現しない場合はすべての文字列をprefixとして返す。
func Prefix(s string) (k Kind, prefix string, suffix string) {
	// 空文字列の場合
	if s == "" {
		return KindEmpty, "", ""
	}

	// 色とマッチしたら返却
	matched := reColorEscapeSequences.FindString(s)
	if matched != "" {
		l := len(matched)
		return KindColor, matched, s[l:]
	}

	// 色の出力と異なるものとマッチしたら返す
	matched = reNotColorEscapeSequences.FindString(s)
	if matched != "" {
		l := len(matched)
		return KindNotColor, matched, s[l:]
	}

	// エスケープシーケンス以外から始まる場合
	// 次のエスケープシーケンスが出現するまで取得
	const pref = "\x1b["
	if !strings.HasPrefix(s, pref) {
		idx := strings.Index(s, pref)
		if idx != -1 {
			return KindText, s[:idx], s[idx:]
		}
	}

	return KindText, s, ""
}

// StringWidth は表示上のテキストの最も幅の長い長さを返却する。
func StringWidth(s []string) (max int) {
	for _, v := range s {
		text := text(v)
		width := runewidth.StringWidth(text)
		if max < width {
			max = width
		}
	}
	return
}

// text はエスケープシーケンスを含む文字列からテキストのみを返す。
func text(s string) (ret string) {
	for s != "" {
		k, pref, suff := Prefix(s)
		if k == KindText {
			ret += pref
		}
		s = suff
	}
	return
}
