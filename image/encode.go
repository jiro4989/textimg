package image

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func (i *Image) Encode(w io.Writer, ext string) error {
	img := i.image
	switch ext {
	case ".png":
		return png.Encode(w, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(w, img, nil)
	case ".gif":
		if i.useAnimation {
			var delays []int
			for x := 0; x < len(i.animationImages); x++ {
				delays = append(delays, i.delay)
			}
			return gif.EncodeAll(w, &gif.GIF{
				Image: toPalettes(i.animationImages),
				Delay: delays,
			})
		}
		return gif.Encode(w, img, nil)
	}
	return fmt.Errorf("%s is not supported extension.", ext)
}

func toPalettes(imgs []image.Image) (ret []*image.Paletted) {
	for _, v := range imgs {
		bounds := v.Bounds()
		p := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(p, p.Rect, v, bounds.Min, draw.Over)
		ret = append(ret, p)
	}
	return
}
