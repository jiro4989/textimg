package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jiro4989/textimg/escseq"
	"github.com/jiro4989/textimg/internal/global"
	"github.com/jiro4989/textimg/ioimage"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/image/font"

	"github.com/spf13/cobra"
)

type applicationConfig struct {
	Foreground        escseq.RGBA // 文字色
	Background        escseq.RGBA // 背景色
	Outpath           string      // 画像の出力ファイルパス
	FontFile          string      // フォントファイルのパス
	EmojiFontFile     string      // 絵文字用のフォントファイルのパス
	UseEmojiFont      bool        // 絵文字TTFを使う
	FontSize          int         // フォントサイズ
	UseAnimation      bool        // アニメーションGIFを生成する
	Delay             int         // アニメーションのディレイ時間
	LineCount         int         // 入力データのうち何行を1フレーム画像に使うか
	UseSlideAnimation bool        // スライドアニメーションする
	SlideWidth        int         // スライドする幅
	SlideForever      bool        // スライドを無限にスライドするように描画する
	ToSlackIcon       bool        // Slackのアイコンサイズにする
}

const shellgeiEmojiFontPath = "/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf"

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().SortFlags = false
	RootCommand.Flags().StringP("foreground", "g", "white", `foreground escseq.
format is [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringP("background", "b", "black", `background escseq.
color format is same as "foreground" option`)

	font := "/usr/share/fonts/truetype/vlgothic/VL-Gothic-Regular.ttf"
	if runtime.GOOS == "darwin" {
		font = "/Library/Fonts/AppleGothic.ttf"
	}
	envFontFile := os.Getenv(global.EnvNameFontFile)
	if envFontFile != "" {
		font = envFontFile
	}
	RootCommand.Flags().StringP("fontfile", "f", font, `font file path.
You can change this default value with environment variables TEXTIMG_FONT_FILE`)

	envEmojiFontFile := os.Getenv(global.EnvNameEmojiFontFile)
	RootCommand.Flags().StringP("emoji-fontfile", "e", envEmojiFontFile, "emoji font file")

	RootCommand.Flags().BoolP("use-emoji-font", "i", false, "use emoji font")
	RootCommand.Flags().BoolP("shellgei-emoji-fontfile", "z", false, `emoji font file for shellgei-bot (path: "`+shellgeiEmojiFontPath+`")`)

	RootCommand.Flags().IntP("fontsize", "F", 20, "font size")
	RootCommand.Flags().StringP("out", "o", "", `output image file path.
available image formats are [png | jpg | gif]`)
	RootCommand.Flags().BoolP("shellgei-imagedir", "s", false, `image directory path for shellgei-bot (path: "/images/t.png")`)

	RootCommand.Flags().BoolP("animation", "a", false, "generate animation gif")
	RootCommand.Flags().IntP("delay", "d", 20, "animation delay time")
	RootCommand.Flags().IntP("line-count", "l", 1, "animation input line count")
	RootCommand.Flags().BoolP("slide", "S", false, "use slide animation")
	RootCommand.Flags().IntP("slide-width", "W", 1, "sliding animation width")
	RootCommand.Flags().BoolP("forever", "E", false, "sliding forever")
	RootCommand.Flags().BoolP("environments", "", false, "print environment variables")
	RootCommand.Flags().BoolP("slack", "", false, "resize to slack icon size (128x128 px)")
}

var RootCommand = &cobra.Command{
	Use:     global.AppName,
	Short:   global.AppName + " is command to convert from colored text (ANSI or 256) to image.",
	Example: global.AppName + ` $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: global.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()

		// コマンドライン引数の取得{{{
		printEnv, err := f.GetBool("environments")
		if err != nil {
			return err
		}
		if printEnv {
			for _, envName := range global.EnvNames {
				text := fmt.Sprintf("%s=%s", envName, os.Getenv(envName))
				fmt.Println(text)
			}
			return nil
		}

		foreground, err := f.GetString("foreground")
		if err != nil {
			return err
		}

		background, err := f.GetString("background")
		if err != nil {
			return err
		}

		outpath, err := f.GetString("out")
		if err != nil {
			return err
		}

		useAnimation, err := f.GetBool("animation")
		if err != nil {
			return err
		}

		delay, err := f.GetInt("delay")
		if err != nil {
			return err
		}

		lineCount, err := f.GetInt("line-count")
		if err != nil {
			return err
		}

		useShellGeiDir, err := f.GetBool("shellgei-imagedir")
		if err != nil {
			return err
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
			return err
		}

		emojiFontpath, err := f.GetString("emoji-fontfile")
		if err != nil {
			return err
		}

		useEmojiFont, err := f.GetBool("use-emoji-font")
		if err != nil {
			return err
		}

		useShellGeiEmojiFont, err := f.GetBool("shellgei-emoji-fontfile")
		if err != nil {
			return err
		}
		if useShellGeiEmojiFont {
			emojiFontpath = shellgeiEmojiFontPath
			useEmojiFont = true
		}

		fontsize, err := f.GetInt("fontsize")
		if err != nil {
			return err
		}

		useSlideAnimation, err := f.GetBool("slide")
		if err != nil {
			return err
		}
		if useSlideAnimation {
			useAnimation = true
		}

		slideWidth, err := f.GetInt("slide-width")
		if err != nil {
			return err
		}

		slideForever, err := f.GetBool("forever")
		if err != nil {
			return err
		}

		confForeground, err := optionColorStringToRGBA(foreground)
		if err != nil {
			return err
		}

		confBackground, err := optionColorStringToRGBA(background)
		if err != nil {
			return err
		}

		slack, err := f.GetBool("slack")
		if err != nil {
			return err
		}

		// }}}

		appconf := applicationConfig{
			Foreground:        confForeground,
			Background:        confBackground,
			Outpath:           outpath,
			FontFile:          fontpath,
			EmojiFontFile:     emojiFontpath,
			UseEmojiFont:      useEmojiFont,
			FontSize:          fontsize,
			UseAnimation:      useAnimation,
			Delay:             delay,
			LineCount:         lineCount,
			UseSlideAnimation: useSlideAnimation,
			SlideWidth:        slideWidth,
			SlideForever:      slideForever,
			ToSlackIcon:       slack,
		}

		// 引数にテキストの指定がなければ標準入力を使用する
		var texts []string
		if len(args) < 1 {
			texts = readStdin()
		} else {
			for _, v := range args {
				for _, line := range strings.Split(v, "\n") {
					texts = append(texts, line)
				}
			}
		}

		// textsが空のときは警告メッセージを出力して異常終了
		var emptyCount int
		for _, v := range texts {
			if len(v) < 1 {
				emptyCount++
			}
		}
		if emptyCount == len(texts) {
			err := errors.New("[WARN] Must need input texts.")
			return err
		}

		// スライドアニメーションを使うときはテキストを加工する
		if appconf.UseSlideAnimation {
			texts = toSlideStrings(texts, appconf.LineCount, appconf.SlideWidth, appconf.SlideForever)
		}

		// 拡張子のみ取得
		imgExt := filepath.Ext(strings.ToLower(outpath))

		var w *os.File
		if outpath == "" {
			// 出力先画像の指定がなく、且つ出力先がパイプならstdout + PNG/GIFと
			// して出力。なければそもそも画像処理しても意味が無いので終了
			fd := os.Stdout.Fd()
			if terminal.IsTerminal(int(fd)) {
				fmt.Fprintln(os.Stderr, "textimg: Image data not written to a terminal. Use -o, -s, pipe or redirect.")
				fmt.Fprintln(os.Stderr, "textimg: For help, type: textimg -h")
				return errors.New("No output target error")
			}
			w = os.Stdout
			if useAnimation {
				imgExt = ".gif"
			} else {
				imgExt = ".png"
			}
		} else {
			var err error
			w, err = os.Create(appconf.Outpath)
			if err != nil {
				return err
			}
			defer w.Close()
		}

		// 拡張子は.png, .jpg, .jpeg, .gifのいずれかでなければならない
		switch imgExt {
		case ".png", ".jpg", ".jpeg", ".gif":
			// 何もしない
		default:
			err := errors.New(fmt.Sprintf("%s is not supported extension.", imgExt))
			return err
		}

		// タブ文字は画像描画時に表示されないので暫定対応で半角スペースに置換
		for i, text := range texts {
			texts[i] = strings.Replace(text, "\t", "  ", -1)
		}

		// ゼロ幅文字を削除
		for i, text := range texts {
			texts[i] = removeZeroWidthCharacters(text)
		}

		face := ioimage.ReadFace(appconf.FontFile, float64(appconf.FontSize))
		var emojiFace font.Face
		if appconf.EmojiFontFile != "" {
			emojiFace = ioimage.ReadFace(appconf.EmojiFontFile, float64(appconf.FontSize))
		}
		emojiDir := os.Getenv(global.EnvNameEmojiDir)

		writeConf := ioimage.WriteConfig{
			Foreground:    confForeground,
			Background:    confBackground,
			FontFace:      face,
			EmojiFontFace: emojiFace,
			EmojiDir:      emojiDir,
			UseEmojiFont:  useEmojiFont,
			FontSize:      fontsize,
			UseAnimation:  useAnimation,
			Delay:         delay,
			LineCount:     lineCount,
			ToSlackIcon:   slack,
		}
		if err := ioimage.Write(w, imgExt, texts, writeConf); err != nil {
			return err
		}

		return nil
	},
}

