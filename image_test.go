package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEncodeFormat(t *testing.T) {
	type TestData struct {
		desc   string
		path   string
		expect encodeFormat
	}
	tds := []TestData{
		{desc: "png", path: "a.png", expect: encodeFormatPNG},
		{desc: "png", path: "jpeg.png", expect: encodeFormatPNG},
		{desc: "png", path: "gif.jpeg.png", expect: encodeFormatPNG},
		{desc: "png", path: "/tmp/a.png", expect: encodeFormatPNG},
		{desc: "Png", path: "/tmp/a.Png", expect: encodeFormatPNG},
		{desc: "PNG", path: "/tmp/a.PNG", expect: encodeFormatPNG},
		{desc: "jpg", path: "/tmp/a.jpg", expect: encodeFormatJPG},
		{desc: "Jpg", path: "/tmp/a.Jpg", expect: encodeFormatJPG},
		{desc: "JPG", path: "/tmp/a.JPG", expect: encodeFormatJPG},
		{desc: "jpeg", path: "/tmp/a.jpeg", expect: encodeFormatJPG},
		{desc: "Jpeg", path: "/tmp/a.Jpeg", expect: encodeFormatJPG},
		{desc: "JPEG", path: "/tmp/a.JPEG", expect: encodeFormatJPG},
		{desc: "gif", path: "/tmp/a.gif", expect: encodeFormatGIF},
		{desc: "Gif", path: "/tmp/a.Gif", expect: encodeFormatGIF},
		{desc: "GIF", path: "/tmp/a.GIF", expect: encodeFormatGIF},
		{desc: "pathが空のときはpng", path: "", expect: encodeFormatPNG},
	}
	for _, v := range tds {
		got, err := getEncodeFormat(v.path)
		assert.Nil(t, err, v.desc)
		assert.Equal(t, v.expect, got, v.desc)
	}

	// 異常系
	got, err := getEncodeFormat("a.webp")
	assert.NotNil(t, err, "エラーが起こるべき")
	assert.Equal(t, encodeFormat(-1), got, "-1になるべき")
}

func TestIsEmojiCodePoint(t *testing.T) {
	type TestData struct {
		desc   string
		r      rune
		expect bool
	}
	tds := []TestData{
		{desc: "# == false", r: []rune("#")[0], expect: false},
		{desc: "* == false", r: []rune("*")[0], expect: false},
		{desc: "0 == false", r: []rune("0")[0], expect: false},
		{desc: "1 == false", r: []rune("1")[0], expect: false},
		{desc: "2 == false", r: []rune("2")[0], expect: false},
		{desc: "3 == false", r: []rune("3")[0], expect: false},
		{desc: "4 == false", r: []rune("4")[0], expect: false},
		{desc: "5 == false", r: []rune("5")[0], expect: false},
		{desc: "6 == false", r: []rune("6")[0], expect: false},
		{desc: "7 == false", r: []rune("7")[0], expect: false},
		{desc: "8 == false", r: []rune("8")[0], expect: false},
		{desc: "9 == false", r: []rune("9")[0], expect: false},
		{desc: "😚 == true", r: []rune("😚")[0], expect: true},
	}
	for _, v := range tds {
		got := isEmojiCodePoint(v.r)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
