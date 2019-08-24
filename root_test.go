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
		t.Run(v.desc, func(t *testing.T) {
			got, err := optionColorStringToRGBA(v.colstr)
			assert.Nil(t, err, v.desc)
			assert.Equal(t, v.expect, got, v.desc)
		})
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
		t.Run(v.desc, func(t *testing.T) {
			_, err := optionColorStringToRGBA(v.colstr)
			assert.NotNil(t, err, v.desc)
		})
	}
}

func TestToSlideStrings(t *testing.T) {
	type TestData struct {
		desc                  string
		src, expect           []string
		lineCount, slideWidth int
		slideForever          bool
	}
	tds := []TestData{
		{
			desc: "2行描画、スライド幅1、無限なし",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2",
				"2", "3",
				"3", "4",
				"4", "5",
			},
			lineCount:    2,
			slideWidth:   1,
			slideForever: false,
		},
		{
			desc: "2行描画、スライド幅2、無限なし",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2",
				"3", "4",
				"5", "",
			},
			lineCount:    2,
			slideWidth:   2,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅1、無限なし",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2", "3",
				"2", "3", "4",
				"3", "4", "5",
			},
			lineCount:    3,
			slideWidth:   1,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅2、無限なし、不足あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"3", "4", "5",
				"5", "6", "",
			},
			lineCount:    3,
			slideWidth:   2,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅2、無限なし、不足なし",
			src:  []string{"1", "2", "3", "4", "5", "6", "7"},
			expect: []string{
				"1", "2", "3",
				"3", "4", "5",
				"5", "6", "7",
			},
			lineCount:    3,
			slideWidth:   2,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅3、無限なし、不足なし",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅3、無限なし、不足あり",
			src:  []string{"1", "2", "3", "4", "5", "6", "7"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
				"7", "", "",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅3、無限なし、不足あり",
			src:  []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
				"7", "8", "",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: false,
		},
		{
			desc: "2行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2",
				"3", "4",
				"5", "1",
			},
			lineCount:    2,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "2行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2",
				"3", "4",
				"5", "6",
			},
			lineCount:    2,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅1、無限あり",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2", "3",
				"2", "3", "4",
				"3", "4", "5",
				"4", "5", "1",
				"5", "1", "2",
			},
			lineCount:    3,
			slideWidth:   1,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅1、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"2", "3", "4",
				"3", "4", "5",
				"4", "5", "6",
				"5", "6", "1",
				"6", "1", "2",
			},
			lineCount:    3,
			slideWidth:   1,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2", "3",
				"3", "4", "5",
				"5", "1", "2",
			},
			lineCount:    3,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"3", "4", "5",
				"5", "6", "1",
			},
			lineCount:    3,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅3、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅3、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6", "7"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
				"7", "1", "2",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: true,
		},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := toSlideStrings(v.src, v.lineCount, v.slideWidth, v.slideForever)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}

func TestRemoveZeroWidthCharacters(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect string
	}
	tds := []TestData{
		{desc: "Zero width space (U+200B)が削除される", s: "A\u200bB", expect: "AB"},
		{desc: "Zero width joiner (U+200C)が削除される", s: "A\u200cB", expect: "AB"},
		{desc: "Zero width joiner (U+200D)が削除される", s: "A\u200dB", expect: "AB"},
		{desc: "U+200B ~ U+200Dが削除される", s: "あ\u200bい\u200cう\u200dえ", expect: "あいうえ"},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := removeZeroWidthCharacters(v.s)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}
