package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSlideStrings(t *testing.T) {
	type TestData struct {
		desc                  string
		src, expect           []string
		lineCount, slideWidth int
		slideForever          bool
	}
	tds := []TestData{
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
			desc: "3行描画、スライド幅1、無限あり",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2", "3",
				"2", "3", "4",
				"3", "4", "5",
				"4", "5", "1",
				"5", "1", "2",
				"1", "2", "3",
			},
			lineCount:    3,
			slideWidth:   1,
			slideForever: false,
		},
	}
	for _, v := range tds {
		got := toSlideStrings(v.src, v.lineCount, v.slideWidth, v.slideForever)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
