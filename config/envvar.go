package config

import (
	"fmt"
	"os"
)

type EnvVars struct {
	EmojiDir      string
	OutputDir     string
	FontFile      string
	EmojiFontFile string
}

const (
	envNameOutputDir     = "TEXTIMG_OUTPUT_DIR"
	envNameFontFile      = "TEXTIMG_FONT_FILE"
	envNameEmojiDir      = "TEXTIMG_EMOJI_DIR"
	envNameEmojiFontFile = "TEXTIMG_EMOJI_FONT_FILE"
)

var (
	envs = map[string]string{
		envNameOutputDir:     os.Getenv(envNameOutputDir),
		envNameFontFile:      os.Getenv(envNameFontFile),
		envNameEmojiDir:      os.Getenv(envNameEmojiDir),
		envNameEmojiFontFile: os.Getenv(envNameEmojiFontFile),
	}
)

func NewEnvVars() EnvVars {
	return EnvVars{
		OutputDir:     envs[envNameOutputDir],
		FontFile:      envs[envNameFontFile],
		EmojiDir:      envs[envNameEmojiDir],
		EmojiFontFile: envs[envNameEmojiFontFile],
	}
}

func PrintEnvs() {
	for k, v := range envs {
		text := fmt.Sprintf("%s=%s", k, v)
		fmt.Println(text)
	}
}