// オプション引数のbackgroundは２つの書き方を許容する。
// 1. black といった色の直接指定
// 2. RGBAのカンマ区切り指定
//    書式: R,G,B,A
//    赤色の例: 255,0,0,255
func optionColorStringToRGBA(colstr string) (escseq.RGBA, error) {
	// "black"といった色名称でマッチするものがあれば返す
	colstr = strings.ToLower(colstr)
	col := escseq.StringMap[colstr]
	zeroColor := escseq.RGBA{}
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
	c := escseq.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
	return c, nil
}

// toSlideStrings は文字列をスライドアニメーション用の文字列に変換する。
func toSlideStrings(src []string, lineCount, slideWidth int, slideForever bool) (ret []string) {
	if 1 < slideWidth {
		var loopCount int
		for i := 0; i < len(src); i += slideWidth {
			loopCount++
		}
		for i := 0; i < (loopCount*slideWidth+1)-len(src); i++ {
			if !slideForever {
				src = append(src, "")
			}
		}
	}

	for i := 0; i < len(src); i += slideWidth {
		n := i + lineCount
		if len(src) < n {
			if slideForever {
				for j := i; j < n; j++ {
					m := j
					if len(src) <= m {
						m -= len(src)
					}
					line := src[m]
					ret = append(ret, line)
				}
				continue
			}
			return
		}
		// lineCountの数ずつ行を取得して戻り値に追加
		for j := i; j < n; j++ {
			line := src[j]
			ret = append(ret, line)
		}
	}
	return
}

// removeZeroWidthSpace はゼロ幅文字が存在したときに削除する。
//
// 参考
// * ゼロ幅スペース https://ja.wikipedia.org/wiki/%E3%82%BC%E3%83%AD%E5%B9%85%E3%82%B9%E3%83%9A%E3%83%BC%E3%82%B9
func removeZeroWidthCharacters(s string) string {
	zwc := []rune{
		0x200b, // zero width space
		0x200c, // zero width joiner
		0x200d, // zero width joiner
		0xfeff, // zero width no-break-space
	}
	var ret []rune
chars:
	for _, v := range []rune(s) {
		for _, c := range zwc {
			if v == c {
				continue chars
			}
		}
		ret = append(ret, v)
	}
	return string(ret)
}
