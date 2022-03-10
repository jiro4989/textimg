package config

import (
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
			desc: "正常系: UseShellgeiImagedirが有効な場合はOutpathを上書きする",
			config: func() Config {
				c := newDefaultConfig()
				c.Outpath = "t.png"
				return c
			}(),
			args: []string{"hello"},
			ev:   EnvVars{},
			want: func() Config {
				c := newDefaultConfig()
				c.ForegroundColor = color.RGBABlack
				c.BackgroundColor = color.RGBAWhite
				c.Texts = []string{"hello"}
				c.FileExtension = ".png"
				return c
			}(),
			wantErr: false,
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
			assert.Equal(tt.want, tt.config)
		})
	}
}
