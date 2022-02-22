package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseColor(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect Colors
	}
	tds := []TestData{
		{desc: "ANSI 文字色 黒", s: "\x1b[30m", expect: Colors{{ColorType: ColorTypeForeground, Color: RGBABlack}}},
		{desc: "ANSI 文字色 赤", s: "\x1b[31m", expect: Colors{{ColorType: ColorTypeForeground, Color: RGBARed}}},
		{desc: "ANSI 文字色 灰", s: "\x1b[37m", expect: Colors{{ColorType: ColorTypeForeground, Color: RGBALightGray}}},
		{desc: "ANSI 文字色 黒", s: "\x1b[90m", expect: Colors{{ColorType: ColorTypeForeground, Color: RGBADarkGray}}},
		{desc: "ANSI 文字色 白", s: "\x1b[97m", expect: Colors{{ColorType: ColorTypeForeground, Color: RGBAWhite}}},

		{desc: "ANSI 背景色 黒", s: "\x1b[40m", expect: Colors{{ColorType: ColorTypeBackground, Color: RGBABlack}}},
		{desc: "ANSI 背景色 赤", s: "\x1b[41m", expect: Colors{{ColorType: ColorTypeBackground, Color: RGBARed}}},
		{desc: "ANSI 背景色 灰", s: "\x1b[47m", expect: Colors{{ColorType: ColorTypeBackground, Color: RGBALightGray}}},
		{desc: "ANSI 背景色 黒", s: "\x1b[100m", expect: Colors{{ColorType: ColorTypeBackground, Color: RGBADarkGray}}},
		{desc: "ANSI 背景色 白", s: "\x1b[107m", expect: Colors{{ColorType: ColorTypeBackground, Color: RGBAWhite}}},
		{desc: "ANSI リセット", s: "\x1b[0m", expect: Colors{{ColorType: ColorTypeReset, Color: RGBA{}}}},
		{desc: "ANSI リセット(省略記法)", s: "\x1b[m", expect: Colors{{ColorType: ColorTypeReset, Color: RGBA{}}}},
		{desc: "ANSI リセットと文字色と背景色", s: "\x1b[0;31;42;01m", expect: Colors{
			{ColorType: ColorTypeReset, Color: RGBA{}},
			{ColorType: ColorTypeForeground, Color: RGBARed},
			{ColorType: ColorTypeBackground, Color: RGBAGreen},
			{ColorType: ColorTypeBold, Color: RGBA{}},
		}},
		{desc: "拡張256色記法 文字色 赤", s: "\x1b[38;5;25m", expect: Colors{{ColorType: ColorTypeForeground, Color: Map256[25]}}},
		{desc: "拡張256色記法 背景色 赤", s: "\x1b[48;5;25m", expect: Colors{{ColorType: ColorTypeBackground, Color: Map256[25]}}},
		{desc: "拡張256色RGB記法 文字色 赤", s: "\x1b[38;2;255;0;0m", expect: Colors{{ColorType: ColorTypeForeground, Color: RGBARed}}},
		{desc: "拡張256色RGB記法 背景色 赤", s: "\x1b[48;2;255;0;0m", expect: Colors{{ColorType: ColorTypeBackground, Color: RGBARed}}},
		{desc: "拡張記法混在", s: "\x1b[38;5;25;48;2;255;0;0m", expect: Colors{
			{ColorType: ColorTypeForeground, Color: Map256[25]},
			{ColorType: ColorTypeBackground, Color: RGBARed},
		}},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := ParseColor(v.s)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}

