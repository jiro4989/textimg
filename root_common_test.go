package main

import "github.com/jiro4989/textimg/v3/config"

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
