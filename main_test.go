package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inDir  = "testdata/in"
	outDir = "testdata/out"
)

func TestMainNormal(t *testing.T) {
	tests := []struct {
		desc string
		in   []string
		want int
	}{
		// {
		// 	desc: "正常系: FontIndex, EmojiFontIndexを指定できる",
		// 	in:   []string{"", "Sample", "-o", outDir + "/main_font_index.png", "-x", "0", "-X", "0"},
		// 	want: 0,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			os.Args = tt.in
			got := Main()
			assert.Equal(tt.want, got)
		})
	}

	// _, err := os.Stat(outDir + "/main_numbering_2.gif")
	// assert.NoError(t, err)
	//
	// _, err = os.Stat(outDir + "/main_numbering_3.gif")
	// assert.NoError(t, err)
}
