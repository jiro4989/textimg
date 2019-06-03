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
	foreground        color.RGBA // 文字色
	background        color.RGBA // 背景色
	outpath           string     // 画像の出力ファイルパス
	fontfile          string     // フォントファイルのパス
	fontsize          int        // フォントサイズ
	useAnimation      bool       // アニメーションGIFを生成する
	delay             int        // アニメーションのディレイ時間
	lineCount         int        // 入力データのうち何行を1フレーム画像に使うか
	useSlideAnimation bool       // スライドアニメーションする
	slideWidth        int        // スライドする幅
	slideForever      bool       // スライドを無限にスライドするように描画する
}

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().SortFlags = false
	RootCommand.Flags().StringP("foreground", "", "white", `foreground color.
format is [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringP("background", "b", "black", `ackground color.
color format is same as "foreground" option`)
	font := "/usr/share/fonts/truetype/vlgothic/VL-Gothic-Regular.ttf"
	if runtime.GOOS == "darwin" {
		font = "/Library/Fonts/AppleGothic.ttf"
	}
	envFontFile := os.Getenv("TEXTIMG_FONT_FILE")
	if envFontFile != "" {
		font = envFontFile
	}
	RootCommand.Flags().StringP("fontfile", "f", font, `font file path.
You can change this default value with environment variables TEXTIMG_FONT_FILE`)
	RootCommand.Flags().IntP("fontsize", "F", 20, "font size")
	RootCommand.Flags().StringP("out", "o", "", `output image file path.
available image formats are [png | jpg | gif]`)
	RootCommand.Flags().BoolP("shellgei-imagedir", "s", false, `image directory path for shell gei bot (path: "/images/t.png")`)

	RootCommand.Flags().BoolP("animation", "a", false, "generate animation gif")
	RootCommand.Flags().IntP("delay", "d", 20, "animation delay time")
	RootCommand.Flags().IntP("line-count", "l", 1, "animation input line count")
	RootCommand.Flags().BoolP("slide", "S", false, "use slide animation")
	RootCommand.Flags().IntP("slide-width", "W", 1, "sliding animation width")
	RootCommand.Flags().BoolP("forever", "E", false, "sliding forever")
}

var RootCommand = &cobra.Command{
	Use:     "textimg",
	Short:   "textimg is command to convert from colored text (ANSI or 256) to image.",
	Example: `textimg $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flags()

		// コマンドライン引数の取得{{{
		foreground, err := f.GetString("foreground")
		if err != nil {
			panic(err)
		}

		background, err := f.GetString("background")
		if err != nil {
			panic(err)
		}

		outpath, err := f.GetString("out")
		if err != nil {
			panic(err)
		}

		useAnimation, err := f.GetBool("animation")
		if err != nil {
			panic(err)
		}

		delay, err := f.GetInt("delay")
		if err != nil {
			panic(err)
		}

		lineCount, err := f.GetInt("line-count")
		if err != nil {
			panic(err)
		}

		useShellGeiDir, err := f.GetBool("shellgei-imagedir")
		if err != nil {
			panic(err)
		}
		if useShellGeiDir {
			if useAnimation {
				outpath = "/images/t.gif"
			} else {
				outpath = "/images/t.png"
			}
		}

		fontpath, err := f.GetString("fontfile")
		if err != nil {
			panic(err)
		}

		fontsize, err := f.GetInt("fontsize")
		if err != nil {
			panic(err)
		}

		useSlideAnimation, err := f.GetBool("slide")
		if err != nil {
			panic(err)
		}
		if useSlideAnimation {
			useAnimation = true
		}

		slideWidth, err := f.GetInt("slide-width")
		if err != nil {
			panic(err)
		}

		slideForever, err := f.GetBool("forever")
		if err != nil {
			panic(err)
		}

		confForeground, err := optionColorStringToRGBA(foreground)
		if err != nil {
			panic(err)
		}

		confBackground, err := optionColorStringToRGBA(background)
		if err != nil {
			panic(err)
		}

		// }}}

		appconf := applicationConfig{
			foreground:        confForeground,
			background:        confBackground,
			outpath:           outpath,
			fontfile:          fontpath,
			fontsize:          fontsize,
			useAnimation:      useAnimation,
			delay:             delay,
			lineCount:         lineCount,
			useSlideAnimation: useSlideAnimation,
			slideWidth:        slideWidth,
			slideForever:      slideForever,
		}

		// 引数にテキストの指定がなければ標準入力を使用する
		var texts []string
		if len(args) < 1 {
			texts = readStdin()
		} else {
			texts = args
		}

		// スライドアニメーションを使うときはテキストを加工する
		if appconf.useSlideAnimation {
			texts = toSlideStrings(texts, appconf.lineCount, appconf.slideWidth, appconf.slideForever)
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

		encFmt, err := getEncodeFormat(outpath)

		writeImage(w, encFmt, texts, appconf)
	},
}

// オプション引数のbackgroundは２つの書き方を許容する。
// 1. black といった色の直接指定
// 2. RGBAのカンマ区切り指定
//    書式: R,G,B,A
//    赤色の例: 255,0,0,255
func optionColorStringToRGBA(colstr string) (color.RGBA, error) {
	// "black"といった色名称でマッチするものがあれば返す
	colstr = strings.ToLower(colstr)
	col := colorStringMap[colstr]
	zeroColor := color.RGBA{}
	if col != zeroColor {
		return col, nil
	}

	// カンマ区切りでの指定があれば返す
	rgba := strings.Split(colstr, ",")
	if len(rgba) != 4 {
		return zeroColor, errors.New("RGBA指定が不正: " + colstr)
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
		return zeroColor, err
	}
	g, err = strconv.ParseUint(gs, 10, 8)
	if err != nil {
		return zeroColor, err
	}
	b, err = strconv.ParseUint(bs, 10, 8)
	if err != nil {
		return zeroColor, err
	}
	a, err = strconv.ParseUint(as, 10, 8)
	if err != nil {
		return zeroColor, err
	}
	c := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
	return c, nil
}
