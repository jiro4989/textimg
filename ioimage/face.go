package ioimage

import (
	"io/ioutil"
	"os"

	"github.com/jiro4989/textimg/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
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
		log.Warnf("%s is not found. please set font path with `-f` option", fontPath)
		fontData = gomono.TTF
	}

	ft, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	opt := opentype.FaceOptions{
		Size:              fontSize,
		DPI:               72,
		Hinting:           0,
	}
	face, err := opentype.NewFace(ft, &opt)
	if err != nil {
		panic(err)
	}
	return face
}
