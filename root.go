package main

import (
	"fmt"
	c "image/color"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/image"
	"github.com/jiro4989/textimg/v3/internal/global"
	"github.com/jiro4989/textimg/v3/parser"

	"github.com/spf13/cobra"
)

const shellgeiEmojiFontPath = "/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf"

var (
	appconf applicationConfig
	envvars EnvVars
)

func init() {
	envvars = NewEnvVars()
	cobra.OnInitialize()

	RootCommand.Flags().SortFlags = false
	RootCommand.Flags().StringVarP(&appconf.Foreground, "foreground", "g", "white", `foreground text color.
available color types are [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringVarP(&appconf.Background, "background", "b", "black", `background text color.
color types are same as "foreground" option`)

	var font string
	envFontFile := envvars.FontFile
	if envFontFile != "" {
		font = envFontFile
	}
	RootCommand.Flags().StringVarP(&appconf.FontFile, "fontfile", "f", font, `font file path.
You can change this default value with environment variables TEXTIMG_FONT_FILE`)
	RootCommand.Flags().IntVarP(&appconf.FontIndex, "fontindex", "x", 0, "")
	appconf.setFontFileAndFontIndex(runtime.GOOS)

	envEmojiFontFile := envvars.EmojiFontFile
	RootCommand.Flags().StringVarP(&appconf.EmojiFontFile, "emoji-fontfile", "e", envEmojiFontFile, "emoji font file")
	RootCommand.Flags().IntVarP(&appconf.EmojiFontIndex, "emoji-fontindex", "X", 0, "")

	RootCommand.Flags().BoolVarP(&appconf.UseEmojiFont, "use-emoji-font", "i", false, "use emoji font")
	RootCommand.Flags().BoolVarP(&appconf.UseShellgeiEmojiFontfile, "shellgei-emoji-fontfile", "z", false, `emoji font file for shellgei-bot (path: "`+shellgeiEmojiFontPath+`")`)

	RootCommand.Flags().IntVarP(&appconf.FontSize, "fontsize", "F", 20, "font size")
	RootCommand.Flags().StringVarP(&appconf.Outpath, "out", "o", "", `output image file path.
available image formats are [png | jpg | gif]`)
	RootCommand.Flags().BoolVarP(&appconf.AddTimeStamp, "timestamp", "t", false, `add time stamp to output image file path.`)
	RootCommand.Flags().BoolVarP(&appconf.SaveNumberedFile, "numbered", "n", false, `add number-suffix to filename when the output file was existed.
ex: t_2.png`)
	RootCommand.Flags().BoolVarP(&appconf.UseShellgeiImagedir, "shellgei-imagedir", "s", false, `image directory path (path: "$HOME/Pictures/t.png" or "$TEXTIMG_OUTPUT_DIR/t.png")`)

	RootCommand.Flags().BoolVarP(&appconf.UseAnimation, "animation", "a", false, "generate animation gif")
	RootCommand.Flags().IntVarP(&appconf.Delay, "delay", "d", 20, "animation delay time")
	RootCommand.Flags().IntVarP(&appconf.LineCount, "line-count", "l", 1, "animation input line count")
	RootCommand.Flags().BoolVarP(&appconf.UseSlideAnimation, "slide", "S", false, "use slide animation")
	RootCommand.Flags().IntVarP(&appconf.SlideWidth, "slide-width", "W", 1, "sliding animation width")
	RootCommand.Flags().BoolVarP(&appconf.SlideForever, "forever", "E", false, "sliding forever")
	RootCommand.Flags().BoolVarP(&appconf.PrintEnvironments, "environments", "", false, "print environment variables")
	RootCommand.Flags().BoolVarP(&appconf.ToSlackIcon, "slack", "", false, "resize to slack icon size (128x128 px)")
	RootCommand.Flags().IntVarP(&appconf.ResizeWidth, "resize-width", "", 0, "resize width")
	RootCommand.Flags().IntVarP(&appconf.ResizeHeight, "resize-height", "", 0, "resize height")
}

type osDefaultFont struct {
	fontFile  string
	fontIndex int
	isLinux   bool
}

const (
	defaultWindowsFont = `C:\Windows\Fonts\msgothic.ttc`
	defaultDarwinFont  = "/System/Library/Fonts/AppleSDGothicNeo.ttc"
	defaultIOSFont     = "/System/Library/Fonts/Core/AppleSDGothicNeo.ttc"
	defaultAndroidFont = "/system/fonts/NotoSansCJK-Regular.ttc"
	defaultLinuxFont1  = "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"
	defaultLinuxFont2  = "/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc"
)

var RootCommand = &cobra.Command{
	Use:     global.AppName,
	Short:   global.AppName + " is command to convert from colored text (ANSI or 256) to image.",
	Example: global.AppName + ` $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: global.Version,
	RunE:    runRootCommand,
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	if appconf.PrintEnvironments {
		for _, envName := range global.EnvNames {
			text := fmt.Sprintf("%s=%s", envName, os.Getenv(envName))
			fmt.Println(text)
		}
		return nil
	}

	if err := appconf.Adjust(args, envvars); err != nil {
		return err
	}
	defer appconf.writer.Close()

	ls := appconf.tokens.StringLines()
	param := &image.ImageParam{
		BaseWidth:          parser.StringWidth(ls),
		BaseHeight:         len(ls),
		ForegroundColor:    c.RGBA(appconf.ForegroundColor),
		BackgroundColor:    c.RGBA(appconf.BackgroundColor),
		FontFace:           appconf.fontFace,
		EmojiFontFace:      appconf.emojiFontFace,
		EmojiDir:           appconf.emojiDir,
		FontSize:           appconf.FontSize,
		Delay:              appconf.Delay,
		AnimationLineCount: appconf.LineCount,
		ResizeWidth:        appconf.ResizeWidth,
		ResizeHeight:       appconf.ResizeHeight,
	}
	img := image.NewImage(param)
	if err := img.Draw(appconf.tokens); err != nil {
		return err
	}
	if err := img.Encode(appconf.writer, appconf.fileExtension); err != nil {
		return err
	}

	return nil
}

// outputImageDir は `-s` オプションで保存するさきのディレクトリパスを返す。
func outputImageDir(outDir string, useAnimation bool) (string, error) {
	if outDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		outDir = filepath.Join(homeDir, "Pictures")
	}

	if useAnimation {
		return filepath.Join(outDir, "t.gif"), nil
	}

	return filepath.Join(outDir, "t.png"), nil
}

// オプション引数のbackgroundは２つの書き方を許容する。
// 1. black といった色の直接指定
// 2. RGBAのカンマ区切り指定
//    書式: R,G,B,A
//    赤色の例: 255,0,0,255
func optionColorStringToRGBA(colstr string) (color.RGBA, error) {
	// "black"といった色名称でマッチするものがあれば返す
	colstr = strings.ToLower(colstr)
	col := color.StringMap[colstr]
	zeroColor := color.RGBA{}
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
	c := color.RGBA{
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
