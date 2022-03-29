//go:build docker

//
// 日本語や絵文字が使えるDocker環境上で実行する想定のテスト。
// どうしてもDocker上でしかテストできないもののみこのファイルに記述する。

package main

import (
	"os"
	"testing"

	"github.com/jiro4989/textimg/v3/config"
	"github.com/stretchr/testify/assert"
)

func TestRunRootCommandOnDocker(t *testing.T) {
	var (
		outDockerDir  = "testdata/out_docker"
		fontFile      = "/tmp/MyricaM.TTC"
		emojiDir      = "/usr/local/src/noto-emoji/png/128"
		emojiFontFile = "/tmp/Symbola_hint.ttf"
	)

	// nolint
	os.Mkdir(outDockerDir, os.ModePerm)

	tests := []struct {
		desc       string
		c          config.Config
		args       []string
		envs       config.EnvVars
		wantErr    bool
		existsFile string
	}{
		{
			desc: "正常系: 日本語や絵文字を描画できる",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDockerDir + "/root_on_docker_test_japanese.png"
				c.Writer = nil
				c.FontFile = fontFile
				// c.EmojiDir = emojiDir
				c.EmojiFontFile = emojiFontFile
				return c
			}(),
			args: []string{"\x1b[31mあいうえお\n\x1b[32;43mあ😃a👍！👀ん👄"},
			envs: config.EnvVars{
				EmojiDir: emojiDir,
			},
			wantErr:    false,
			existsFile: outDockerDir + "/root_on_docker_test_japanese.png",
		},
		{
			desc: "正常系: 特殊な絵文字を使う",
			c: func() config.Config {
				c := newDefaultConfig()
				c.Outpath = outDockerDir + "/root_on_docker_test_shellgei_emoji.png"
				c.Writer = nil
				c.UseEmojiFont = true
				c.FontFile = fontFile
				// c.EmojiDir = emojiDir
				c.EmojiFontFile = emojiFontFile
				return c
			}(),
			args:       []string{"\x1b[31mあいうえお\n\x1b[32;43mあ😃a👍！👀ん👄"},
			envs: config.EnvVars{
				EmojiDir: emojiDir,
			},
			wantErr:    false,
			existsFile: outDockerDir + "/root_on_docker_test_shellgei_emoji.png",
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
			if tt.existsFile != "" {
				_, err := os.Stat(tt.existsFile)
				got := os.IsNotExist(err)
				assert.False(got)
			}
		})
	}
}
