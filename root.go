package main

import (
	"errors"
	"image/color"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type applicationConfig struct {
	background color.RGBA
	outpath    string
	fontfile   string
	fontsize   int
}

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().StringP("background", "b", "black", `background color.
format is [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringP("out", "o", "", "output image file path")

	font := "/usr/share/fonts/truetype/vlgothic/VL-Gothic-Regular.ttf"
	if runtime.GOOS == "darwin" {
		font = "/Library/Fonts/AppleGothic.ttf"
	}
	RootCommand.Flags().StringP("fontfile", "f", font, "font file path")
	RootCommand.Flags().IntP("fontsize", "F", 64, "font size")
}

var RootCommand = &cobra.Command{
	Use:     "textimg",
	Short:   "textimg is command to convert from colored text (ANSI or 256) to image.",
	Example: `textimg $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flags()

		background, err := f.GetString("background")
		if err != nil {
			panic(err)
		}

		outpath, err := f.GetString("out")
		if err != nil {
			panic(err)
		}

		fontpath, err := f.GetString("fontfile")
		if err != nil {
			panic(err)
		}

		fontsize, err := f.GetInt("fontsize")
		if err != nil {
			panic(err)
		}

		appconf := applicationConfig{
			background: optionBackgrounToRGBA(background),
			outpath:    outpath,
			fontfile:   fontpath,
			fontsize:   fontsize,
		}

		// 引数にテキストの指定がなければ標準入力を使用する
		var texts []string
		if len(args) < 1 {
			texts = readStdin()
		} else {
			texts = args
		}

		// 出力先画像の指定がなければ標準出力を出力先にする
		var w *os.File
		if outpath == "" {
			w = os.Stdout
		} else {
			var err error
			w, err = os.Create(appconf.outpath)
			if err != nil {
				panic(err)
			}
			defer w.Close()
		}

		writeImage(w, texts, appconf)
	},
}

// オプション引数のbackgroundは２つの書き方を許容する。
// 1. black といった色の直接指定
// 2. RGBAのカンマ区切り指定
//    書式: R,G,B,A
//    赤色の例: 255,0,0,255
func optionBackgrounToRGBA(bg string) color.RGBA {
	// "black"といった色名称でマッチするものがあれば返す
	bg = strings.ToLower(bg)
	for k, v := range colorStringMap {
		if bg == k {
			return v
		}
	}

	// カンマ区切りでの指定があれば返す
	rgba := strings.Split(bg, ",")
	if len(rgba) != 4 {
		panic(errors.New("RGBA指定が不正: " + bg))
	}

	var (
		r   uint64
		g   uint64
		b   uint64
		a   uint64
		err error
		rs  = rgba[0]
		gs  = rgba[1]
		bs  = rgba[2]
		as  = rgba[3]
	)
	r, err = strconv.ParseUint(rs, 10, 8)
	if err != nil {
		panic(err)
	}
	g, err = strconv.ParseUint(gs, 10, 8)
	if err != nil {
		panic(err)
	}
	b, err = strconv.ParseUint(bs, 10, 8)
	if err != nil {
		panic(err)
	}
	a, err = strconv.ParseUint(as, 10, 8)
	if err != nil {
		panic(err)
	}
	c := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
	return c
}
