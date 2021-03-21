package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jiro4989/textimg/escseq"
	"github.com/jiro4989/textimg/internal/global"
	"github.com/stretchr/testify/assert"
)

func TestOptionColorStringToRGBA(t *testing.T) {
	type TestData struct {
		desc   string
		colstr string
		expect escseq.RGBA
	}
	tds := []TestData{
		{desc: "BLACK", colstr: "BLACK", expect: escseq.RGBABlack},
		{desc: "black", colstr: "black", expect: escseq.RGBABlack},
		{desc: "red", colstr: "red", expect: escseq.RGBARed},
		{desc: "green", colstr: "green", expect: escseq.RGBAGreen},
		{desc: "yellow", colstr: "yellow", expect: escseq.RGBAYellow},
		{desc: "blue", colstr: "blue", expect: escseq.RGBABlue},
		{desc: "magenta", colstr: "magenta", expect: escseq.RGBAMagenta},
		{desc: "cyan", colstr: "cyan", expect: escseq.RGBACyan},
		{desc: "white", colstr: "white", expect: escseq.RGBAWhite},
		{desc: "0,0,0,255", colstr: "0,0,0,255", expect: escseq.RGBA{R: 0, G: 0, B: 0, A: 255}},
		{desc: "255,255,255,255", colstr: "255,255,255,255", expect: escseq.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{desc: "0,0,0,0", colstr: "0,0,0,0", expect: escseq.RGBA{R: 0, G: 0, B: 0, A: 0}},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got, err := optionColorStringToRGBA(v.colstr)
			assert.Nil(t, err, v.desc)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}

	// 異常系
	tds = []TestData{
		{desc: "不正な色文字列", colstr: "unko"},
		{desc: "RGBAの書式不正(値の数不足)", colstr: "1,2,3"},
		{desc: "RGBAの書式不正(値の数過多)", colstr: "1,2,3,4,5"},
		{desc: "RGBAの書式不正(値がない)", colstr: "1,2,3,"},
		{desc: "RGBAの書式不正(値に文字が混じっている)", colstr: "1,2,3,a"},
		{desc: "RGBAの書式不正(255以上の値)", colstr: "1,2,3,256"},
		{desc: "RGBAの書式不正(負の値)", colstr: "-1,2,3,255"},
		{desc: "RGBAの書式不正(空文字)", colstr: ""},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			_, err := optionColorStringToRGBA(v.colstr)
			assert.NotNil(t, err, v.desc)
		})
	}
}

func TestToSlideStrings(t *testing.T) {
	type TestData struct {
		desc                  string
		src, expect           []string
		lineCount, slideWidth int
		slideForever          bool
	}
	tds := []TestData{
		{
			desc: "2行描画、スライド幅1、無限なし",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2",
				"2", "3",
				"3", "4",
				"4", "5",
			},
			lineCount:    2,
			slideWidth:   1,
			slideForever: false,
		},
		{
			desc: "2行描画、スライド幅2、無限なし",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2",
				"3", "4",
				"5", "",
			},
			lineCount:    2,
			slideWidth:   2,
			slideForever: false,
		},
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
			desc: "3行描画、スライド幅3、無限なし、不足なし",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅3、無限なし、不足あり",
			src:  []string{"1", "2", "3", "4", "5", "6", "7"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
				"7", "", "",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: false,
		},
		{
			desc: "3行描画、スライド幅3、無限なし、不足あり",
			src:  []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
				"7", "8", "",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: false,
		},
		{
			desc: "2行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2",
				"3", "4",
				"5", "1",
			},
			lineCount:    2,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "2行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2",
				"3", "4",
				"5", "6",
			},
			lineCount:    2,
			slideWidth:   2,
			slideForever: true,
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
			},
			lineCount:    3,
			slideWidth:   1,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅1、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"2", "3", "4",
				"3", "4", "5",
				"4", "5", "6",
				"5", "6", "1",
				"6", "1", "2",
			},
			lineCount:    3,
			slideWidth:   1,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5"},
			expect: []string{
				"1", "2", "3",
				"3", "4", "5",
				"5", "1", "2",
			},
			lineCount:    3,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅2、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"3", "4", "5",
				"5", "6", "1",
			},
			lineCount:    3,
			slideWidth:   2,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅3、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: true,
		},
		{
			desc: "3行描画、スライド幅3、無限あり",
			src:  []string{"1", "2", "3", "4", "5", "6", "7"},
			expect: []string{
				"1", "2", "3",
				"4", "5", "6",
				"7", "1", "2",
			},
			lineCount:    3,
			slideWidth:   3,
			slideForever: true,
		},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := toSlideStrings(v.src, v.lineCount, v.slideWidth, v.slideForever)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}

func TestRemoveZeroWidthCharacters(t *testing.T) {
	type TestData struct {
		desc   string
		s      string
		expect string
	}
	tds := []TestData{
		{desc: "Zero width space (U+200B)が削除される", s: "A\u200bB", expect: "AB"},
		{desc: "Zero width joiner (U+200C)が削除される", s: "A\u200cB", expect: "AB"},
		{desc: "Zero width joiner (U+200D)が削除される", s: "A\u200dB", expect: "AB"},
		{desc: "U+200B ~ U+200Dが削除される", s: "あ\u200bい\u200cう\u200dえ", expect: "あいうえ"},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			got := removeZeroWidthCharacters(v.s)
			assert.Equal(t, v.expect, got, v.desc)
		})
	}
}

