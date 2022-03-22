package main

import (
	"testing"

	"github.com/jiro4989/textimg/v3/config"
	"github.com/stretchr/testify/assert"
)

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
