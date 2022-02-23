package image

import (
	"fmt"
	"os"
)

var (
	// 絵文字描画の際に、普通に描画してほしいけれど絵文字としても定義されている
	// 文字のコードポイント
	exRunes = []rune{
		0x0023, // #
		0x002A, // *
		0x0030, // 0
		0x0031, // 1
		0x0032, // 2
		0x0033, // 3
		0x0034, // 4
		0x0035, // 5
		0x0036, // 6
		0x0037, // 7
		0x0038, // 8
		0x0039, // 9
		0x00A9, // ©
		0x00AE, // ®️
	}
)

// コードポイントに対応する画像ファイルかどうかを判定する。
// 画像ファイルだった場合は当該画像ファイルのパスを返却する。
func isEmoji(r rune, emojiDir string) (bool, string) {
	path := fmt.Sprintf("%s/emoji_u%.4x.png", emojiDir, r)
	_, err := os.Stat(path)
	if err == nil && !isExceptionallyCodePoint(r) {
		return true, path
	}
	return false, ""
}

// r が例外的なコードポイントに存在するかを判定する。
// http://unicode.org/Public/emoji/4.0/emoji-data.txt
//
// ここでtrueを返す文字は、絵文字データ的には絵文字ではあるものの、
// シェル芸bot環境ではテキストとして表示したいので例外的に除外するために指定して
// いる。
func isExceptionallyCodePoint(r rune) bool {
	for _, ex := range exRunes {
		if r == ex {
			return true
		}
	}

	return false
}
