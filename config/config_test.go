package config

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/stretchr/testify/assert"
)

func newDefaultConfig() Config {
	return Config{
		Foreground:               "white",
		Background:               "black",
		Outpath:                  "",
		AddTimeStamp:             false,
		SaveNumberedFile:         false,
		FontFile:                 "",
		FontIndex:                0,
		EmojiFontFile:            "",
		EmojiFontIndex:           0,
		UseEmojiFont:             false,
		FontSize:                 20,
		UseAnimation:             false,
		Delay:                    20,
		LineCount:                1,
		UseSlideAnimation:        false,
		SlideWidth:               1,
		SlideForever:             false,
		ToSlackIcon:              false,
		PrintEnvironments:        false,
		UseShellgeiImagedir:      false,
		UseShellgeiEmojiFontfile: false,
		ResizeWidth:              0,
		ResizeHeight:             0,
		Writer:                   NewMockWriter(false, false),
	}
}

func TestConfig_Adjust(t *testing.T) {
	tests := []struct {
		desc    string
		config  Config
		args    []string
		ev      EnvVars
		want    Config
		wantErr bool
	}{
		{
			desc: "正常系: Outpathが設定されている",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				return c
			}(),
			args: []string{"hello"},
			ev:   EnvVars{},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: UseShellgeiImagedirが有効なときはt.pngが設定される",
			config: func() Config {
				c := newDefaultConfig()
				c.UseShellgeiImagedir = true
				return c
			}(),
			args: []string{"hello"},
			ev: EnvVars{
				OutputDir: "sushi",
			},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = filepath.Join("sushi", "t.png")
				c.UseShellgeiImagedir = true
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: UseShellgeiEmojiFontfileが有効な時は組み込みの絵文字パスが設定されて、UseEmojiFont=trueになる",
			config: func() Config {
				c := newDefaultConfig()
				c.UseShellgeiImagedir = true
				c.UseShellgeiEmojiFontfile = true
				return c
			}(),
			args: []string{"hello"},
			ev: EnvVars{
				OutputDir: "sushi",
			},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = filepath.Join("sushi", "t.png")
				c.UseShellgeiImagedir = true
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				c.UseShellgeiEmojiFontfile = true
				c.UseEmojiFont = true
				c.EmojiFontFile = ShellgeiEmojiFontPath
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: UseShellgeiImagedirが有効でUseAnimationが設定されているときはt.gifになる",
			config: func() Config {
				c := newDefaultConfig()
				c.UseShellgeiImagedir = true
				c.UseAnimation = true
				return c
			}(),
			args: []string{"hello"},
			ev: EnvVars{
				OutputDir: "sushi",
			},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = filepath.Join("sushi", "t.gif")
				c.UseShellgeiImagedir = true
				c.UseAnimation = true
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".gif"
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: ToSlackIconが有効なときはResizeWidthとResizeHeightが設定される",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.ToSlackIcon = true
				return c
			}(),
			args: []string{"hello"},
			ev:   EnvVars{},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				c.ToSlackIcon = true
				c.ResizeWidth = 128
				c.ResizeHeight = 128
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: UseSlideAnimationが有効なときはUseAnimationも有効になる",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.UseSlideAnimation = true
				return c
			}(),
			args: []string{"hello"},
			ev:   EnvVars{},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				c.UseSlideAnimation = true
				c.UseAnimation = true
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: SlideWidthが2以上の時はテキストの処理が変化する",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.UseSlideAnimation = true
				c.LineCount = 2
				c.SlideWidth = 2
				c.SlideForever = true
				return c
			}(),
			args: []string{"hello", "hello", "hello", "hello"},
			ev:   EnvVars{},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello", "hello", "hello", "hello"}
				c.FileExtension = ".png"
				c.UseSlideAnimation = true
				c.UseAnimation = true
				c.LineCount = 2
				c.SlideWidth = 2
				c.SlideForever = true
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "正常系: EmojiFontFileに存在しないファイルを指定してもエラーにはならない",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.EmojiFontFile = "sushi.otf"
				return c
			}(),
			args: []string{"hello"},
			ev:   EnvVars{},
			want: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.ForegroundColor = color.RGBAWhite
				c.BackgroundColor = color.RGBABlack
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				c.EmojiFontFile = "sushi.otf"
				return c
			}(),
			wantErr: false,
		},
		{
			desc: "異常系: Outpathが空のときで出力先が存在しないときはエラーが返る",
			config: func() Config {
				c := newDefaultConfig()
				return c
			}(),
			args:    []string{"hello"},
			ev:      EnvVars{},
			want:    Config{},
			wantErr: true,
		},
		{
			desc: "異常系: Foregroundに不正な色指定をした時はエラーを返す",
			config: func() Config {
				c := newDefaultConfig()
				c.Foreground = "sushi"
				return c
			}(),
			args:    []string{"hello"},
			ev:      EnvVars{},
			want:    Config{},
			wantErr: true,
		},
		{
			desc: "異常系: Backgroundに不正な色指定をした時はエラーを返す",
			config: func() Config {
				c := newDefaultConfig()
				c.Background = "sushi"
				return c
			}(),
			args:    []string{"hello"},
			ev:      EnvVars{},
			want:    Config{},
			wantErr: true,
		},
		{
			desc: "異常系: textsが空の時はエラーを返す",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				return c
			}(),
			args:    []string{},
			ev:      EnvVars{},
			want:    Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			err := tt.config.Adjust(tt.args, tt.ev)
			if tt.wantErr {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(tt.want.Foreground, tt.config.Foreground)
			assert.Equal(tt.want.Background, tt.config.Background)
			assert.Equal(tt.want.Outpath, tt.config.Outpath)
			assert.Equal(tt.want.AddTimeStamp, tt.config.AddTimeStamp)
			assert.Equal(tt.want.SaveNumberedFile, tt.config.SaveNumberedFile)
			assert.Equal(tt.want.FontFile, tt.config.FontFile)
			assert.Equal(tt.want.FontIndex, tt.config.FontIndex)
			assert.Equal(tt.want.EmojiFontFile, tt.config.EmojiFontFile)
			assert.Equal(tt.want.EmojiFontIndex, tt.config.EmojiFontIndex)
			assert.Equal(tt.want.UseEmojiFont, tt.config.UseEmojiFont)
			assert.Equal(tt.want.FontSize, tt.config.FontSize)
			assert.Equal(tt.want.UseAnimation, tt.config.UseAnimation)
			assert.Equal(tt.want.Delay, tt.config.Delay)
			assert.Equal(tt.want.LineCount, tt.config.LineCount)
			assert.Equal(tt.want.UseSlideAnimation, tt.config.UseSlideAnimation)
			assert.Equal(tt.want.SlideWidth, tt.config.SlideWidth)
			assert.Equal(tt.want.SlideForever, tt.config.SlideForever)
			assert.Equal(tt.want.ToSlackIcon, tt.config.ToSlackIcon)
			assert.Equal(tt.want.PrintEnvironments, tt.config.PrintEnvironments)
			assert.Equal(tt.want.UseShellgeiImagedir, tt.config.UseShellgeiImagedir)
			assert.Equal(tt.want.UseShellgeiEmojiFontfile, tt.config.UseShellgeiEmojiFontfile)
			assert.Equal(tt.want.ForegroundColor, tt.config.ForegroundColor)
			assert.Equal(tt.want.BackgroundColor, tt.config.BackgroundColor)
			assert.Equal(tt.want.Texts, tt.config.Texts)
			assert.Equal(tt.want.FileExtension, tt.config.FileExtension)
			// NOTE: ここはテストするのが難しいので無視
			// assert.Equal(tt.want.Writer, tt.config.Writer)
			// assert.Equal(tt.want.FontFace, tt.config.FontFace)
			// assert.Equal(tt.want.EmojiFontFace, tt.config.EmojiFontFace)
			assert.Equal(tt.want.EmojiDir, tt.config.EmojiDir)
		})
	}
}

