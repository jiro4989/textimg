package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainNormal(t *testing.T) {
	testdataDir := filepath.Join(".", "testdata", "in")

	tests := []struct {
		desc string
		in   []string
		want int
	}{
		{
			desc: "正常系: 普通に描画する",
			in:   []string{"", "test"},
			want: 0,
		},
		{
			desc: "正常系: 赤色のエスケープシーケンス",
			in:   []string{"", "\x1b[31mRed\x1b[m"},
			want: 0,
		},
		{
			desc: "正常系: 赤色のエスケープシーケンスでファイル出力する",
			in:   []string{"", "\x1b[31mRed\x1b[m", "-o", outDir + "/main_normal_red.png"},
			want: 0,
		},
		{
			desc: "正常系: 文字色、背景色を変更する",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "green", "-b", "blue"},
			want: 0,
		},
		{
			desc: "正常系: カンマ区切り指定",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "0,255,0,255", "-b", "0,0,255,255"},
			want: 0,
		},
		{
			desc: "正常系: Slackアイコンサイズで生成できる",
			in:   []string{"", "Sample", "-o", outDir + "/main_slack.png", "--slack"},
			want: 0,
		},
		{
			desc: "正常系: FontIndex, EmojiFontIndexを指定できる",
			in:   []string{"", "Sample", "-o", outDir + "/main_font_index.png", "-x", "0", "-X", "0"},
			want: 0,
		},
		{
			desc: "正常系: SlackアイコンサイズでアニメーションGIFを生成できる",
			in:   []string{"", "Sample", "-o", outDir + "/main_slack_anim.gif", "--slack", "-a"},
			want: 0,
		},
		{
			desc: "異常系: 不正な文字色",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "gggg", "-b", "blue"},
			want: -1,
		},
		{
			desc: "異常系: 不正な背景色",
			in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "green", "-b", "bbbb"},
			want: -1,
		},
		{
			desc: "異常系: 不正なフォントを指定",
			in:   []string{"", "Sample", "-o", outDir + "/illegal_case1.png", "-f", filepath.Join(testdataDir, "illegal_font.ttc")},
			want: -1,
		},
		{
			desc: "異常系: 不正な絵文字フォントを指定",
			in:   []string{"", "Sample", "-o", outDir + "/illegal_case2.png", "-e", filepath.Join(testdataDir, "illegal_font.ttc")},
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
