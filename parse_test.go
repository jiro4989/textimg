package main

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseColorEscapeSequence(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect colorEscapeSequences
	}
	tds := []TestData{
		{desc: "ANSI 文字色 黒", s: "\x1b[30m", expect: colorEscapeSequences{{colorType: colorTypeForeground, color: colorRGBABlack}}},
		{desc: "ANSI 文字色 赤", s: "\x1b[31m", expect: colorEscapeSequences{{colorType: colorTypeForeground, color: colorRGBARed}}},
		{desc: "ANSI 文字色 白", s: "\x1b[37m", expect: colorEscapeSequences{{colorType: colorTypeForeground, color: colorRGBAWhite}}},
		{desc: "ANSI 背景色 黒", s: "\x1b[40m", expect: colorEscapeSequences{{colorType: colorTypeBackground, color: colorRGBABlack}}},
		{desc: "ANSI 背景色 赤", s: "\x1b[41m", expect: colorEscapeSequences{{colorType: colorTypeBackground, color: colorRGBARed}}},
		{desc: "ANSI 背景色 白", s: "\x1b[47m", expect: colorEscapeSequences{{colorType: colorTypeBackground, color: colorRGBAWhite}}},
		{desc: "ANSI リセット", s: "\x1b[0m", expect: colorEscapeSequences{{colorType: colorTypeReset, color: color.RGBA{}}}},
		{desc: "ANSI リセット(省略記法)", s: "\x1b[m", expect: colorEscapeSequences{{colorType: colorTypeReset, color: color.RGBA{}}}},
		{desc: "ANSI リセットと文字色と背景色", s: "\x1b[0;31;42;01m", expect: colorEscapeSequences{
			{colorType: colorTypeReset, color: color.RGBA{}},
			{colorType: colorTypeForeground, color: colorRGBARed},
			{colorType: colorTypeBackground, color: colorRGBAGreen},
			{colorType: colorTypeBold, color: color.RGBA{}},
		}},
		{desc: "拡張256色記法 文字色 赤", s: "\x1b[38;5;25m", expect: colorEscapeSequences{{colorType: colorTypeForeground, color: color256Map[25]}}},
		{desc: "拡張256色記法 背景色 赤", s: "\x1b[48;5;25m", expect: colorEscapeSequences{{colorType: colorTypeBackground, color: color256Map[25]}}},
		{desc: "拡張256色RGB記法 文字色 赤", s: "\x1b[38;2;255;0;0m", expect: colorEscapeSequences{{colorType: colorTypeForeground, color: colorRGBARed}}},
		{desc: "拡張256色RGB記法 背景色 赤", s: "\x1b[48;2;255;0;0m", expect: colorEscapeSequences{{colorType: colorTypeBackground, color: colorRGBARed}}},
		{desc: "拡張記法混在", s: "\x1b[38;5;25;48;2;255;0;0m", expect: colorEscapeSequences{
			{colorType: colorTypeForeground, color: color256Map[25]},
			{colorType: colorTypeBackground, color: colorRGBARed},
		}},
	}
	for _, v := range tds {
		got := parseColorEscapeSequence(v.s)
		assert.Equal(t, v.expect, got, v.desc)
	}
}

func TestGetPrefix(t *testing.T) {
	type TestData struct {
		desc         string
		s            string
		expectKind   kind
		expectPrefix string
		expectSuffix string
	}
	tds := []TestData{
		// 色 (\x1b[Nm)
		{desc: "赤文字", s: "\x1b[31mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[31m", expectSuffix: "Test"},
		{desc: "緑文字", s: "\x1b[32mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[32m", expectSuffix: "Test"},
		{desc: "リセット文字", s: "\x1b[0mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[0m", expectSuffix: "Test"},
		{desc: "リセット文字(省略記法)", s: "\x1b[mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[m", expectSuffix: "Test"},
		{desc: "セミコロン区切り", s: "\x1b[31;42mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[31;42m", expectSuffix: "Test"},
		{desc: "拡張256色記法(文字色)", s: "\x1b[38;5;255mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[38;5;255m", expectSuffix: "Test"},
		{desc: "拡張256色記法(背景色)", s: "\x1b[48;5;255mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[48;5;255m", expectSuffix: "Test"},
		{desc: "拡張256色RGB指定記法(文字色)", s: "\x1b[38;2;1;2;3mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[38;2;1;2;3m", expectSuffix: "Test"},
		{desc: "拡張256色RGB指定記法(背景色)", s: "\x1b[48;2;1;2;3mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[48;2;1;2;3m", expectSuffix: "Test"},
		{desc: "いろいろ混在", s: "\x1b[5;38;5;124;48;2;1;2;3mTest", expectKind: kindEscapeSequenceColor, expectPrefix: "\x1b[5;38;5;124;48;2;1;2;3m", expectSuffix: "Test"},
		// 色以外のエスケープシーケンス
		{desc: "A", s: "\x1b[1ATest", expectKind: kindEscapeSequenceNotColor, expectPrefix: "\x1b[1A", expectSuffix: "Test"},
		{desc: "H", s: "\x1b[1HTest", expectKind: kindEscapeSequenceNotColor, expectPrefix: "\x1b[1H", expectSuffix: "Test"},
		{desc: "f", s: "\x1b[1fTest", expectKind: kindEscapeSequenceNotColor, expectPrefix: "\x1b[1f", expectSuffix: "Test"},
		{desc: "K", s: "\x1b[1KTest", expectKind: kindEscapeSequenceNotColor, expectPrefix: "\x1b[1K", expectSuffix: "Test"},
		// エスケープシーケンス以外
		{desc: "エスケープシーケンス以外 空文字列", s: "", expectKind: kindEmpty, expectPrefix: "", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 通常文字列", s: "a", expectKind: kindText, expectPrefix: "a", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 通常文字列", s: "abc\x1b[31mTest", expectKind: kindText, expectPrefix: "abc", expectSuffix: "\x1b[31mTest"},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "x1b[31mTest", expectKind: kindText, expectPrefix: "x1b[31mTest", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "x1b[31amTest", expectKind: kindText, expectPrefix: "x1b[31amTest", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "\\x1b[31mTest", expectKind: kindText, expectPrefix: "\\x1b[31mTest", expectSuffix: ""},
		{desc: "エスケープシーケンス以外 エスケープシーケンスっぽい文字列", s: "\x1b[31xTest", expectKind: kindText, expectPrefix: "\x1b[31xTest", expectSuffix: ""},
	}
	for _, v := range tds {
		gk, gp, gs := getPrefix(v.s)
		assert.Equal(t, v.expectKind, gk, v.desc)
		assert.Equal(t, v.expectPrefix, gp, v.desc)
		assert.Equal(t, v.expectSuffix, gs, v.desc)
	}
}

func TestGetText(t *testing.T) {
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
		got := getText(v.s)
		assert.Equal(t, v.expect, got, v.desc)
	}
}

func TestMaxStringWidth(t *testing.T) {
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
		got := maxStringWidth(v.s)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
