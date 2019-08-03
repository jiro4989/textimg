package ioimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
