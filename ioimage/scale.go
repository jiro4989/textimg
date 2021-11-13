package ioimage

import (
	"image"

	"golang.org/x/image/draw"
)

type scaleOption struct {
	w, h        int
	toSlackIcon bool
}

// scale は w, h の幅に画像を拡縮する。
func scale(img *image.RGBA, opt scaleOption) *image.RGBA {
	rect := img.Bounds()
	w, h := opt.w, opt.h
	if opt.toSlackIcon {
		w, h = 128, 128
	}
	size := rect.Size()
	w, h = complementWidthHeight(size.X, size.Y, w, h)
	if w == -1 && h == -1 {
		return img
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, rect, draw.Over, nil)
	return dst
}

func complementWidthHeight(x, y, w, h int) (int, int) {
	if w == -1 {
		hh := y
		d := float64(h) / float64(hh)
		w = int(float64(x) * d)
		return w, h
	}
	if h == -1 {
		ww := x
		d := float64(w) / float64(ww)
		h = int(float64(y) * d)
		return w, h
	}
	return x, y
}
