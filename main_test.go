package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	outDir = "testdata/out"
)

func TestMainNormal(t *testing.T) {
	tests := []struct {
		desc string
		in   []string
		want int
	}{
		// {
		// 	desc: "正常系: FontIndex, EmojiFontIndexを指定できる",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_font_index.png", "-x", "0", "-X", "0"},
		// 	want: 0,
		// },
		// {
		// 	desc: "正常系: SlackアイコンサイズでアニメーションGIFを生成できる",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_slack_anim.gif", "--slack", "-a"},
		// 	want: 0,
		// },
		// {
		// 	desc: "正常系: 連番付与オプションのテスト",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_numbering.gif", "-n"},
		// 	want: 0,
		// },
		// {
		// 	desc: "正常系: 同じファイルがすでに存在するため別名で保存される (_2)",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_numbering.gif", "-n"},
		// 	want: 0,
		// },
		// {
		// 	desc: "正常系: 同じファイルがすでに存在するため別名で保存される (_3)",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_numbering.gif", "-n"},
		// 	want: 0,
		// },
		// {
		// 	desc: "異常系: 不正な文字色",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "gggg", "-b", "blue"},
		// 	want: -1,
		// },
		// {
		// 	desc: "異常系: 不正な背景色",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_normal_red.png", "-g", "green", "-b", "bbbb"},
		// 	want: -1,
		// },
		// {
		// 	desc: "異常系: 不正なフォントを指定",
		// 	in:   []string{"", "Sample", "-o", outDir + "/illegal_case1.png", "-f", filepath.Join(testdataInputDir, "illegal_font.ttc")},
		// 	want: -1,
		// },
		// {
		// 	desc: "異常系: 不正な絵文字フォントを指定",
		// 	in:   []string{"", "Sample", "-o", outDir + "/illegal_case2.png", "-e", filepath.Join(testdataInputDir, "illegal_font.ttc")},
		// 	want: -1,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			os.Args = tt.in
			got := Main()
			assert.Equal(tt.want, got)
		})
	}

	// _, err := os.Stat(outDir + "/main_numbering_2.gif")
	// assert.NoError(t, err)
	//
	// _, err = os.Stat(outDir + "/main_numbering_3.gif")
	// assert.NoError(t, err)
}
