package ioimage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFace(t *testing.T) {
	testdataDir := filepath.Join("..", "testdata", "in")

	type TestData struct {
		desc        string
		inFontPath  string
		inFontIndex int
		wantErr     bool
	}
	tests := []TestData{
		{
			desc:        "正常系: font.Faceが取得できる",
			inFontPath:  "/tmp/MyricaM.TTC",
			inFontIndex: 0,
			wantErr:     false,
		},
		{
			desc:        "正常系: 存在しないファイルの場合もエラーにはならない",
			inFontPath:  "/tmp/寿司",
			inFontIndex: 0,
			wantErr:     false,
		},
		{
			desc:        "異常系: パスとしては存在するがディレクトリの場合はエラー",
			inFontPath:  "/tmp",
			inFontIndex: 0,
			wantErr:     true,
		},
		{
			desc:        "異常系: ファイルは存在するけれど、フォントファイルじゃない時はエラー (ttc)",
			inFontPath:  filepath.Join(testdataDir, "illegal_font.ttc"),
			inFontIndex: 0,
			wantErr:     true,
		},
		{
			desc:        "異常系: ファイルは存在するけれど、フォントファイルじゃない時はエラー (otc)",
			inFontPath:  filepath.Join(testdataDir, "illegal_font.otc"),
			inFontIndex: 0,
			wantErr:     true,
		},
		{
			desc:        "異常系: ファイルは存在するけれど、フォントファイルじゃない時はエラー (txt)",
			inFontPath:  filepath.Join(testdataDir, "illegal_font.txt"),
			inFontIndex: 0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, err := ReadFace(tt.inFontPath, tt.inFontIndex, 20)
			if tt.wantErr {
				assert.Nil(got)
				assert.Error(err)
				return
			}

			assert.NotNil(got)
			assert.NoError(err)
		})
	}
}
