package parser

import (
	"strconv"

	"github.com/jiro4989/textimg/v3/token"
)

type ParserFunc struct {
	// pegが生成するTokensと名前が衝突するので別名にする
	Tk token.Tokens
}

func Parse(s string) (token.Tokens, error) {
	p := &Parser{Buffer: s}
	p.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()
	return p.Tk, nil
}

func (p *ParserFunc) pushResetColor() {
	p.Tk = append(p.Tk, token.NewResetColor())
}

func (p *ParserFunc) pushText(text string) {
	p.Tk = append(p.Tk, token.NewText(text))
}

func (p *ParserFunc) pushStandardColorWithCategory(text string) {
	p.Tk = append(p.Tk, token.NewStandardColorWithCategory(text))
}

func (p *ParserFunc) pushExtendedColor256(text string) {
	p.Tk = append(p.Tk, token.NewExtendedColor256(text))
}

func (p *ParserFunc) pushExtendedColorRGB() {
	p.Tk = append(p.Tk, token.NewExtendedColorRGB())
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