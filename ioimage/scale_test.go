package ioimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScale(t *testing.T) {
	type TestData struct {
		desc       string
		x, y, w, h int
		wantWidth  int
		wantHeight int
	}
	tds := []TestData{
		{
			desc:       "正常系: wが-1のときはwidthが調整される",
			x:          200,
			y:          100,
			w:          -1,
			h:          200,
			wantWidth:  400,
			wantHeight: 200,
		},
		{
			desc:       "正常系: hが-1のときはheightが調整される",
			x:          200,
			y:          100,
			w:          100,
			h:          -1,
			wantWidth:  100,
			wantHeight: 50,
		},
		{
			desc:       "正常系: hが-1のときはheightが調整される",
			x:          200,
			y:          100,
			w:          100,
			h:          -1,
			wantWidth:  100,
			wantHeight: 50,
		},
		{
			desc:       "正常系: wとhが-1出ないときはwとhが返る",
			x:          200,
			y:          100,
			w:          400,
			h:          300,
			wantWidth:  400,
			wantHeight: 300,
		},
	}
	for _, v := range tds {
		t.Run(v.desc, func(t *testing.T) {
			a := assert.New(t)
			w, h := complementWidthHeight(v.x, v.y, v.w, v.h)
			a.Equal(v.wantWidth, w)
			a.Equal(v.wantHeight, h)
		})
	}
}
