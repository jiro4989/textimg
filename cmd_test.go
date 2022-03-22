// 実行可能バイナリをビルドして、バイナリのオプションやファイル生成の挙動をテス
// トする。

package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// type (
// 	ANSIColor struct {
// 		text           string
// 		escapeSequence string
// 	}
// )
//
// const (
// 	outDir               = "testdata/out"
// 	bin                  = "./textimg"
// 	ansiColorReset       = "\x1b[0m"
// 	ansiColorFGBlack     = "\x1b[30m"
// 	ansiColorFGRed       = "\x1b[31m"
// 	ansiColorFGGreen     = "\x1b[32m"
// 	ansiColorFGYellow    = "\x1b[33m"
// 	ansiColorFGBlue      = "\x1b[34m"
// 	ansiColorFGMagenta   = "\x1b[35m"
// 	ansiColorFGCyan      = "\x1b[36m"
// 	ansiColorFGLightGray = "\x1b[37m"
// 	ansiColorBGBlack     = "\x1b[40m"
// 	ansiColorBGRed       = "\x1b[41m"
// 	ansiColorBGGreen     = "\x1b[42m"
// 	ansiColorBGYellow    = "\x1b[43m"
// 	ansiColorBGBlue      = "\x1b[44m"
// 	ansiColorBGMagenta   = "\x1b[45m"
// 	ansiColorBGCyan      = "\x1b[46m"
// 	ansiColorBGLightGray = "\x1b[47m"
// )
//
// var (
// 	foregrounds = []ANSIColor{
// 		{text: "black", escapeSequence: ansiColorFGBlack},
// 		{text: "red", escapeSequence: ansiColorFGRed},
// 		{text: "green", escapeSequence: ansiColorFGGreen},
// 		{text: "yellow", escapeSequence: ansiColorFGYellow},
// 		{text: "blue", escapeSequence: ansiColorFGBlue},
// 		{text: "magenta", escapeSequence: ansiColorFGMagenta},
// 		{text: "cyan", escapeSequence: ansiColorFGCyan},
// 		{text: "light_gray", escapeSequence: ansiColorFGLightGray},
// 	}
// 	backgrounds = []ANSIColor{
// 		{text: "black", escapeSequence: ansiColorBGBlack},
// 		{text: "red", escapeSequence: ansiColorBGRed},
// 		{text: "green", escapeSequence: ansiColorBGGreen},
// 		{text: "yellow", escapeSequence: ansiColorBGYellow},
// 		{text: "blue", escapeSequence: ansiColorBGBlue},
// 		{text: "magenta", escapeSequence: ansiColorBGMagenta},
// 		{text: "cyan", escapeSequence: ansiColorBGCyan},
// 		{text: "light_gray", escapeSequence: ansiColorBGLightGray},
// 	}
// 	fgBgLine string
//
// 	testdataInputDir = filepath.Join(".", "testdata", "in")
// )
//
// func setup() {
// 	var lines []string
// 	for _, e := range foregrounds {
// 		line := e.escapeSequence + "Text" + ansiColorReset
// 		lines = append(lines, line)
// 	}
// 	for _, e := range backgrounds {
// 		line := e.escapeSequence + "Text" + ansiColorReset
// 		lines = append(lines, line)
// 	}
// 	fgBgLine = strings.Join(lines, "\n")
//
// 	// ディレクトリが存在しない場合はエラーになるけれど無視
// 	os.RemoveAll(outDir)
// 	err := os.MkdirAll(outDir, os.ModePerm)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := exec.Command("go", "build").Run(); err != nil {
// 		panic(err)
// 	}
// }
//
// func teardown() {
// 	// 何もしない
// }
//
// func TestMain(m *testing.M) {
// 	setup()
// 	ret := m.Run()
// 	if ret == 0 {
// 		teardown()
// 	}
// 	os.Exit(ret)
// }
//
// func TestNormalText(t *testing.T) {
// 	tests := []struct {
// 		desc string
// 		in   string
// 	}{
// 		{
// 			desc: "エスケープシーケンスの無いテキストは普通にテキストを描画する",
// 			in:   fmt.Sprintf("echo test | %s > %s/%s", bin, outDir, "normal_text.png"),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			err := exec.Command("bash", "-c", tt.in).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestANSIColorFrom30To37AndFrom40To47(t *testing.T) {
// 	tests := []struct {
// 		desc    string
// 		inCmd   string
// 		inColor []ANSIColor
// 	}{
// 		{
// 			desc:    "ANSIColorの前景色が変わる (%s)",
// 			inCmd:   fmt.Sprintf("echo -e '%s' | %s > %s/%s", "%s%s%s", bin, outDir, "ansi_color_foreground_%s.png"),
// 			inColor: foregrounds,
// 		},
// 		{
// 			desc:    "ANSIColorの背景色が変わる (%s)",
// 			inCmd:   fmt.Sprintf("echo -e '%s' | %s > %s/%s", "%s%s%s", bin, outDir, "ansi_color_background_%s.png"),
// 			inColor: backgrounds,
// 		},
// 	}
// 	for _, tt := range tests {
// 		for _, c := range tt.inColor {
// 			desc := fmt.Sprintf(tt.desc, c.text)
// 			cmd := fmt.Sprintf(tt.inCmd, c.escapeSequence, c.text, ansiColorReset, c.text)
//
// 			t.Run(desc, func(t *testing.T) {
// 				err := exec.Command("bash", "-c", cmd).Run()
// 				assert.Nil(t, err, tt.desc)
// 			})
// 		}
// 	}
// }
//
// func TestANSIColorSwitch(t *testing.T) {
// 	tests := []struct {
// 		desc string
// 		in   string
// 	}{
// 		{
// 			desc: "前景色が1行で複数の色に切り替わる",
// 			in: fmt.Sprintf(
// 				"echo -e '%s' | %s > %s/%s",
// 				strings.Join([]string{ansiColorFGRed, "Red", ansiColorFGGreen, "緑", ansiColorFGBlue, "あお", ansiColorReset}, ""),
// 				bin, outDir, "ansi_color_switch_foreground.png",
// 			),
// 		},
// 		{
// 			desc: "背景色が1行で複数の色に切り替わる",
// 			in: fmt.Sprintf(
// 				"echo -e '%s' | %s > %s/%s",
// 				strings.Join([]string{ansiColorBGRed, "Red", ansiColorBGGreen, "緑", ansiColorBGBlue, "あお", ansiColorReset}, ""),
// 				bin, outDir, "ansi_color_switch_background.png",
// 			),
// 		},
// 		{
// 			desc: "前景色と背景色が1行で複数の色に切り替わる",
// 			in: fmt.Sprintf(
// 				"echo -e '%s' | %s > %s/%s",
// 				strings.Join([]string{ansiColorFGRed, "Red", ansiColorBGGreen, "緑", ansiColorFGBlue, "あお", ansiColorReset}, ""),
// 				bin, outDir, "ansi_color_switch_foreground_background.png",
// 			),
// 		},
// 		{
// 			desc: "複数行のテキスト",
// 			in: fmt.Sprintf(
// 				"echo -e '%s' | %s > %s/%s",
// 				strings.Repeat(strings.Join([]string{ansiColorFGRed, "あかRed", ansiColorBGGreen, "緑Green", ansiColorFGBlue, "あおBlue", ansiColorReset, "\n"}, ""), 10),
// 				bin, outDir, "ansi_color_switch_multiline.png",
// 			),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			err := exec.Command("bash", "-c", tt.in).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestGrep(t *testing.T) {
// 	cmd := fmt.Sprintf(`echo TestAbcTest | grep --color=always Abc | %s > %s/%s`, bin, outDir, "grep.png")
// 	err := exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "grepで色がつく")
// }
//
// func TestFGBGOptionRGBAFormat(t *testing.T) {
// 	tests := []struct {
// 		desc   string
// 		inOpt  string
// 		inRGBA string
// 		fn     string
// 	}{
// 		{
// 			desc:   "前景色のショートオプション",
// 			inOpt:  "-g",
// 			inRGBA: "50,100,12,255",
// 			fn:     "foreground_short_random",
// 		},
// 		{
// 			desc:   "前景色のロングオプション",
// 			inOpt:  "--foreground",
// 			inRGBA: "255,0,0,255",
// 			fn:     "foreground_long_red",
// 		},
// 		{
// 			desc:   "背景色のショートオプション",
// 			inOpt:  "-b",
// 			inRGBA: "0,100,0,160",
// 			fn:     "background_short_blue",
// 		},
// 		{
// 			desc:   "背景色のロングオプション",
// 			inOpt:  "--background",
// 			inRGBA: "0,0,255,100",
// 			fn:     "background_long_green",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			fn := fmt.Sprintf("fg_bg_option_rgba_format_%s.png", tt.fn)
// 			cmd := fmt.Sprintf("echo -e 'A%sRed%sB' | %s %s %s > %s/%s", ansiColorFGRed, ansiColorReset, bin, tt.inOpt, tt.inRGBA, outDir, fn)
// 			err := exec.Command("bash", "-c", cmd).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestOutFileOption(t *testing.T) {
// 	for i, opt := range []string{"-o", "--out"} {
// 		msg := opt + "オプションで画像ファイルが生成できる"
// 		t.Run(msg, func(t *testing.T) {
// 			fn := fmt.Sprintf("output_option_%d.png", i)
// 			cmd := fmt.Sprintf(`echo -e '%sText%s' | %s %s %s/%s`, ansiColorFGBlue, ansiColorReset, bin, opt, outDir, fn)
// 			err := exec.Command("bash", "-c", cmd).Run()
// 			assert.Nil(t, err, msg)
// 		})
// 	}
// }
//
// func TestSimple(t *testing.T) {
// 	const t1 = ansiColorFGBlue + "Text" + ansiColorReset
// 	tests := []struct {
// 		desc   string
// 		inText string
// 		outFn  string
// 	}{
// 		{
// 			desc:   "JPGファイルが生成できる",
// 			inText: t1,
// 			outFn:  "simple_jpg.jpg",
// 		},
// 		{
// 			desc:   "GIFファイルが生成できる",
// 			inText: t1,
// 			outFn:  "simple_gif.gif",
// 		},
// 		{
// 			desc:   "前景色と背景色が反転する",
// 			inText: "\x1b[31;42mRedGreen\x1b[7mRedGreen",
// 			outFn:  "simple_reverse.png",
// 		},
// 		{
// 			desc:   "1行の絵文字",
// 			inText: "あ😃a👍！👀ん👄",
// 			outFn:  "emoji1.png",
// 		},
// 		{
// 			desc:   "3行の絵文字",
// 			inText: "ab😃cd👍ef👀gh👄\n😃12👍34👀5a👄あ\n😃a👍b👀c👄dabcd",
// 			outFn:  "emoji2.png",
// 		},
// 		{desc: "ゼロ幅文字 (U+200B)", inText: "A \u200B B", outFn: "zws_u002b.png"},
// 		{desc: "ゼロ幅文字 (U+200C)", inText: "A \u200C B", outFn: "zws_u002c.png"},
// 		{desc: "ゼロ幅文字 (U+200D)", inText: "A \u200D B", outFn: "zws_u002d.png"},
// 		{desc: "ゼロ幅文字 (U+FEFF)", inText: "A \uFEFF B", outFn: "zws_ufeff.png"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			cmd := fmt.Sprintf("echo -e '%s' | %s -o %s/%s", tt.inText, bin, outDir, tt.outFn)
// 			err := exec.Command("bash", "-c", cmd).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestAnimationGIF(t *testing.T) {
// 	tests := []struct {
// 		desc       string
// 		inText     string
// 		inLine     int
// 		inDuration int
// 	}{
// 		{desc: "1行のアニメ (デフォルト)", inText: fgBgLine, inLine: 1, inDuration: -1},
// 		{desc: "1行のアニメ", inText: fgBgLine, inLine: 1, inDuration: 5},
// 		{desc: "2行のアニメ", inText: fgBgLine, inLine: 2, inDuration: 10},
// 		{desc: "4行のアニメ", inText: fgBgLine, inLine: 4, inDuration: 20},
// 		{desc: "8行のアニメ", inText: fgBgLine, inLine: 8, inDuration: 30},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			fn := fmt.Sprintf("animation_gif_%d_line.gif", tt.inLine)
// 			var cmd string
// 			if tt.inDuration == -1 {
// 				cmd = fmt.Sprintf("echo -e '%s' | %s -a -l %d -o %s/%s", tt.inText, bin, tt.inLine, outDir, fn)
// 			} else {
// 				cmd = fmt.Sprintf("echo -e '%s' | %s -a -l %d -d %d -o %s/%s", tt.inText, bin, tt.inLine, tt.inDuration, outDir, fn)
// 			}
// 			err := exec.Command("bash", "-c", cmd).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestANSIColorExt(t *testing.T) {
// 	tests := []struct {
// 		desc   string
// 		inCode int
// 		inFmt  string
// 		outFn  string
// 	}{
// 		{desc: "前景色256色", inCode: 38, inFmt: "\x1b[%d;5;%dm%03d%s", outFn: "anci_color_ext_256_foreground.png"},
// 		{desc: "背景色256色", inCode: 48, inFmt: "\x1b[%d;5;%dm%03d%s", outFn: "anci_color_ext_256_background.png"},
// 		{desc: "前景色RGBグラデーション色", inCode: 38, inFmt: "\x1b[%d;2;%d;0;0m%03d%s", outFn: "anci_color_ext_rgb_foreground.png"},
// 		{desc: "背景色RGBグラデーション色", inCode: 48, inFmt: "\x1b[%d;2;%d;0;0m%03d%s", outFn: "anci_color_ext_rgb_background.png"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			var text string
// 			for i := 0; i < 256; i++ {
// 				es := fmt.Sprintf(tt.inFmt, tt.inCode, i, i, ansiColorReset)
// 				text += es
// 				if (i+1)%16 == 0 {
// 					text += "\n"
// 				}
// 			}
// 			cmd := fmt.Sprintf("echo -e '%s' | %s -o %s/%s", text, bin, outDir, tt.outFn)
// 			err := exec.Command("bash", "-c", cmd).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestSlideAnimation(t *testing.T) {
// 	format := "seq 5 | %s -S -l 3 %s -o %s/%s"
//
// 	var cmd string
// 	var err error
// 	cmd = fmt.Sprintf(format, bin, "", outDir, "slide_animation_line_3_width_default.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "行3でスライド幅デフォルトのGIFアニメ")
//
// 	cmd = fmt.Sprintf(format, bin, "-W 2", outDir, "slide_animation_line_3_width_2.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "行3でスライド幅2のGIFアニメ")
// }
//
// func TestSlideAnimationRainbow(t *testing.T) {
// 	format := "echo -e '%s' | %s -S -l 5 %s -o %s/%s"
//
// 	var cmd string
// 	var err error
// 	cmd = fmt.Sprintf(format, fgBgLine, bin, "", outDir, "slide_animation_rainbow.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "スライドアニメ5行")
//
// 	cmd = fmt.Sprintf(format, fgBgLine, bin, "-E", outDir, "slide_animation_rainbow_forever.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "スライドアニメ5行を無限に")
// }
//
// func TestEmojiZ(t *testing.T) {
// 	format := "echo -e '%s' | %s -z -o %s/%s"
//
// 	var cmd string
// 	var err error
// 	cmd = fmt.Sprintf(format, "ab😃cd👍ef👀gh👄\n😃12👍34👀5a👄あ\n😃a👍b👀c👄dabcd", bin, outDir, "emoji_z.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "絵文字Z")
// }
//
// func TestError(t *testing.T) {
// 	var cmd string
// 	var err error
//
// 	cmd = fmt.Sprintf(`echo -n "" | %s -o %s/%s`, bin, outDir, "empty.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.NotNil(t, err, "空の標準入力はエラーを返す")
//
// 	cmd = fmt.Sprintf(`%s "" -o %s/%s`, bin, outDir, "empty.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.NotNil(t, err, "空の引数はエラーを返す")
//
// 	cmd = fmt.Sprintf(`%s "%s" -o %s/%s`, "\n\n\n", bin, outDir, "empty.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.NotNil(t, err, "改行のみの入力はエラーを返す")
// }
//
// func TestNoRedirectNoPipe(t *testing.T) {
// 	var cmd string
// 	var err error
//
// 	const msg = "test"
//
// 	// Note:
// 	//   go testで実行する時は端末ではないためか、絶対にnilになってしまう
// 	//   端末からコマンドを実行する時では動作確認できているためコメントアウト
// 	// err = exec.Command(bin, msg).Run()
// 	// assert.NotNil(t, err, "出力先を未指定の時は異常終了")
//
// 	cmd = fmt.Sprintf(`%s %s > /dev/null`, bin, msg)
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "出力先(リダイレクト)が指定されている時は正常終了")
//
// 	cmd = fmt.Sprintf(`%s %s | cat > /dev/null`, bin, msg)
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "出力先(パイプ)が指定されている時は正常終了")
// }
