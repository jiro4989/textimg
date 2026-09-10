package image

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 絵文字画像ファイルを模したダミーファイルを作成する。
// isEmoji はファイルの存在のみを判定に使うため、中身は空でよい。
func createDummyEmojiFiles(t *testing.T, runes ...rune) string {
	t.Helper()

	dir := t.TempDir()
	for _, r := range runes {
		path := filepath.Join(dir, fmt.Sprintf("emoji_u%.4x.png", r))
		if err := os.WriteFile(path, []byte{}, 0600); err != nil {
			t.Fatalf("ダミーの絵文字画像ファイルの作成に失敗した: %v", err)
		}
	}

	return dir
}

func TestIsEmoji(t *testing.T) {
	// 0x00a9 (©) は例外的なコードポイントだが、ファイルは存在させる。
	// 「ファイルがあっても例外なら false」を確かめるため。
	emojiDir := createDummyEmojiFiles(t, 0x1f600, 0x2600, 0x00b6, 0x0023, 0x00a9)

	// os.Stat はディレクトリに対しても成功するため、絵文字画像と同じ名前の
	// ディレクトリがあると絵文字として扱われる。現在の挙動を固定しておく。
	dirPath := filepath.Join(emojiDir, fmt.Sprintf("emoji_u%.4x.png", 0x1f601))
	if err := os.Mkdir(dirPath, 0700); err != nil {
		t.Fatalf("ディレクトリの作成に失敗した: %v", err)
	}

	notExistDir := filepath.Join(t.TempDir(), "notexist")

	tests := []struct {
		desc     string
		r        rune
		emojiDir string
		wantOK   bool
		wantPath string
	}{
		{
			desc:     "画像ファイルが存在する絵文字はパスを返す",
			r:        0x1f600, // 😀
			emojiDir: emojiDir,
			wantOK:   true,
			wantPath: filepath.Join(emojiDir, "emoji_u1f600.png"),
		},
		{
			desc:     "コードポイントが4桁の絵文字もパスを返す",
			r:        0x2600, // ☀
			emojiDir: emojiDir,
			wantOK:   true,
			wantPath: filepath.Join(emojiDir, "emoji_u2600.png"),
		},
		{
			desc:     "コードポイントが4桁未満のときは0埋めしたファイル名で探す",
			r:        0x00b6, // ¶
			emojiDir: emojiDir,
			wantOK:   true,
			wantPath: filepath.Join(emojiDir, "emoji_u00b6.png"),
		},
		{
			desc:     "同名のディレクトリがあると絵文字として扱う",
			r:        0x1f601, // 😁
			emojiDir: emojiDir,
			wantOK:   true,
			wantPath: dirPath,
		},
		{
			desc:     "画像ファイルが存在しない絵文字はfalseを返す",
			r:        0x1f602, // 😂
			emojiDir: emojiDir,
			wantOK:   false,
			wantPath: "",
		},
		{
			desc:     "画像ファイルが存在しても#はfalseを返す",
			r:        '#',
			emojiDir: emojiDir,
			wantOK:   false,
			wantPath: "",
		},
		{
			desc:     "画像ファイルが存在しても©はfalseを返す",
			r:        0x00a9, // ©
			emojiDir: emojiDir,
			wantOK:   false,
			wantPath: "",
		},
		{
			desc:     "通常の文字はfalseを返す",
			r:        'a',
			emojiDir: emojiDir,
			wantOK:   false,
			wantPath: "",
		},
		{
			desc:     "絵文字ディレクトリが存在しないときはfalseを返す",
			r:        0x1f600, // 😀
			emojiDir: notExistDir,
			wantOK:   false,
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)
			gotOK, gotPath := isEmoji(tt.r, tt.emojiDir)
			assert.Equal(tt.wantOK, gotOK)
			assert.Equal(tt.wantPath, gotPath)
		})
	}
}

func TestIsExceptionallyCodePoint(t *testing.T) {
	// 期待する例外集合は実装の exRunes を参照せず、ここに書き下す。
	// 実装のデータで実装を検証すると、集合が変わっても検出できないため。
	// http://unicode.org/Public/emoji/4.0/emoji-data.txt
	wantExceptions := []rune{
		'#', '*',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		0x00A9, // ©
		0x00AE, // ®
	}

	// 例外集合の隣接文字。境界がずれたら落ちるように入れてある。
	wantNormals := []rune{
		'"', '$', ')', '+', '/', ':',
		0x00A8, // ¨
		0x00AA, // ª
		0x00AD, // ソフトハイフン
		0x00AF, // ¯
		'a', 'A',
		0x0000,
		0x1f600,  // 😀
		0x10FFFF, // Unicodeの上限
	}

	for _, r := range wantExceptions {
		t.Run(fmt.Sprintf("例外: %#U", r), func(t *testing.T) {
			assert.True(t, isExceptionallyCodePoint(r))
		})
	}

	for _, r := range wantNormals {
		t.Run(fmt.Sprintf("例外ではない: %#U", r), func(t *testing.T) {
			assert.False(t, isExceptionallyCodePoint(r))
		})
	}
}

func TestIsLinefeed(t *testing.T) {
	// LF だけを改行とみなし、他の改行系文字は改行として扱わない。
	tests := []struct {
		desc string
		r    rune
		want bool
	}{
		{desc: "LFは改行", r: '\n', want: true},
		{desc: "CRは改行ではない", r: '\r', want: false},
		{desc: "垂直タブは改行ではない", r: '\v', want: false},
		{desc: "改ページは改行ではない", r: '\f', want: false},
		{desc: "NEL(U+0085)は改行ではない", r: 0x0085, want: false},
		{desc: "LINE SEPARATOR(U+2028)は改行ではない", r: 0x2028, want: false},
		{desc: "PARAGRAPH SEPARATOR(U+2029)は改行ではない", r: 0x2029, want: false},
		{desc: "タブは改行ではない", r: '\t', want: false},
		{desc: "通常の文字は改行ではない", r: 'a', want: false},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, isLinefeed(tt.r))
		})
	}
}
