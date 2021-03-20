package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jiro4989/textimg/escseq"
	"github.com/jiro4989/textimg/internal/global"
	"github.com/jiro4989/textimg/ioimage"
	"github.com/jiro4989/textimg/log"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/image/font"

	"github.com/spf13/cobra"
)

type applicationConfig struct {
	Foreground               string // 文字色
	Background               string // 背景色
	Outpath                  string // 画像の出力ファイルパス
	FontFile                 string // フォントファイルのパス
	FontIndex                int    // フォントコレクションのインデックス
	EmojiFontFile            string // 絵文字用のフォントファイルのパス
	EmojiFontIndex           int    // 絵文字用のフォントコレクションのインデックス
	UseEmojiFont             bool   // 絵文字TTFを使う
	FontSize                 int    // フォントサイズ
	UseAnimation             bool   // アニメーションGIFを生成する
	Delay                    int    // アニメーションのディレイ時間
	LineCount                int    // 入力データのうち何行を1フレーム画像に使うか
	UseSlideAnimation        bool   // スライドアニメーションする
	SlideWidth               int    // スライドする幅
	SlideForever             bool   // スライドを無限にスライドするように描画する
	ToSlackIcon              bool   // Slackのアイコンサイズにする
	PrintEnvironments        bool
	UseShellgeiImagedir      bool
	UseShellgeiEmojiFontfile bool
}

const shellgeiEmojiFontPath = "/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf"

var (
	appconf applicationConfig
)

