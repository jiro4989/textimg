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
		{desc: "colorBlackを取得", s: "\x1b[30mtest\x1b[0m", col: colorBlack, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorRedを取得", s: "\x1b[31mtest\x1b[0m", col: colorRed, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorGreenを取得", s: "\x1b[32mtest\x1b[0m", col: colorGreen, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorYellowを取得", s: "\x1b[33mtest\x1b[0m", col: colorYellow, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorBlueを取得", s: "\x1b[34mtest\x1b[0m", col: colorBlue, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorMagentaを取得", s: "\x1b[35mtest\x1b[0m", col: colorMagenta, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorCyanを取得", s: "\x1b[36mtest\x1b[0m", col: colorCyan, matched: "test", suffix: "\x1b[0m"},
		{desc: "colorWhiteを取得", s: "\x1b[37mtest\x1b[0m", col: colorWhite, matched: "test", suffix: "\x1b[0m"},
		{desc: "途中で色が変わる", s: "\x1b[30mBlack\x1b[31mRed\x1b[0m", col: colorBlack, matched: "Black", suffix: "\x1b[31mRed\x1b[0m"},
		// 前提として色と直接関係のないエスケープ文字は削除していないといけない
		// ので、このテストケースは不要
		// {desc: "混合文字からcolorRedを取得", s: "\x1b[01;31m\x1b[Kmtest\x1b[m\x1b[Km", col: colorRed, matched: "test", suffix: "\x1b[m\x1b[Km"},
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
		{desc: "赤文字を取得", s: "\x1b[31mte", expect: colorRed},
		{desc: "緑文字を取得", s: "\x1b[32mte", expect: colorGreen},
		{desc: "リセット文字を取得", s: "\x1b[0mte", expect: colorReset},
		{desc: "余計な文字が混じっていても取得", s: "\x1b[01;31mte", expect: colorRed},
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
		{desc: "ボールドわりこみ", s: "\x1b[01;31m\x1b[KmRed\x1b[0m", expect: "\x1b[31mRed\x1b[0m"},
		{desc: "途中で色がきりかわる(末尾にリセット色がない)", s: "\x1b[01;31m\x1b[KmRed\x1b[32mGreen", expect: "\x1b[31mRed\x1b[32mGreen"},
		{desc: "途中で色がきりかわる(末尾にリセット色がある)", s: "\x1b[01;31m\x1b[KmRed\x1b[32mGreen\x1b[0m", expect: "\x1b[31mRed\x1b[32mGreen\x1b[0m"},
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
			desc: "通常ケース", s: "\x1b[31m\x1b[KmRed\x1b[0m", expect: []ClassifiedString{
				{class: classEscape, text: "\x1b[31m"},
				{class: classEscape, text: "\x1b[Km"},
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
	}
	for _, v := range tds {
		got := classifyString(v.s)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
