package parser

import (
	"testing"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/token"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type TestData struct {
		desc    string
		s       string
		want    token.Tokens
		wantErr bool
	}
	tds := []TestData{
		{
			desc: "正常系: 黒",
			s:    "\x1b[30m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBABlack,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 赤とテキストとリセット",
			s:    "\x1b[31mhello world\x1b[0m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBARed,
				},
				{
					Kind: token.KindText,
					Text: "hello world",
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeReset,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 前景色と背景色の同時指定",
			s:    "\x1b[32;43mhello world",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBAGreen,
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeBackground,
					Color:     color.RGBAYellow,
				},
				{
					Kind: token.KindText,
					Text: "hello world",
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 0埋めありの指定",
			s:    "\x1b[032;00043mhello world",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBAGreen,
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeBackground,
					Color:     color.RGBAYellow,
				},
				{
					Kind: token.KindText,
					Text: "hello world",
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 拡張系 256色",
			s:    "\x1b[38;5;1m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.Map256[1],
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 拡張系 RGB指定",
			s:    "\x1b[48;2;1;2;3m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeBackground,
					Color: color.RGBA{
						R: 1,
						G: 2,
						B: 3,
						A: 255,
					},
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 拡張系の混在",
			s:    "\x1b[38;5;2;48;2;1;2;3mこんばんは",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.Map256[2],
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeBackground,
					Color: color.RGBA{
						R: 1,
						G: 2,
						B: 3,
						A: 255,
					},
				},
				{
					Kind: token.KindText,
					Text: "こんばんは",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tds {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			got, err := Parse(tt.s)

			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
				return
			}
			assert.Equal(tt.want, got)
		})
	}
}