func init() {
	cobra.OnInitialize()

	RootCommand.Flags().SortFlags = false
	RootCommand.Flags().StringVarP(&appconf.Foreground, "foreground", "g", "white", `foreground text color.
available color types are [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringVarP(&appconf.Background, "background", "b", "black", `background text color.
color types are same as "foreground" option`)

	var font string
	envFontFile := os.Getenv(global.EnvNameFontFile)
	if envFontFile != "" {
		font = envFontFile
	}
	RootCommand.Flags().StringVarP(&appconf.FontFile, "fontfile", "f", font, `font file path.
You can change this default value with environment variables TEXTIMG_FONT_FILE`)
	RootCommand.Flags().IntVarP(&appconf.FontIndex, "fontindex", "x", 0, "")
	appconf.setFontFileAndFontIndex(runtime.GOOS)

	envEmojiFontFile := os.Getenv(global.EnvNameEmojiFontFile)
	RootCommand.Flags().StringVarP(&appconf.EmojiFontFile, "emoji-fontfile", "e", envEmojiFontFile, "emoji font file")
	RootCommand.Flags().IntVarP(&appconf.EmojiFontIndex, "emoji-fontindex", "X", 0, "")

	RootCommand.Flags().BoolVarP(&appconf.UseEmojiFont, "use-emoji-font", "i", false, "use emoji font")
	RootCommand.Flags().BoolVarP(&appconf.UseShellgeiEmojiFontfile, "shellgei-emoji-fontfile", "z", false, `emoji font file for shellgei-bot (path: "`+shellgeiEmojiFontPath+`")`)

	RootCommand.Flags().IntVarP(&appconf.FontSize, "fontsize", "F", 20, "font size")
	RootCommand.Flags().StringVarP(&appconf.Outpath, "out", "o", "", `output image file path.
available image formats are [png | jpg | gif]`)
	RootCommand.Flags().BoolVarP(&appconf.UseShellgeiImagedir, "shellgei-imagedir", "s", false, `image directory path for shellgei-bot (path: "/images/t.png")`)

	RootCommand.Flags().BoolVarP(&appconf.UseAnimation, "animation", "a", false, "generate animation gif")
	RootCommand.Flags().IntVarP(&appconf.Delay, "delay", "d", 20, "animation delay time")
	RootCommand.Flags().IntVarP(&appconf.LineCount, "line-count", "l", 1, "animation input line count")
	RootCommand.Flags().BoolVarP(&appconf.UseSlideAnimation, "slide", "S", false, "use slide animation")
	RootCommand.Flags().IntVarP(&appconf.SlideWidth, "slide-width", "W", 1, "sliding animation width")
	RootCommand.Flags().BoolVarP(&appconf.SlideForever, "forever", "E", false, "sliding forever")
	RootCommand.Flags().BoolVarP(&appconf.PrintEnvironments, "environments", "", false, "print environment variables")
	RootCommand.Flags().BoolVarP(&appconf.ToSlackIcon, "slack", "", false, "resize to slack icon size (128x128 px)")
}

func (a *applicationConfig) setFontFileAndFontIndex(runtimeOS string) {
	if a.FontFile != "" {
		return
	}

	switch runtimeOS {
	case "linux":
		if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
			a.FontFile = "/mnt/c/Windows/Fonts/msgothic.ttc"
			a.FontIndex = 0
			return
		}

		a.FontFile = "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"
		if _, err := os.Stat(a.FontFile); err != nil {
			a.FontFile = "/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc"
		}
		a.FontIndex = 4
	case "windows":
		a.FontFile = `C:\Windows\Fonts\msgothic.ttc`
		a.FontIndex = 0
	case "darwin":
		a.FontFile = "/System/Library/Fonts/AppleSDGothicNeo.ttc"
		a.FontIndex = 0
	case "ios":
		a.FontFile = "/System/Library/Fonts/Core/AppleSDGothicNeo.ttc"
		a.FontIndex = 0
	case "android":
		a.FontFile = "/system/fonts/NotoSansCJK-Regular.ttc"
		a.FontIndex = 4
	}
}

var RootCommand = &cobra.Command{
	Use:     global.AppName,
	Short:   global.AppName + " is command to convert from colored text (ANSI or 256) to image.",
	Example: global.AppName + ` $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: global.Version,
	RunE:    runRootCommand,
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	// コマンドライン引数の取得{{{
	if appconf.PrintEnvironments {
		for _, envName := range global.EnvNames {
			text := fmt.Sprintf("%s=%s", envName, os.Getenv(envName))
			fmt.Println(text)
		}
		return nil
	}

	// シェル芸イメージディレクトリの指定がある時はパスを変更する
	if appconf.UseShellgeiImagedir {
		if appconf.UseAnimation {
			appconf.Outpath = "/images/t.gif"
		} else {
			appconf.Outpath = "/images/t.png"
		}
	}

	if appconf.UseShellgeiEmojiFontfile {
		appconf.EmojiFontFile = shellgeiEmojiFontPath
		appconf.UseEmojiFont = true
	}

	if appconf.UseSlideAnimation {
		appconf.UseAnimation = true
	}

	confForeground, err := optionColorStringToRGBA(appconf.Foreground)
	if err != nil {
		return err
	}

	confBackground, err := optionColorStringToRGBA(appconf.Background)
	if err != nil {
		return err
	}

	// }}}

	// 引数にテキストの指定がなければ標準入力を使用する
	var texts []string
	if len(args) < 1 {
		texts = readStdin()
	} else {
		for _, v := range args {
			texts = append(texts, strings.Split(v, "\n")...)
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
		err := fmt.Errorf("must need input texts.")
		return err
	}

	// スライドアニメーションを使うときはテキストを加工する
	if appconf.UseSlideAnimation {
		texts = toSlideStrings(texts, appconf.LineCount, appconf.SlideWidth, appconf.SlideForever)
	}

	// 拡張子のみ取得
	imgExt := filepath.Ext(strings.ToLower(appconf.Outpath))

	var w *os.File
	if appconf.Outpath == "" {
		// 出力先画像の指定がなく、且つ出力先がパイプならstdout + PNG/GIFと
		// して出力。なければそもそも画像処理しても意味が無いので終了
		fd := os.Stdout.Fd()
		if terminal.IsTerminal(int(fd)) {
			log.Error("image data not written to a terminal. use -o, -s, pipe or redirect.")
			log.Error("for help, type: textimg -h")
			return fmt.Errorf("no output target error")
		}
		w = os.Stdout
		if appconf.UseAnimation {
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
		err := fmt.Errorf("%s is not supported extension.", imgExt)
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

	face := ioimage.ReadFace(appconf.FontFile, appconf.FontIndex, float64(appconf.FontSize))
	var emojiFace font.Face
	if appconf.EmojiFontFile != "" {
		emojiFace = ioimage.ReadFace(appconf.EmojiFontFile, appconf.EmojiFontIndex, float64(appconf.FontSize))
	}
	emojiDir := os.Getenv(global.EnvNameEmojiDir)

	writeConf := ioimage.WriteConfig{
		Foreground:    confForeground,
		Background:    confBackground,
		FontFace:      face,
		EmojiFontFace: emojiFace,
		EmojiDir:      emojiDir,
		UseEmojiFont:  appconf.UseEmojiFont,
		FontSize:      appconf.FontSize,
		UseAnimation:  appconf.UseAnimation,
		Delay:         appconf.Delay,
		LineCount:     appconf.LineCount,
		ToSlackIcon:   appconf.ToSlackIcon,
	}
	if err := ioimage.Write(w, imgExt, texts, writeConf); err != nil {
		return err
	}

	return nil
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
		return zeroColor, fmt.Errorf("illegal RGBA format: " + colstr)
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
	for _, v := range s {
		for _, c := range zwc {
			if v == c {
				continue chars
			}
		}
		ret = append(ret, v)
	}
	return string(ret)
}
