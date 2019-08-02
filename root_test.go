package main

import (
	"testing"

	"github.com/jiro4989/textimg/escseq"
	"github.com/stretchr/testify/assert"
)

func TestOptionColorStringToRGBA(t *testing.T) {
	type TestData struct {
		desc   string
		colstr string
		expect escseq.RGBA
	}
	tds := []TestData{
		{desc: "BLACK", colstr: "BLACK", expect: escseq.RGBABlack},
		{desc: "black", colstr: "black", expect: escseq.RGBABlack},
		{desc: "red", colstr: "red", expect: escseq.RGBARed},
		{desc: "green", colstr: "green", expect: escseq.RGBAGreen},
		{desc: "yellow", colstr: "yellow", expect: escseq.RGBAYellow},
		{desc: "blue", colstr: "blue", expect: escseq.RGBABlue},
		{desc: "magenta", colstr: "magenta", expect: escseq.RGBAMagenta},
		{desc: "cyan", colstr: "cyan", expect: escseq.RGBACyan},
		{desc: "white", colstr: "white", expect: escseq.RGBAWhite},
		{desc: "0,0,0,255", colstr: "0,0,0,255", expect: escseq.RGBA{R: 0, G: 0, B: 0, A: 255}},
		{desc: "255,255,255,255", colstr: "255,255,255,255", expect: escseq.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{desc: "0,0,0,0", colstr: "0,0,0,0", expect: escseq.RGBA{R: 0, G: 0, B: 0, A: 0}},
	}
	for _, v := range tds {
		got, err := optionColorStringToRGBA(v.colstr)
		assert.Nil(t, err, v.desc)
		assert.Equal(t, v.expect, got, v.desc)
	}

	// 異常系
	tds = []TestData{
		{desc: "不正な色文字列", colstr: "unko"},
		{desc: "RGBAの書式不正(値の数不足)", colstr: "1,2,3"},
		{desc: "RGBAの書式不正(値の数過多)", colstr: "1,2,3,4,5"},
		{desc: "RGBAの書式不正(値がない)", colstr: "1,2,3,"},
		{desc: "RGBAの書式不正(値に文字が混じっている)", colstr: "1,2,3,a"},
		{desc: "RGBAの書式不正(255以上の値)", colstr: "1,2,3,256"},
		{desc: "RGBAの書式不正(負の値)", colstr: "-1,2,3,255"},
		{desc: "RGBAの書式不正(空文字)", colstr: ""},
	}
	for _, v := range tds {
		_, err := optionColorStringToRGBA(v.colstr)
		assert.NotNil(t, err, v.desc)
	}
}
