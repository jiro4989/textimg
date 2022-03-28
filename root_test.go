package main

import (
	"os"
	"testing"

	"github.com/jiro4989/textimg/v3/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testBefore()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func testBefore() {
	if err := os.RemoveAll(outDir); err != nil {
		panic(err)
	}
	os.Mkdir(outDir, os.ModePerm)
}

func newDefaultConfig() config.Config {
	return config.Config{
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
		Writer:                   config.NewMockWriter(false, false),
	}
}

func TestRunRootCommand(t *testing.T) {
	tests := []struct {
		desc    string
		c       config.Config
		args    []string
		envs    config.EnvVars
		wantErr bool
	}{
		{
			desc: "正常系: PrintEnvironmentsが設定されていると環境変数を出力して終了",
			c: config.Config{
				PrintEnvironments: true,
			},
			args:    []string{"hello"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: 正常系がパスする",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				return c
			}(),
			args:    []string{"hello"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "異常系: Writerがエラーを返す",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				c.Writer = config.NewMockWriter(true, false)
				return c
			}(),
			args:    []string{"hello"},
			envs:    config.EnvVars{},
			wantErr: true,
		},
		// 旧 main_test.go を移行してきたもの
		{
			desc: "正常系: 画像ファイルに出力する",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_red_and_background_is_black.png"
				c.Writer = nil
				return c
			}(),
			args:    []string{"\x1b[31mred\x1b[m"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: 文字色と背景色を変更する",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_green_and_background_is_blue.png"
				c.Writer = nil
				c.Foreground = "green"
				c.Background = "blue"
				return c
			}(),
			args:    []string{"green"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: カンマ区切りで指定",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_blue_and_background_is_red.png"
				c.Writer = nil
				c.Foreground = "0,0,255,255"
				c.Background = "255,0,0,255"
				return c
			}(),
			args:    []string{"blue"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: Slackアイコンサイズで生成する",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_blue_and_background_is_red_slack_icon_size.png"
				c.Writer = nil
				c.Foreground = "0,0,255,255"
				c.Background = "255,0,0,255"
				c.ToSlackIcon = true
				return c
			}(),
			args:    []string{"slack"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: 明示的に幅を指定できる",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_blue_and_background_is_red_100x200.png"
				c.Writer = nil
				c.Foreground = "0,0,255,255"
				c.Background = "255,0,0,255"
				c.ResizeWidth = 100
				c.ResizeHeight = 200
				return c
			}(),
			args:    []string{"100x200"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: Widthのみを指定した場合はHeightが調整される",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_blue_and_background_is_red_100w.png"
				c.Writer = nil
				c.Foreground = "0,0,255,255"
				c.Background = "255,0,0,255"
				c.ResizeWidth = 100
				return c
			}(),
			args:    []string{"100w"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
		{
			desc: "正常系: Heightのみを指定した場合はWidthが調整される",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDir + "/root_test_font_is_blue_and_background_is_red_100h.png"
				c.Writer = nil
				c.Foreground = "0,0,255,255"
				c.Background = "255,0,0,255"
				c.ResizeHeight = 100
				return c
			}(),
			args:    []string{"100h"},
			envs:    config.EnvVars{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			err := RunRootCommand(tt.c, tt.args, tt.envs)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			assert.NoError(err)
		})
	}
}

func TestComplementWidthHeight(t *testing.T) {
	type TestData struct {
		desc       string
		x, y, w, h int
		wantWidth  int
		wantHeight int
	}
	tds := []TestData{
		{
			desc:       "正常系: wが0のときはwidthが調整される",
			x:          200,
			y:          100,
			w:          0,
			h:          200,
			wantWidth:  400,
			wantHeight: 200,
		},
		{
			desc:       "正常系: hが0のときはheightが調整される",
			x:          200,
			y:          100,
			w:          100,
			h:          0,
			wantWidth:  100,
			wantHeight: 50,
		},
		{
			desc:       "正常系: hが0のときはheightが調整される",
			x:          200,
			y:          100,
			w:          100,
			h:          0,
			wantWidth:  100,
			wantHeight: 50,
		},
		{
			desc:       "正常系: wとhが0出ないときはwとhが返る",
			x:          200,
			y:          100,
			w:          400,
			h:          300,
			wantWidth:  400,
			wantHeight: 300,
		},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			a := assert.New(t)
			w, h := complementWidthHeight(v.x, v.y, v.w, v.h)
			a.Equal(v.wantWidth, w)
			a.Equal(v.wantHeight, h)
		})
	}
}
