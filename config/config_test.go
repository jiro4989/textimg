package config

import (
	"path/filepath"
	"testing"

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
