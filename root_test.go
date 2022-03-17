package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationConfig_AddNumberSuffixToOutPath(t *testing.T) {
	testdataDir := filepath.Join(".", "testdata", "in")

	existedFile := filepath.Join(testdataDir, "appconf_add_number_suffix_test_case1.png")
	existedFileWant := filepath.Join(testdataDir, "appconf_add_number_suffix_test_case1_2.png")

	notExistedFile := filepath.Join(testdataDir, "appconf_add_number_suffix_sushi.png")

	type TestData struct {
		desc               string
		inOutpath          string
		inSaveNumberedFile bool
		want               string
	}
	tests := []TestData{
		{
			desc:               "正常系: フラグfalseの場合は変更なし",
			inOutpath:          existedFile,
			inSaveNumberedFile: false,
			want:               existedFile,
		},
		{
			desc:               "正常系: フラグtrueの場合は連番を付与する",
			inOutpath:          existedFile,
			inSaveNumberedFile: true,
			want:               existedFileWant,
		},
		{
			desc:               "正常系: フラグtrueの場合でも、ファイルが存在しなければ何もしない",
			inOutpath:          notExistedFile,
			inSaveNumberedFile: true,
			want:               notExistedFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			a := applicationConfig{
				Outpath:          tt.inOutpath,
				SaveNumberedFile: tt.inSaveNumberedFile,
			}
			a.addNumberSuffixToOutPath()

			assert.Equal(tt.want, a.Outpath)
		})
	}
}
