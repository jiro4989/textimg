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

func TestIsExceptionallyCodePoint(t *testing.T) {
	type TestData struct {
		desc   string
		r      rune
		expect bool
	}
	tds := []TestData{
		{desc: "# == true", r: []rune("#")[0], expect: true},
		{desc: "* == true", r: []rune("*")[0], expect: true},
		{desc: "0 == true", r: []rune("0")[0], expect: true},
		{desc: "1 == true", r: []rune("1")[0], expect: true},
		{desc: "2 == true", r: []rune("2")[0], expect: true},
		{desc: "3 == true", r: []rune("3")[0], expect: true},
		{desc: "4 == true", r: []rune("4")[0], expect: true},
		{desc: "5 == true", r: []rune("5")[0], expect: true},
		{desc: "6 == true", r: []rune("6")[0], expect: true},
		{desc: "7 == true", r: []rune("7")[0], expect: true},
		{desc: "8 == true", r: []rune("8")[0], expect: true},
		{desc: "9 == true", r: []rune("9")[0], expect: true},
		{desc: "© == true", r: []rune("©")[0], expect: true},
		{desc: "®️ == true", r: []rune("®️")[0], expect: true},
		{desc: "この次", r: rune(0x00AF), expect: false},
		{desc: "😚 == false", r: []rune("😚")[0], expect: false},
	}
	for _, v := range tds {
		got := isExceptionallyCodePoint(v.r)
		assert.Equal(t, v.expect, got, v.desc)
	}
}
