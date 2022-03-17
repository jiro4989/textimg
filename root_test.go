package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputImageDir(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.NoError(t, err)
	pictDir := filepath.Join(home, "Pictures")

	type TestData struct {
		desc           string
		inEnvDir       string
		inUseAnimation bool
		wantPath       string
		wantErr        bool
	}
	tests := []TestData{
		{
			desc:           "正常系: Env未設定の場合はホームディレクトリ配下のPictures配下が返る",
			inEnvDir:       "",
			inUseAnimation: false,
			wantPath:       filepath.Join(pictDir, "t.png"),
			wantErr:        false,
		},
		{
			desc:           "正常系: animation trueの場合は Basenameが t.gif になる",
			inEnvDir:       "",
			inUseAnimation: true,
			wantPath:       filepath.Join(pictDir, "t.gif"),
			wantErr:        false,
		},
		{
			desc:           "正常系: Envが設定されている場合は設定されている値が優先される",
			inEnvDir:       filepath.Join(".", "sushi"),
			inUseAnimation: false,
			wantPath:       filepath.Join(".", "sushi", "t.png"),
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, err := outputImageDir(tt.inEnvDir, tt.inUseAnimation)
			if tt.wantErr {
				assert.Equal("", got)
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(tt.wantPath, got)
		})
	}
}

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
