package main

import (
	"fmt"
	"testing"
)

func TestParseText(t *testing.T) {
	type TestData struct {
		s       string
		col     string
		matched string
		suffix  string
	}
	tds := []TestData{
		{s: "\x1b[30mtest\x1b[0m", col: colorBlack, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[31mtest\x1b[0m", col: colorRed, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[32mtest\x1b[0m", col: colorGreen, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[33mtest\x1b[0m", col: colorYellow, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[34mtest\x1b[0m", col: colorBlue, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[35mtest\x1b[0m", col: colorMagenta, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[36mtest\x1b[0m", col: colorSyan, matched: "test", suffix: "\x1b[0m"},
		{s: "\x1b[37mtest\x1b[0m", col: colorWhite, matched: "test", suffix: "\x1b[0m"},
	}
	for i, v := range tds {
		col, matched, suffix := parseText(v.s)
		if v.col != col {
			t.Error(fmt.Sprintf("[%2d] NG: color doesn't equals. expect = %s, got = %s", i, v.col, col))
		}
		if v.matched != matched {
			t.Error(fmt.Sprintf("[%2d] NG: matched doesn't equals. expect = %s, got = %s", i, v.matched, matched))
		}
		if v.suffix != suffix {
			t.Error(fmt.Sprintf("[%2d] NG: suffix doesn't equals. expect = %s, got = %s", i, v.suffix, suffix))
		}
		t.Log(fmt.Sprintf("[%2d] OK:", i))
	}
}
