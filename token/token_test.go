package token

import (
	"testing"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/stretchr/testify/assert"
)

func TestToken_MaxStringWidth(t *testing.T) {
	tests := []struct {
		desc string
		t    Tokens
		want int
	}{
		{
			desc: "正常系: hello = 5",
			t: Tokens{
				{
					Kind: KindText,
					Text: "hello",
				},
			},
			want: 5,
		},
		{
			desc: "正常系: あいうえお = 10",
			t: Tokens{
				{
					Kind: KindText,
					Text: "あいうえお",
				},
			},
			want: 10,
		},
		{
			desc: "正常系: he\nllo = 3",
			t: Tokens{
				{
					Kind: KindText,
					Text: "he\nllo",
				},
			},
			want: 3,
		},
		{
			desc: "正常系: hel\nlo = 3",
			t: Tokens{
				{
					Kind: KindText,
					Text: "hel\nlo",
				},
			},
			want: 3,
		},
		{
			desc: "正常系: aREDb = 5",
			t: Tokens{
				{
					Kind: KindText,
					Text: "a",
				},
				{
					Kind:      KindColor,
					ColorType: ColorTypeForeground,
					Color:     color.RGBARed,
				},
				{
					Kind: KindText,
					Text: "RED",
				},
				{
					Kind:      KindColor,
					ColorType: ColorTypeReset,
				},
				{
					Kind: KindText,
					Text: "b",
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			got := tt.t.MaxStringWidth()
			assert.Equal(tt.want, got)
		})
	}
}

func TestToken_StringLines(t *testing.T) {
	tests := []struct {
		desc string
		t    Tokens
		want []string
	}{
		{
			desc: "正常系: hello",
			t: Tokens{
				{
					Kind: KindText,
					Text: "hello",
				},
			},
			want: []string{"hello"},
		},
		{
			desc: "正常系: hel\nlo\nworld",
			t: Tokens{
				{
					Kind: KindText,
					Text: "hel\nlo\nworld",
				},
			},
			want: []string{"hel", "lo", "world"},
		},
		{
			desc: "正常系: aREDb = 5",
			t: Tokens{
				{
					Kind: KindText,
					Text: "a",
				},
				{
					Kind:      KindColor,
					ColorType: ColorTypeForeground,
					Color:     color.RGBARed,
				},
				{
					Kind: KindText,
					Text: "RED",
				},
				{
					Kind:      KindColor,
					ColorType: ColorTypeReset,
				},
				{
					Kind: KindText,
					Text: "b",
				},
			},
			want: []string{"aREDb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			got := tt.t.StringLines()
			assert.Equal(tt.want, got)
		})
	}
}
