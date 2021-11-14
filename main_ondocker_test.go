//go:build ondocker
// +build ondocker

//
// ディレクトリを作成してテストしたりするので明示的にビルドタグを指定しないとテ
// ストされないようにしている。

package main

// NOTE: 他のテストを破壊してしまうので一時的にコメントアウト

// import (
// 	"os"
// 	"path/filepath"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
//
// 	"github.com/jiro4989/textimg/internal/global"
// )
//
// func TestMain_SaveImageToDefaultDirectory(t *testing.T) {
// 	assert := assert.New(t)
//
// 	// 画像ディレクトリを作成
// 	pictDir, err := os.UserHomeDir()
// 	assert.NoError(err)
// 	pictDir = filepath.Join(pictDir, "Pictures")
// 	assert.NoError(os.Mkdir(pictDir, os.ModePerm))
//
// 	os.Args = []string{"", "sample1", "-s"}
// 	assert.Equal(0, Main())
//
// 	// 画像が生成されていること
// 	_, err = os.Stat(filepath.Join(pictDir, "t.png"))
// 	assert.NoError(err)
// }
//
// func TestMain_SaveImageToEnvDirectory(t *testing.T) {
// 	assert := assert.New(t)
//
// 	// 画像ディレクトリを作成
// 	pictDir := filepath.Join("/", "tmp", "tmpout")
// 	os.Setenv(global.EnvNameOutputDir, pictDir)
// 	assert.NoError(os.Mkdir(pictDir, os.ModePerm))
//
// 	os.Args = []string{"", "sample2", "-s"}
// 	assert.Equal(0, Main())
//
// 	// 画像が生成されていること
// 	_, err := os.Stat(filepath.Join(pictDir, "t.png"))
// 	assert.NoError(err)
// }
