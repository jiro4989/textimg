package parser

import (
	"strconv"

	"github.com/jiro4989/textimg/v3/color"
	"github.com/jiro4989/textimg/v3/token"
)

type ParserFunc struct {
	// pegが生成するTokensと名前が衝突するので別名にする
	Tk token.Tokens
}

func Parse(s string) (token.Tokens, error) {
	p := &Parser{Buffer: s}
	if err := p.Init(); err != nil {
		return nil, err
	}
	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()
	return p.Tk, nil
}

func (p *ParserFunc) pushResetColor() {
	p.Tk = append(p.Tk, token.NewResetColor())
}

func (p *ParserFunc) pushResetForegroundColor() {
	p.Tk = append(p.Tk, token.NewResetForegroundColor())
}

func (p *ParserFunc) pushResetBackgroundColor() {
	p.Tk = append(p.Tk, token.NewResetBackgroundColor())
}

func (p *ParserFunc) pushReverseColor() {
	p.Tk = append(p.Tk, token.NewReverseColor())
}

func (p *ParserFunc) pushText(text string) {
	p.Tk = append(p.Tk, token.NewText(text))
}

func (p *ParserFunc) pushStandardColorWithCategory(text string) {
	p.Tk = append(p.Tk, token.NewStandardColorWithCategory(text))
}

func (p *ParserFunc) pushExtendedColor(text string) {
	p.Tk = append(p.Tk, token.NewExtendedColor(text))
}

func (p *ParserFunc) setExtendedColor256(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Tk[len(p.Tk)-1].Color = color.Map256[int(n)]
}

func (p *ParserFunc) setExtendedColorR(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Tk[len(p.Tk)-1].Color.R = uint8(n)
}

func (p *ParserFunc) setExtendedColorG(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Tk[len(p.Tk)-1].Color.G = uint8(n)
}

func (p *ParserFunc) setExtendedColorB(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Tk[len(p.Tk)-1].Color.B = uint8(n)
}
