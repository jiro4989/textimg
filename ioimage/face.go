package ioimage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiro4989/textimg/v3/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

// ReadFace はfontPathのフォントファイルからfaceを返す。
func ReadFace(fontPath string, fontIndex int, fontSize float64) (font.Face, error) {
	var ft *opentype.Font

	// ファイルが存在しなければビルトインのフォントをデフォルトとして使う
	_, err := os.Stat(fontPath)
	if err == nil {
		fontData, err := ioutil.ReadFile(fontPath)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(filepath.Ext(fontPath)) {
		case ".otc", ".ttc":
			ftc, err := opentype.ParseCollection(fontData)
			if err != nil {
				return nil, err
			}
			ft, err = ftc.Font(fontIndex)
			if err != nil {
				return nil, err
			}
		default:
			ft, err = opentype.Parse(fontData)
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Warnf("%s is not found. please set font path with `-f` option", fontPath)
		ft, err = opentype.Parse(gomono.TTF)
		if err != nil {
			return nil, err
		}
	}

	opt := opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: 0,
	}
	face, err := opentype.NewFace(ft, &opt)
	if err != nil {
		return nil, err
	}
	return face, nil
}