func TestApplicationConfigSetFontFileAndFontIndex(t *testing.T) {
	type TestData struct {
		desc          string
		inFontFile    string
		inFontIndex   int
		inRuntimeOS   string
		wantFontFile  string
		wantFontIndex int
	}
	tests := []TestData{
		{
			desc:          "正常系: FontFileが設定済みの場合は変更なし",
			inFontFile:    "/usr/share/fonts/寿司",
			inFontIndex:   0,
			inRuntimeOS:   "linux",
			wantFontFile:  "/usr/share/fonts/寿司",
			wantFontIndex: 0,
		},
		{
			desc:          "正常系: フォント未設定でwindowsの場合はwindows用のフォントが設定される",
			inRuntimeOS:   "windows",
			wantFontFile:  defaultWindowsFont,
			wantFontIndex: 0,
		},
		{
			desc:          "正常系: フォント未設定でdarwinの場合はdarwin用のフォントが設定される",
			inRuntimeOS:   "darwin",
			wantFontFile:  defaultDarwinFont,
			wantFontIndex: 0,
		},
		{
			desc:          "正常系: フォント未設定でiosの場合はios用のフォントが設定される",
			inRuntimeOS:   "ios",
			wantFontFile:  defaultIOSFont,
			wantFontIndex: 0,
		},
		{
			desc:          "正常系: フォント未設定でandroidの場合はandroid用のフォントが設定される",
			inRuntimeOS:   "android",
			wantFontFile:  defaultAndroidFont,
			wantFontIndex: 4,
		},
		{
			desc:          "正常系: フォント未設定でlinuxの場合はlinux用のフォントが設定される。Linux用のフォントは2つ存在するが、1つ目のフォントはalpineコンテナ内にデフォルトでは存在しないため2つ目が設定される",
			inRuntimeOS:   "linux",
			wantFontFile:  defaultLinuxFont2,
			wantFontIndex: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			a := applicationConfig{
				FontFile:  tt.inFontFile,
				FontIndex: tt.inFontIndex,
			}
			a.setFontFileAndFontIndex(tt.inRuntimeOS)

			assert.Equal(tt.wantFontFile, a.FontFile)
			assert.Equal(tt.wantFontIndex, a.FontIndex)
		})
	}
}

func TestApplicationConfig_AddTimeStampToOutPath(t *testing.T) {
	type TestData struct {
		desc           string
		inOutpath      string
		inAddTimeStamp bool
		inTime         time.Time
		want           string
	}
	tests := []TestData{
		{
			desc:           "正常系: フラグfalseの場合は変更なし",
			inOutpath:      "t.png",
			inAddTimeStamp: false,
			inTime:         time.Date(2000, 1, 1, 12, 10, 5, 2, time.Local),
			want:           "t.png",
		},
		{
			desc:           "正常系: フラグtrueの場合はタイムスタンプがつく",
			inOutpath:      "t.png",
			inAddTimeStamp: true,
			inTime:         time.Date(2000, 1, 1, 12, 10, 5, 0, time.Local),
			want:           "t_2000-01-01-121005.png",
		},
		{
			desc:           "正常系: フルパスでも同様に動作する",
			inOutpath:      "/images/t.png",
			inAddTimeStamp: true,
			inTime:         time.Date(2000, 1, 1, 12, 10, 5, 0, time.Local),
			want:           "/images/t_2000-01-01-121005.png",
		},
		{
			desc:           "正常系: ファイル拡張子が多重についていても動作する",
			inOutpath:      "/images/t.png.1",
			inAddTimeStamp: true,
			inTime:         time.Date(2000, 1, 1, 12, 10, 5, 0, time.Local),
			want:           "/images/t.png_2000-01-01-121005.1",
		},
		{
			desc:           "正常系: Windowsのパス表現でも動作する",
			inOutpath:      `C:\Users\foobar\Pictures\t.png`,
			inAddTimeStamp: true,
			inTime:         time.Date(2000, 1, 1, 12, 10, 5, 0, time.Local),
			want:           `C:\Users\foobar\Pictures\t_2000-01-01-121005.png`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			a := applicationConfig{
				Outpath:      tt.inOutpath,
				AddTimeStamp: tt.inAddTimeStamp,
			}
			a.addTimeStampToOutPath(tt.inTime)

			assert.Equal(tt.want, a.Outpath)
		})
	}
}

func TestOutputImageDir(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.NoError(t, err)
	pictDir := filepath.Join(home, "Pictures")

	type TestData struct {
		desc           string
		inEnvDir       string
		inUseAnimation bool
		wantPath       string
		wantErr        bool
	}
	tests := []TestData{
		{
			desc:           "正常系: Env未設定の場合はホームディレクトリ配下のPictures配下が返る",
			inEnvDir:       "",
			inUseAnimation: false,
			wantPath:       filepath.Join(pictDir, "t.png"),
			wantErr:        false,
		},
		{
			desc:           "正常系: animation trueの場合は Basenameが t.gif になる",
			inEnvDir:       "",
			inUseAnimation: true,
			wantPath:       filepath.Join(pictDir, "t.gif"),
			wantErr:        false,
		},
		{
			desc:           "正常系: Envが設定されている場合は設定されている値が優先される",
			inEnvDir:       filepath.Join(".", "sushi"),
			inUseAnimation: false,
			wantPath:       filepath.Join(".", "sushi", "t.png"),
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			os.Setenv(global.EnvNameOutputDir, tt.inEnvDir)

			got, err := outputImageDir(tt.inUseAnimation)
			if tt.wantErr {
				assert.Equal("", got)
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(tt.wantPath, got)
		})
	}
}
