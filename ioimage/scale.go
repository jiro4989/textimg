package ioimage

import (
	"image"

	"golang.org/x/image/draw"
)

// scaleToSlackIconSize はslackがtrueの時だけ画像をSlackアイコンサイズ(128x128)
// に拡縮する。
func scaleToSlackIconSize(img *image.RGBA, slack bool) *image.RGBA {
	if !slack {
		return img
	}
	w, h := 128, 128
	rect := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, rect, draw.Over, nil)
	return dst
}
