package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/ioimage"
	"github.com/jiro4989/textimg/v3/log"
	"github.com/jiro4989/textimg/v3/parser"
	"github.com/jiro4989/textimg/v3/token"
	"golang.org/x/image/font"
	"golang.org/x/term"
)

type applicationConfig struct {
	Foreground               string // 文字色
	Background               string // 背景色
	Outpath                  string // 画像の出力ファイルパス
	AddTimeStamp             bool   // ファイル名末尾にタイムスタンプ付与
	SaveNumberedFile         bool   // 保存しようとしたファイルがすでに存在する場合に連番を付与する
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
	ResizeWidth              int // 画像の横幅
	ResizeHeight             int // 画像の縦幅

	ForegroundColor color.RGBA // 文字色
	BackgroundColor color.RGBA // 背景色
	Texts           []string
	FileExtension   string
	Writer          *os.File
	FontFace        font.Face
	EmojiFontFace   font.Face
	EmojiDir        string
	Tokens          token.Tokens
}

// adjust はパラメータを調整する。
// 副作用を持つ。
func (a *applicationConfig) Adjust(args []string, ev EnvVars) error {
	a.EmojiDir = ev.EmojiDir

	// シェル芸イメージディレクトリの指定がある時はパスを変更する
	if a.UseShellgeiImagedir {
		var err error
		outDir := ev.OutputDir
		a.Outpath, err = outputImageDir(outDir, a.UseAnimation)
		if err != nil {
			return err
		}
	}

	a.addTimeStampToOutPath(time.Now())
	a.addNumberSuffixToOutPath()

	if a.UseShellgeiEmojiFontfile {
		a.EmojiFontFile = shellgeiEmojiFontPath
		a.UseEmojiFont = true
	}

	if a.UseSlideAnimation {
		a.UseAnimation = true
	}

	var err error
	a.ForegroundColor, err = optionColorStringToRGBA(a.Foreground)
	if err != nil {
		return err
	}

	a.BackgroundColor, err = optionColorStringToRGBA(a.Background)
	if err != nil {
		return err
	}

	// 引数にテキストの指定がなければ標準入力を使用する
	a.Texts = readInputText(args)

	// textsが空のときは警告メッセージを出力して異常終了
	if err := validateInputText(a.Texts); err != nil {
		return err
	}

	// スライドアニメーションを使うときはテキストを加工する
	if a.UseSlideAnimation {
		a.Texts = toSlideStrings(a.Texts, a.LineCount, a.SlideWidth, a.SlideForever)
	}

	// 拡張子のみ取得
	a.FileExtension = filepath.Ext(strings.ToLower(a.Outpath))

	if err := a.setWriter(); err != nil {
		return err
	}

	if err := validateFileExtension(a.FileExtension); err != nil {
		return err
	}

	a.Texts = normalizeTexts(a.Texts)

	a.FontFace, err = ioimage.ReadFace(a.FontFile, a.FontIndex, float64(a.FontSize))
	if err != nil {
		return err
	}

	if a.EmojiFontFile != "" {
		a.EmojiFontFace, err = ioimage.ReadFace(a.EmojiFontFile, a.EmojiFontIndex, float64(a.FontSize))
		if err != nil {
			return err
		}
	}

	a.Tokens, err = parser.Parse(strings.Join(a.Texts, "\n"))
	if err != nil {
		return err
	}

	if a.ToSlackIcon {
		a.ResizeWidth = 128
		a.ResizeHeight = 128
	}

	return nil
}

