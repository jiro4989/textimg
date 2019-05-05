package main

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

type (
	kind                int // エスケープシーケンスの種類
	colorType           int // 文字色か背景色か
	colorEscapeSequence struct {
		colorType colorType
		color     color.RGBA
	}
	colorEscapeSequences []colorEscapeSequence
)

const (
	kindEmpty kind = iota
	kindText
	kindEscapeSequenceColor
	kindEscapeSequenceNotColor

	colorTypeReset       colorType = iota // \x1b[0m 指定をリセット
	colorTypeBold                         // \x1b[1m 太字
	colorTypeDim                          // \x1b[2m 薄く表示
	colorTypeItalic                       // \x1b[3m イタリック
	colorTypeUnderline                    // \x1b[4m アンダーライン
	colorTypeBlink                        // \x1b[5m ブリンク
	colorTypeSpeedyBlink                  // \x1b[6m 高速ブリンク
	colorTypeReverse                      // \x1b[7m 文字色と背景色の反転
	colorTypeHide                         // \x1b[8m 表示を隠す
	colorTypeDelete                       // \x1b[9m 取り消し
	colorTypeForeground
	colorTypeBackground
)

var (
	// 出力の仕方や色を制御するエスケープシーケンスにマッチする
	reColorEscapeSequences = regexp.MustCompile(`^\x1b\[[\d;]*m`)
	// 出力の仕方や色を制御するエスケープシーケンス以外のエスケープシーケンスに
	// マッチする
	reNotColorEscapeSequences = regexp.MustCompile(`^\x1b\[\d*[A-HfSTJK]`)
)

// parseColorEscapeSequence は色のエスケープシーケンスを解析してRGBAに変換する。
func parseColorEscapeSequence(s string) (colors colorEscapeSequences) {
	s = strings.Replace(s, "\x1b[", "", -1)
	s = strings.Replace(s, "m", "", -1)
	spl := strings.Split(s, ";")

	colorReset := colorEscapeSequence{
		colorType: colorTypeReset,
		color:     color.RGBA{},
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
			c := colorEscapeSequence{
				colorType: colorAttributeMap[n],
			}
			colors = append(colors, c)
			continue
		}

		if (30 <= n && n <= 37) || (90 <= n && n <= 97) {
			c := colorEscapeSequence{
				colorType: colorTypeForeground,
				color:     colorANSIMap[n],
			}
			colors = append(colors, c)
			continue
		}

		if (40 <= n && n <= 47) || (100 <= n && n <= 107) {
			c := colorEscapeSequence{
				colorType: colorTypeBackground,
				color:     colorANSIMap[n],
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
				panic(errors.New(fmt.Sprintf("%v is illegal format.", s)))
			}

			var ct colorType
			if n == 38 {
				ct = colorTypeForeground
			} else {
				ct = colorTypeBackground
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
				c := colorEscapeSequence{
					colorType: ct,
					color:     color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255},
				}
				colors = append(colors, c)
				i = i + 4
				continue
			} else {
				// 255指定
				c255, err := strconv.Atoi(spl[i+2])
				if err != nil {
					panic(err)
				}
				c := colorEscapeSequence{
					colorType: ct,
					color:     color256Map[c255],
				}
				colors = append(colors, c)
				i = i + 2
				continue
			}
		}
	}
	return
}

// getPrefix は文字列の先頭の要素を種類とともに返却と残り部分とともに返す。
// エスケープシーケンス系とマッチしたらマッチした部分を返す。
// エスケープシーケンス以外の場合は次のエスケープシーケンスが出現するまでを返す。
// エスケープシーケンスが出現しない場合はすべての文字列をprefixとして返す。
func getPrefix(s string) (k kind, prefix string, suffix string) {
	// 空文字列の場合
	if s == "" {
		return kindEmpty, "", ""
	}

	// 色とマッチしたら返却
	matched := reColorEscapeSequences.FindString(s)
	if matched != "" {
		l := len(matched)
		return kindEscapeSequenceColor, matched, s[l:]
	}

	// 色の出力と異なるものとマッチしたら返す
	matched = reNotColorEscapeSequences.FindString(s)
	if matched != "" {
		l := len(matched)
		return kindEscapeSequenceNotColor, matched, s[l:]
	}

	// エスケープシーケンス以外から始まる場合
	// 次のエスケープシーケンスが出現するまで取得
	const pref = "\x1b["
	if !strings.HasPrefix(s, pref) {
		idx := strings.Index(s, pref)
		if idx != -1 {
			return kindText, s[:idx], s[idx:]
		}
	}

	return kindText, s, ""
}

// getText はエスケープシーケンスを含む文字列からテキストのみを返す。
func getText(s string) (ret string) {
	for s != "" {
		k, pref, suff := getPrefix(s)
		if k == kindText {
			ret += pref
		}
		s = suff
	}
	return
}

// maxStringWidth は表示上のテキストの最も幅の長い長さを返却する。
func maxStringWidth(s []string) (max int) {
	for _, v := range s {
		text := getText(v)
		width := runewidth.StringWidth(text)
		if max < width {
			max = width
		}
	}
	return
}
