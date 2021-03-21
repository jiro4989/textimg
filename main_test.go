package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainNormal(t *testing.T) {
	tests := []struct {
		desc string
		in   []string
		want int
	}{
		{
			desc: "普通に描画する",
			in:   []string{"", "test"},
			want: 0,
		},
		{
			desc: "赤色のエスケープシーケンス",
			in:   []string{"", "\x1b[31mRed\x1b[m"},
			want: 0,
		},
		{
			desc: "赤色のエスケープシーケンスでファイル出力する",
			in:   []string{"", "\x1b[31mRed\x1b[m", "-o", outDir + "/main_normal_red.png"},
			want: 0,
		},
		{
			desc: "文字色、背景色を変更する",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "green", "-b", "blue"},
			want: 0,
		},
		{
			desc: "カンマ区切り指定",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "0,255,0,255", "-b", "0,0,255,255"},
			want: 0,
		},
		{
			desc: "Slackアイコンサイズで生成できる",
			in:   []string{"", "Sample", "-o", outDir + "/main_slack.png", "--slack"},
			want: 0,
		},
		{
			desc: "SlackアイコンサイズでアニメーションGIFを生成できる",
			in:   []string{"", "Sample", "-o", outDir + "/main_slack_anim.gif", "--slack", "-a"},
			want: 0,
		},
		{
			desc: "不正な文字色",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "gggg", "-b", "blue"},
			want: -1,
		},
		{
			desc: "不正な背景色",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "green", "-b", "bbbb"},
			want: -1,
		},
		{
			desc: "異常系: 存在しないFontを指定",
			in:   []string{"", "Sample", "-o", outDir + "/illegal_case1.png", "-f", "./寿司.txt"},
			want: -1,
		},
		{
			desc: "異常系: 存在しない絵文字Fontを指定",
			in:   []string{"", "Sample", "-o", outDir + "/illegal_case2.png", "-e", "./寿司.txt"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			os.Args = tt.in
			got := Main()
			assert.Equal(t, tt.want, got, tt.desc)
		})
	}

}
