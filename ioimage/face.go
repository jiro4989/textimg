package ioimage

import (
	"io/ioutil"
	"os"

	"github.com/goki/freetype/truetype"
	"github.com/jiro4989/textimg/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

// ReadFace はfontPathのフォントファイルからfaceを返す。
func ReadFace(fontPath string, fontSize float64) font.Face {
	var fontData []byte

	// ファイルが存在しなければビルトインのフォントをデフォルトとして使う
	_, err := os.Stat(fontPath)
	if err == nil {
		fontData, err = ioutil.ReadFile(fontPath)
		if err != nil {
			panic(err)
		}
	} else {
		log.Warnf("%s is not found. please set font path with `-f` option\n", fontPath)
		fontData = gomono.TTF
	}

	ft, err := truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	opt := truetype.Options{
		Size:              fontSize,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face := truetype.NewFace(ft, &opt)
	return face
}