func (a *applicationConfig) setFontFileAndFontIndex(runtimeOS string) {
	if a.FontFile != "" {
		return
	}

	m := map[string]osDefaultFont{
		"linux": {
			isLinux: true,
		},
		"windows": {
			fontFile:  defaultWindowsFont,
			fontIndex: 0,
		},
		"darwin": {
			fontFile:  defaultDarwinFont,
			fontIndex: 0,
		},
		"ios": {
			fontFile:  defaultIOSFont,
			fontIndex: 0,
		},
		"android": {
			fontFile:  defaultAndroidFont,
			fontIndex: 5,
		},
	}

	if f, ok := m[runtimeOS]; ok {
		// linux だけ特殊なので特別に分岐
		if !f.isLinux {
			a.FontFile = f.fontFile
			a.FontIndex = f.fontIndex
			return
		}

		if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
			a.FontFile = "/mnt/c/Windows/Fonts/msgothic.ttc"
			a.FontIndex = 0
			return
		}

		a.FontFile = defaultLinuxFont1
		if _, err := os.Stat(a.FontFile); err != nil {
			a.FontFile = defaultLinuxFont2
		}
		a.FontIndex = 5
		return
	}
}

// addTimeStampToOutPath はOutpathに指定日時のタイムスタンプを付与する。
func (a *applicationConfig) addTimeStampToOutPath(t time.Time) {
	if !a.AddTimeStamp {
		return
	}

	ext := filepath.Ext(a.Outpath)
	file := strings.TrimSuffix(a.Outpath, ext)
	timestamp := t.Format("2006-01-02-150405")
	a.Outpath = file + "_" + timestamp + ext
}

// addTimeStampToOutPath はOutpathに指定日時のタイムスタンプを付与する。
func (a *applicationConfig) addNumberSuffixToOutPath() {
	if !a.SaveNumberedFile {
		return
	}

	// ファイルが存在しない時は何もしない
	// NOTE: 並列に実行されるとチェックしきれない場合があるけれど許容する
	if _, err := os.Stat(a.Outpath); err != nil {
		return
	}

	fileExt := filepath.Ext(a.Outpath)
	fileName := strings.TrimSuffix(a.Outpath, fileExt)
	i := 2
	for {
		a.Outpath = fmt.Sprintf("%s_%d%s", fileName, i, fileExt)
		_, err := os.Stat(a.Outpath)
		if err != nil {
			return
		}
		i++
	}
}

func (a *applicationConfig) setWriter() error {
	if a.Outpath == "" {
		// 出力先画像の指定がなく、且つ出力先がパイプならstdout + PNG/GIFと
		// して出力。なければそもそも画像処理しても意味が無いので終了
		fd := os.Stdout.Fd()
		if term.IsTerminal(int(fd)) {
			log.Error("image data not written to a terminal. use -o, -s, pipe or redirect.")
			log.Error("for help, type: textimg -h")
			return fmt.Errorf("no output target error")
		}
		a.Writer = os.Stdout
		if a.UseAnimation {
			a.FileExtension = ".gif"
		} else {
			a.FileExtension = ".png"
		}

		return nil
	}

	var err error
	a.Writer, err = os.Create(a.Outpath)
	if err != nil {
		return err
	}
	// NOTE: writerは呼び出し元でクローズする

	return nil
}

func validateInputText(texts []string) error {
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
	return nil
}

// validateFileExtension はファイル拡張子をチェックする。
func validateFileExtension(ext string) error {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		// 何もしない
	default:
		err := fmt.Errorf("%s is not supported extension.", ext)
		return err
	}
	return nil
}

// normalizeTexts はテキストを正規化する。
func normalizeTexts(texts []string) []string {
	result := texts

	// タブ文字は画像描画時に表示されないので暫定対応で半角スペースに置換
	for i, text := range result {
		result[i] = strings.Replace(text, "\t", "  ", -1)
	}

	// ゼロ幅文字を削除
	for i, text := range result {
		result[i] = removeZeroWidthCharacters(text)
	}

	return result
}

func readInputText(args []string) []string {
	var texts []string
	if len(args) < 1 {
		texts = readStdin()
	} else {
		for _, v := range args {
			texts = append(texts, strings.Split(v, "\n")...)
		}
	}
	return texts
}