func TestPrefix(t *testing.T) {
	type TestData struct {
		desc         string
		s            string
		expectKind   Kind
		expectPrefix string
		expectSuffix string
	}
	tds := []TestData{
		// 色 (\x1b[Nm)
		{desc: "赤文字", s: "\x1b[31mTest", expectKind: KindColor, expectPrefix: "\x1b[31m", expectSuffix: "Test"},
		{desc: "緑文字", s: "\x1b[32mTest", expectKind: KindColor, expectPrefix: "\x1b[32m", expectSuffix: "Test"},
		{desc: "リセット文字", s: "\x1b[0mTest", expectKind: KindColor, expectPrefix: "\x1b[0m", expectSuffix: "Test"},
		{desc: "リセット文字(省略記法)", s: "\x1b[mTest", expectKind: KindColor, expectPrefix: "\x1b[m", expectSuffix: "Test"},
		{desc: "セミコロン区切り", s: "\x1b[31;42mTest", expectKind: KindColor, expectPrefix: "\x1b[31;42m", expectSuffix: "Test"},
		{desc: "拡張256色記法(文字色)", s: "\x1b[38;5;255mTest", expectKind: KindColor, expectPrefix: "\x1b[38;5;255m", expectSuffix: "Test"},
		{desc: "拡張256色記法(背景色)", s: "\x1b[48;5;255mTest", expectKind: KindColor, expectPrefix: "\x1b[48;5;255m", expectSuffix: "Test"},
		{desc: "拡張256色RGB指定記法(文字色)", s: "\x1b[38;2;1;2;3mTest", expectKind: KindColor, expectPrefix: "\x1b[38;2;1;2;3m", expectSuffix: "Test"},
		{desc: "拡張256色RGB指定記法(背景色)", s: "\x1b[48;2;1;2;3mTest", expectKind: KindColor, expectPrefix: "\x1b[48;2;1;2;3m", expectSuffix: "Test"},
		{desc: "いろいろ混在", s: "\x1b[5;38;5;124;48;2;1;2;3mTest", expectKind: KindColor, expectPrefix: "\x1b[5;38;5;124;48;2;1;2;3m", expectSuffix: "Test"},
		// 色以外のエスケープシーケンス
		{desc: "A", s: "\x1b[1ATest", expectKind: KindNotColor, expectPrefix: "\x1b[1A", expectSuffix: "Test"},
		{desc: "H", s: "\x1b[1HTest", expectKind: KindNotColor, expectPrefix: "\x1b[1H", expectSuffix: "Test"},
		{desc: "f", s: "\x1b[1fTest", expectKind: KindNotColor, expectPrefix: "\x1b[1f", expectSuffix: "Test"},
		{desc: "K", s: "\x1b[1KTest", expectKind: KindNotColor, expectPrefix: "\x1b[1K", expectSuffix: "Test"},
		// エスケープシーケンス以外
		{desc: "エスケープシーケンス以外 空文字列", s: "", expectKind: KindEmpty, expectPrefix: "", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 通常文字列", s: "a", expectKind: KindText, expectPrefix: "a", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 通常文字列", s: "abc\x1b[31mTest", expectKind: KindText, expectPrefix: "abc", expectSuffix: "\x1b[31mTest"},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "x1b[31mTest", expectKind: KindText, expectPrefix: "x1b[31mTest", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "x1b[31amTest", expectKind: KindText, expectPrefix: "x1b[31amTest", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "\\x1b[31mTest", expectKind: KindText, expectPrefix: "\\x1b[31mTest", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "\x1b[31xTest", expectKind: KindText, expectPrefix: "\x1b[31xTest", expectSuffix: ""},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			gk, gp, gs := Prefix(v.s)
			assert.Equal(t, v.expectKind, gk, v.desc)
			assert.Equal(t, v.expectPrefix, gp, v.desc)
			assert.Equal(t, v.expectSuffix, gs, v.desc)
		})
	}
}

func TestText(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect string
	}
	tds := []TestData{
		{desc: "赤文字", s: "\x1b[31mTest\x1b[0mTest2", expect: "TestTest2"},
		{desc: "リセット文字", s: "\x1b[mTest\x1b[0mTest2", expect: "TestTest2"},
		{desc: "255文字", s: "\x1b[38;5;255mTest\x1b[0mTest2", expect: "TestTest2"},
		{desc: "RGB文字", s: "\x1b[48;2;255;2;2mTest\x1b[0mTest2", expect: "TestTest2"},
		{desc: "色以外文字", s: "\x1b[1ATest\x1b[2BTest2", expect: "TestTest2"},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := text(v.s)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}

func TestStringWidth(t *testing.T) {
	type TestData struct {
		desc   string
		s      []string
		expect int
	}
	tds := []TestData{
		{desc: "エスケープ文字 有", s: []string{"\x1b[31mTest\x1b[0mTest2"}, expect: 9},
		{desc: "エスケープ文字 有", s: []string{"\x1b[31mTest", "\x1b[0mTest2"}, expect: 5},
		{desc: "エスケープ文字 無", s: []string{"TestTest2"}, expect: 9},
		{desc: "マルチバイト文字", s: []string{"あいうえお"}, expect: 10},
		{desc: "空文字列", s: []string{""}, expect: 0},
		{desc: "空配列", s: []string{}, expect: 0},
		{desc: "nil", s: nil, expect: 0},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := StringWidth(v.s)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}