func TestOptionColorStringToRGBA(t *testing.T) {
	type TestData struct {
		desc   string
		colstr string
		expect color.RGBA
	}
	tds := []TestData{
		{desc: "BLACK", colstr: "BLACK", expect: color.RGBABlack},
		{desc: "black", colstr: "black", expect: color.RGBABlack},
		{desc: "red", colstr: "red", expect: color.RGBARed},
		{desc: "green", colstr: "green", expect: color.RGBAGreen},
		{desc: "yellow", colstr: "yellow", expect: color.RGBAYellow},
		{desc: "blue", colstr: "blue", expect: color.RGBABlue},
		{desc: "magenta", colstr: "magenta", expect: color.RGBAMagenta},
		{desc: "cyan", colstr: "cyan", expect: color.RGBACyan},
		{desc: "white", colstr: "white", expect: color.RGBAWhite},
		{desc: "0,0,0,255", colstr: "0,0,0,255", expect: color.RGBA{R: 0, G: 0, B: 0, A: 255}},
		{desc: "255,255,255,255", colstr: "255,255,255,255", expect: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{desc: "0,0,0,0", colstr: "0,0,0,0", expect: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
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
			wantFontIndex: 5,
		},
		// FIXME: ローカル環境で実行するとエラーになるので一旦無効化
		// {
		// 	desc:          "正常系: フォント未設定でlinuxの場合はlinux用のフォントが設定される。Linux用のフォントは2つ存在するが、1つ目のフォントはalpineコンテナ内にデフォルトでは存在しないため2つ目が設定される",
		// 	inRuntimeOS:   "linux",
		// 	wantFontFile:  defaultLinuxFont2,
		// 	wantFontIndex: 5,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			a := Config{
				FontFile:  tt.inFontFile,
				FontIndex: tt.inFontIndex,
			}
			a.SetFontFileAndFontIndex(tt.inRuntimeOS)

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

			a := Config{
				Outpath:      tt.inOutpath,
				AddTimeStamp: tt.inAddTimeStamp,
			}
			a.addTimeStampToOutPath(tt.inTime)

			assert.Equal(tt.want, a.Outpath)
		})
	}
}
