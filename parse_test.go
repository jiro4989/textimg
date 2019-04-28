package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseText(t *testing.T) {
	type TestData struct {
		desc    string
		s       string
		col     string
		matched string
		suffix  string
	}
	tds := []TestData{
		{desc: "colorBlackを取得", s: "\x1b[30mtest\x1b[0m", col: colorEscapeSequenceBlack, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorRedを取得", s: "\x1b[31mtest\x1b[0m", col: colorEscapeSequenceRed, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorGreenを取得", s: "\x1b[32mtest\x1b[0m", col: colorEscapeSequenceGreen, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorYellowを取得", s: "\x1b[33mtest\x1b[0m", col: colorEscapeSequenceYellow, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorBlueを取得", s: "\x1b[34mtest\x1b[0m", col: colorEscapeSequenceBlue, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorMagentaを取得", s: "\x1b[35mtest\x1b[0m", col: colorEscapeSequenceMagenta, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorCyanを取得", s: "\x1b[36mtest\x1b[0m", col: colorEscapeSequenceCyan, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorWhiteを取得", s: "\x1b[37mtest\x1b[0m", col: colorEscapeSequenceWhite, matched: "test", suffix: "\x1b[0m"},
		{desc: "途中で色が変わる", s: "\x1b[30mBlack\x1b[31mRed\x1b[0m", col: colorEscapeSequenceBlack, matched: "Black", suffix: "\x1b[31mRed\x1b[0m"},
		{desc: "リセット文字", s: "\x1b[0mReset\x1b[m", col: colorEscapeSequenceReset, matched: "Reset", suffix: "\x1b[m"},
		{desc: "省略リセット文字", s: "\x1b[mReset\x1b[m", col: colorEscapeSequenceResetShort, matched: "Reset", suffix: "\x1b[m"},
		{desc: "Red to Green", s: "\x1b[31mRed\x1b[m\x1b[32mGreen\x1b[m", col: colorEscapeSequenceRed, matched: "Red", suffix: "\x1b[m\x1b[32mGreen\x1b[m"},
		{desc: "エスケープシーケンスが連続する", s: "\x1b[m\x1b[31mRed", col: colorEscapeSequenceResetShort, matched: "", suffix: "\x1b[31mRed"},
		// 前提として色と直接関係のないエスケープ文字は削除していないといけない
		// ので、このテストケースは不要
		// {desc: "混合文字からcolorRedを取得", s: "\x1b[01;31m\x1b[Ktest\x1b[m\x1b[K", col: colorRed, matched: "test", suffix: "\x1b[m\x1b[K"},
	}
	for i, v := range tds {
		col, matched, suffix := parseText(v.s)
		foundNG := false
		if v.col != col {
			t.Error(fmt.Sprintf("[%2d] NG: %s color doesn't equals. expect = %s, got = %s", i, v.desc, v.col, col))
			foundNG = true
		}
		if v.matched != matched {
			t.Error(fmt.Sprintf("[%2d] NG: %s matched doesn't equals. expect = %s, got = %s", i, v.desc, v.matched, matched))
			foundNG = true
		}
		if v.suffix != suffix {
			t.Error(fmt.Sprintf("[%2d] NG: %s suffix doesn't equals. expect = %s, got = %s", i, v.desc, v.suffix, suffix))
			foundNG = true
		}
		if !foundNG {
			t.Log(fmt.Sprintf("[%2d] OK: %s", i, v.desc))
		}
	}
}

func TestGetOnlyColorEscapeSequence(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect string
	}
	tds := []TestData{
		{desc: "赤文字を取得", s: "\x1b[31mte", expect: colorEscapeSequenceRed},
		{desc: "緑文字を取得", s: "\x1b[32mte", expect: colorEscapeSequenceGreen},
		{desc: "リセット文字を取得", s: "\x1b[0mte", expect: colorEscapeSequenceReset},
		{desc: "省略リセット文字を取得", s: "\x1b[mte", expect: colorEscapeSequenceResetShort},
		{desc: "余計な文字が混じっていても取得", s: "\x1b[0;31mte", expect: colorEscapeSequenceRed},
		{desc: "余計な文字が混じっていても取得", s: "\x1b[01;31mte", expect: colorEscapeSequenceRed},
		{desc: "先頭がエスケープシーケンス以外で開始", s: "test", expect: colorEscapeSequenceNone},
	}
	for i, v := range tds {
		got := getOnlyColorEscapeSequence(v.s)
		if v.expect != got {
			t.Error(fmt.Sprintf("[%2d] NG: %s NG: expect doesn't equals. expect = %s, got = %s", i, v.desc, v.expect, got))
			continue
		}
		t.Log(fmt.Sprintf("[%2d] OK: %s", i, v.desc))
	}
}

func TestRemoveNotColorEscapeSequences(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect string
	}
	tds := []TestData{
		{desc: "消すものが何もない", s: "\x1b[31mRed\x1b[0m", expect: "\x1b[31mRed\x1b[0m"},
		{desc: "赤文字の直前に別色がわりこむ", s: "\x1b[01;31mRed\x1b[0m", expect: "\x1b[31mRed\x1b[0m"},
		{desc: "ボールドわりこみ", s: "\x1b[01;31m\x1b[KRed\x1b[0m", expect: "\x1b[31mRed\x1b[0m"},
		{desc: "途中で色がきりかわる(末尾にリセット色がない)", s: "\x1b[01;31m\x1b[KRed\x1b[32mGreen", expect: "\x1b[31mRed\x1b[32mGreen"},
		{desc: "途中で色がきりかわる(末尾にリセット色がある)", s: "\x1b[01;31m\x1b[KRed\x1b[32mGreen\x1b[0m", expect: "\x1b[31mRed\x1b[32mGreen\x1b[0m"},
		{desc: "出力消去文字  J", s: "\x1b[JTest", expect: "Test"},
		{desc: "出力消去文字 0J", s: "\x1b[0JTest", expect: "Test"},
		{desc: "出力消去文字 1J", s: "\x1b[1JTest", expect: "Test"},
		{desc: "出力消去文字  K", s: "\x1b[KTest", expect: "Test"},
		{desc: "出力消去文字 0K", s: "\x1b[0KTest", expect: "Test"},
		{desc: "出力消去文字 1K", s: "\x1b[1KTest", expect: "Test"},
	}
	for i, v := range tds {
		got := removeNotColorEscapeSequences(v.s)
		if v.expect != got {
			t.Error(fmt.Sprintf("[%2d] NG: %s NG: expect doesn't equals. expect = %s, got = %s", i, v.desc, v.expect, got))
			continue
		}
		t.Log(fmt.Sprintf("[%2d] OK: %s", i, v.desc))
	}
}

func TestClassifyString(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect ClassifiedStrings
	}
	tds := []TestData{
		{
			desc: "通常ケース", s: "\x1b[31m\x1b[KRed\x1b[0m", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[31m"},
				{class: classEscape, text: "\x1b[K"},
				{class: classText, text: "Red"},
				{class: classEscape, text: "\x1b[0m"},
			},
		},
		{
			desc: "省略記法", s: "\x1b[mTest\x1b[31mRed\x1b[0m", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[m"},
				{class: classText, text: "Test"},
				{class: classEscape, text: "\x1b[31m"},
				{class: classText, text: "Red"},
				{class: classEscape, text: "\x1b[0m"},
			},
		},
		{
			desc: "エスケープ混合", s: "\x1b[01;31mRed\x1b[0m", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[01;31m"},
				{class: classText, text: "Red"},
				{class: classEscape, text: "\x1b[0m"},
			},
		},
		{
			desc: "出力消去文字", s: "\x1b[KTest\x1b[0KTest\x1b[1KTest\x1b[JTest\x1b[1JTest", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[K"},
				{class: classText, text: "Test"},
				{class: classEscape, text: "\x1b[0K"},
				{class: classText, text: "Test"},
				{class: classEscape, text: "\x1b[1K"},
				{class: classText, text: "Test"},
				{class: classEscape, text: "\x1b[J"},
				{class: classText, text: "Test"},
				{class: classEscape, text: "\x1b[1J"},
				{class: classText, text: "Test"},
			},
		},
		{
			desc: "カーソル移動", s: "\x1b[1A\x1b[1B\x1b[1C\x1b[1D\x1b[1E\x1b[1F\x1b[1G\x1b[1;2H\x1b[1;2f", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[1A"},
				{class: classEscape, text: "\x1b[1B"},
				{class: classEscape, text: "\x1b[1C"},
				{class: classEscape, text: "\x1b[1D"},
				{class: classEscape, text: "\x1b[1E"},
				{class: classEscape, text: "\x1b[1F"},
				{class: classEscape, text: "\x1b[1G"},
				{class: classEscape, text: "\x1b[1;2H"},
				{class: classEscape, text: "\x1b[1;2f"},
			},
		},
		{
			desc: "スクロール", s: "\x1b[1S\x1b[2T", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[1S"},
				{class: classEscape, text: "\x1b[2T"},
			},
		},
	}
	for _, v := range tds {
		got := classifyString(v.s)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
