package parser

import (
	"strconv"
)

type ParserFunc struct {
	Colors Colors
}

func Parse(s string) (Colors, error) {
	p := &Parser{Buffer: s}
	p.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()
	return p.Colors, nil
}

func (p *ParserFunc) pushResetColor() {
	p.Colors = append(p.Colors, newResetColor())
}

func (p *ParserFunc) pushText(text string) {
	p.Colors = append(p.Colors, newText(text))
}

func (p *ParserFunc) pushStandardColorWithCategory(text string) {
	p.Colors = append(p.Colors, newStandardColorWithCategory(text))
}

func (p *ParserFunc) pushExtendedColor256(text string) {
	p.Colors = append(p.Colors, newExtendedColor256(text))
}

func (p *ParserFunc) pushExtendedColorRGB() {
	p.Colors = append(p.Colors, newExtendedColorRGB())
}

func (p *ParserFunc) setExtendedColorR(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Colors[len(p.Colors)-1].Color.R = uint8(n)
}

func (p *ParserFunc) setExtendedColorG(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Colors[len(p.Colors)-1].Color.G = uint8(n)
}

func (p *ParserFunc) setExtendedColorB(text string) {
	n, _ := strconv.ParseUint(text, 10, 8)
	p.Colors[len(p.Colors)-1].Color.B = uint8(n)
}
