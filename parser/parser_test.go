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
			desc: "正常系: 90系と100系",
			s:    "\x1b[90;100m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBADarkGray,
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeBackground,
					Color:     color.RGBADarkGray,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 黒赤緑",
			s:    "\x1b[30m\x1b[31m\x1b[32m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBABlack,
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBARed,
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBAGreen,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 赤とテキストとリセット",
			s:    "\x1b[31m\n hello\tworld \n\x1b[0m",
			want: token.Tokens{
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeForeground,
					Color:     color.RGBARed,
				},
				{
					Kind: token.KindText,
					Text: "\n hello\tworld \n",
				},
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeReset,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 前景色と背景色の同時指定 + リセット省略系",
			s:    "\x1b[32;43mhello world\x1b[m",
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
				{
					Kind:      token.KindColor,
					ColorType: token.ColorTypeReset,
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
		{
			desc: "正常系: 拡張系の混在 + 0埋め",
			s:    "\x1b[038;005;002;048;002;001;002;003mx1bこんば\nんはx1b",
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
					Text: "x1bこんば\nんはx1b",
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 関係ないエスケープシーケンス系無視される",
			s:    "\x1b[1A寿司",
			want: token.Tokens{
				{
					Kind: token.KindText,
					Text: "寿司",
				},
			},
			wantErr: false,
		},
		{
			desc:    "正常系: 空文字列の場合は空",
			s:       "",
			want:    nil,
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
