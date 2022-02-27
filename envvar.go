package main

import (
	"os"

	"github.com/jiro4989/textimg/v3/internal/global"
)

type EnvVars struct {
	EmojiDir      string
	OutputDir     string
	FontFile      string
	EmojiFontFile string
}

func NewEnvVars() EnvVars {
	return EnvVars{
		EmojiDir:      os.Getenv(global.EnvNameEmojiDir),
		OutputDir:     os.Getenv(global.EnvNameOutputDir),
		FontFile:      os.Getenv(global.EnvNameFontFile),
		EmojiFontFile: os.Getenv(global.EnvNameEmojiFontFile),
	}
}
