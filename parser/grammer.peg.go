package parser

// Code generated by peg parser/grammer.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleroot
	ruleignore
	rulecolors
	ruletext
	rulecolor
	rulereset_color
	rulestandard_color
	ruleextended_color
	ruleextended_color_256
	ruleextended_color_rgb
	ruleextended_color_prefix
	rulezero
	rulenumber
	ruleprefix
	ruleescape_sequence
	rulecolor_suffix
	rulenon_color_suffix
	ruledelimiter
	ruleAction0
	rulePegText
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
)

var rul3s = [...]string{
	"Unknown",
	"root",
	"ignore",
	"colors",
	"text",
	"color",
	"reset_color",
	"standard_color",
	"extended_color",
	"extended_color_256",
	"extended_color_rgb",
	"extended_color_prefix",
	"zero",
	"number",
	"prefix",
	"escape_sequence",
	"color_suffix",
	"non_color_suffix",
	"delimiter",
	"Action0",
	"PegText",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Parser struct {
	ParserFunc

	Buffer string
	buffer []rune
	rules  [29]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Parser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Parser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Parser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *Parser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Parser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *Parser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func (p *Parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.pushResetColor()
		case ruleAction1:
			p.pushText(text)
		case ruleAction2:
			p.pushResetColor()
		case ruleAction3:
			p.pushStandardColorWithCategory(text)
		case ruleAction4:
			p.pushExtendedColor256(text)
		case ruleAction5:
			p.pushExtendedColorRGB()
		case ruleAction6:
			p.setExtendedColorR(text)
		case ruleAction7:
			p.setExtendedColorG(text)
		case ruleAction8:
			p.setExtendedColorB(text)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func Pretty(pretty bool) func(*Parser) error {
	return func(p *Parser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*Parser) error {
	return func(p *Parser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *Parser) Init(options ...func(*Parser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 root <- <(ignore / colors / text)*> */
		func() bool {
			{
				position1 := position
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					{
						position4, tokenIndex4 := position, tokenIndex
						if !_rules[ruleignore]() {
							goto l5
						}
						goto l4
					l5:
						position, tokenIndex = position4, tokenIndex4
						if !_rules[rulecolors]() {
							goto l6
						}
						goto l4
					l6:
						position, tokenIndex = position4, tokenIndex4
						if !_rules[ruletext]() {
							goto l3
						}
					}
				l4:
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				add(ruleroot, position1)
			}
			return true
		},
		/* 1 ignore <- <(prefix number non_color_suffix)> */
		func() bool {
			position7, tokenIndex7 := position, tokenIndex
			{
				position8 := position
				if !_rules[ruleprefix]() {
					goto l7
				}
				if !_rules[rulenumber]() {
					goto l7
				}
				if !_rules[rulenon_color_suffix]() {
					goto l7
				}
				add(ruleignore, position8)
			}
			return true
		l7:
			position, tokenIndex = position7, tokenIndex7
			return false
		},
		/* 2 colors <- <((prefix color_suffix Action0) / (prefix color (delimiter color)* color_suffix))> */
		func() bool {
			position9, tokenIndex9 := position, tokenIndex
			{
				position10 := position
				{
					position11, tokenIndex11 := position, tokenIndex
					if !_rules[ruleprefix]() {
						goto l12
					}
					if !_rules[rulecolor_suffix]() {
						goto l12
					}
					if !_rules[ruleAction0]() {
						goto l12
					}
					goto l11
				l12:
					position, tokenIndex = position11, tokenIndex11
					if !_rules[ruleprefix]() {
						goto l9
					}
					if !_rules[rulecolor]() {
						goto l9
					}
				l13:
					{
						position14, tokenIndex14 := position, tokenIndex
						if !_rules[ruledelimiter]() {
							goto l14
						}
						if !_rules[rulecolor]() {
							goto l14
						}
						goto l13
					l14:
						position, tokenIndex = position14, tokenIndex14
					}
					if !_rules[rulecolor_suffix]() {
						goto l9
					}
				}
			l11:
				add(rulecolors, position10)
			}
			return true
		l9:
			position, tokenIndex = position9, tokenIndex9
			return false
		},
		/* 3 text <- <(<.+> Action1)> */
		func() bool {
			position15, tokenIndex15 := position, tokenIndex
			{
				position16 := position
				{
					position17 := position
					if !matchDot() {
						goto l15
					}
				l18:
					{
						position19, tokenIndex19 := position, tokenIndex
						if !matchDot() {
							goto l19
						}
						goto l18
					l19:
						position, tokenIndex = position19, tokenIndex19
					}
					add(rulePegText, position17)
				}
				if !_rules[ruleAction1]() {
					goto l15
				}
				add(ruletext, position16)
			}
			return true
		l15:
			position, tokenIndex = position15, tokenIndex15
			return false
		},
		/* 4 color <- <(reset_color / standard_color / extended_color)> */
		func() bool {
			position20, tokenIndex20 := position, tokenIndex
			{
				position21 := position
				{
					position22, tokenIndex22 := position, tokenIndex
					if !_rules[rulereset_color]() {
						goto l23
					}
					goto l22
				l23:
					position, tokenIndex = position22, tokenIndex22
					if !_rules[rulestandard_color]() {
						goto l24
					}
					goto l22
				l24:
					position, tokenIndex = position22, tokenIndex22
					if !_rules[ruleextended_color]() {
						goto l20
					}
				}
			l22:
				add(rulecolor, position21)
			}
			return true
		l20:
			position, tokenIndex = position20, tokenIndex20
			return false
		},
		/* 5 reset_color <- <('0'+ Action2)> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				if buffer[position] != rune('0') {
					goto l25
				}
				position++
			l27:
				{
					position28, tokenIndex28 := position, tokenIndex
					if buffer[position] != rune('0') {
						goto l28
					}
					position++
					goto l27
				l28:
					position, tokenIndex = position28, tokenIndex28
				}
				if !_rules[ruleAction2]() {
					goto l25
				}
				add(rulereset_color, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 6 standard_color <- <(zero <(('3' / '4' / '9' / ('1' '0')) [0-7])> Action3)> */
		func() bool {
			position29, tokenIndex29 := position, tokenIndex
			{
				position30 := position
				if !_rules[rulezero]() {
					goto l29
				}
				{
					position31 := position
					{
						position32, tokenIndex32 := position, tokenIndex
						if buffer[position] != rune('3') {
							goto l33
						}
						position++
						goto l32
					l33:
						position, tokenIndex = position32, tokenIndex32
						if buffer[position] != rune('4') {
							goto l34
						}
						position++
						goto l32
					l34:
						position, tokenIndex = position32, tokenIndex32
						if buffer[position] != rune('9') {
							goto l35
						}
						position++
						goto l32
					l35:
						position, tokenIndex = position32, tokenIndex32
						if buffer[position] != rune('1') {
							goto l29
						}
						position++
						if buffer[position] != rune('0') {
							goto l29
						}
						position++
					}
				l32:
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l29
					}
					position++
					add(rulePegText, position31)
				}
				if !_rules[ruleAction3]() {
					goto l29
				}
				add(rulestandard_color, position30)
			}
			return true
		l29:
			position, tokenIndex = position29, tokenIndex29
			return false
		},
		/* 7 extended_color <- <(extended_color_256 / extended_color_rgb)> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				{
					position38, tokenIndex38 := position, tokenIndex
					if !_rules[ruleextended_color_256]() {
						goto l39
					}
					goto l38
				l39:
					position, tokenIndex = position38, tokenIndex38
					if !_rules[ruleextended_color_rgb]() {
						goto l36
					}
				}
			l38:
				add(ruleextended_color, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 8 extended_color_256 <- <(extended_color_prefix delimiter zero '5' delimiter <number> Action4)> */
		func() bool {
			position40, tokenIndex40 := position, tokenIndex
			{
				position41 := position
				if !_rules[ruleextended_color_prefix]() {
					goto l40
				}
				if !_rules[ruledelimiter]() {
					goto l40
				}
				if !_rules[rulezero]() {
					goto l40
				}
				if buffer[position] != rune('5') {
					goto l40
				}
				position++
				if !_rules[ruledelimiter]() {
					goto l40
				}
				{
					position42 := position
					if !_rules[rulenumber]() {
						goto l40
					}
					add(rulePegText, position42)
				}
				if !_rules[ruleAction4]() {
					goto l40
				}
				add(ruleextended_color_256, position41)
			}
			return true
		l40:
			position, tokenIndex = position40, tokenIndex40
			return false
		},
		/* 9 extended_color_rgb <- <(extended_color_prefix delimiter zero '2' Action5 delimiter <number> Action6 delimiter <number> Action7 delimiter <number> Action8)> */
		func() bool {
			position43, tokenIndex43 := position, tokenIndex
			{
				position44 := position
				if !_rules[ruleextended_color_prefix]() {
					goto l43
				}
				if !_rules[ruledelimiter]() {
					goto l43
				}
				if !_rules[rulezero]() {
					goto l43
				}
				if buffer[position] != rune('2') {
					goto l43
				}
				position++
				if !_rules[ruleAction5]() {
					goto l43
				}
				if !_rules[ruledelimiter]() {
					goto l43
				}
				{
					position45 := position
					if !_rules[rulenumber]() {
						goto l43
					}
					add(rulePegText, position45)
				}
				if !_rules[ruleAction6]() {
					goto l43
				}
				if !_rules[ruledelimiter]() {
					goto l43
				}
				{
					position46 := position
					if !_rules[rulenumber]() {
						goto l43
					}
					add(rulePegText, position46)
				}
				if !_rules[ruleAction7]() {
					goto l43
				}
				if !_rules[ruledelimiter]() {
					goto l43
				}
				{
					position47 := position
					if !_rules[rulenumber]() {
						goto l43
					}
					add(rulePegText, position47)
				}
				if !_rules[ruleAction8]() {
					goto l43
				}
				add(ruleextended_color_rgb, position44)
			}
			return true
		l43:
			position, tokenIndex = position43, tokenIndex43
			return false
		},
		/* 10 extended_color_prefix <- <(zero ('3' / '4') '8')> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				if !_rules[rulezero]() {
					goto l48
				}
				{
					position50, tokenIndex50 := position, tokenIndex
					if buffer[position] != rune('3') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex = position50, tokenIndex50
					if buffer[position] != rune('4') {
						goto l48
					}
					position++
				}
			l50:
				if buffer[position] != rune('8') {
					goto l48
				}
				position++
				add(ruleextended_color_prefix, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 11 zero <- <'0'*> */
		func() bool {
			{
				position53 := position
			l54:
				{
					position55, tokenIndex55 := position, tokenIndex
					if buffer[position] != rune('0') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex = position55, tokenIndex55
				}
				add(rulezero, position53)
			}
			return true
		},
		/* 12 number <- <[0-9]+> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l56
				}
				position++
			l58:
				{
					position59, tokenIndex59 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex = position59, tokenIndex59
				}
				add(rulenumber, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 13 prefix <- <(escape_sequence '[')> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				if !_rules[ruleescape_sequence]() {
					goto l60
				}
				if buffer[position] != rune('[') {
					goto l60
				}
				position++
				add(ruleprefix, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
			return false
		},
		/* 14 escape_sequence <- <'\x1b'> */
		func() bool {
			position62, tokenIndex62 := position, tokenIndex
			{
				position63 := position
				if buffer[position] != rune('\x1b') {
					goto l62
				}
				position++
				add(ruleescape_sequence, position63)
			}
			return true
		l62:
			position, tokenIndex = position62, tokenIndex62
			return false
		},
		/* 15 color_suffix <- <'m'> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				if buffer[position] != rune('m') {
					goto l64
				}
				position++
				add(rulecolor_suffix, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 16 non_color_suffix <- <([A-H] / 'f' / 'S' / 'T' / 'J' / 'K')> */
		func() bool {
			position66, tokenIndex66 := position, tokenIndex
			{
				position67 := position
				{
					position68, tokenIndex68 := position, tokenIndex
					if c := buffer[position]; c < rune('A') || c > rune('H') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex = position68, tokenIndex68
					if buffer[position] != rune('f') {
						goto l70
					}
					position++
					goto l68
				l70:
					position, tokenIndex = position68, tokenIndex68
					if buffer[position] != rune('S') {
						goto l71
					}
					position++
					goto l68
				l71:
					position, tokenIndex = position68, tokenIndex68
					if buffer[position] != rune('T') {
						goto l72
					}
					position++
					goto l68
				l72:
					position, tokenIndex = position68, tokenIndex68
					if buffer[position] != rune('J') {
						goto l73
					}
					position++
					goto l68
				l73:
					position, tokenIndex = position68, tokenIndex68
					if buffer[position] != rune('K') {
						goto l66
					}
					position++
				}
			l68:
				add(rulenon_color_suffix, position67)
			}
			return true
		l66:
			position, tokenIndex = position66, tokenIndex66
			return false
		},
		/* 17 delimiter <- <';'> */
		func() bool {
			position74, tokenIndex74 := position, tokenIndex
			{
				position75 := position
				if buffer[position] != rune(';') {
					goto l74
				}
				position++
				add(ruledelimiter, position75)
			}
			return true
		l74:
			position, tokenIndex = position74, tokenIndex74
			return false
		},
		/* 19 Action0 <- <{ p.pushResetColor() }> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 21 Action1 <- <{ p.pushText(text) }> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 22 Action2 <- <{ p.pushResetColor() }> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 23 Action3 <- <{ p.pushStandardColorWithCategory(text) }> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 24 Action4 <- <{ p.pushExtendedColor256(text) }> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 25 Action5 <- <{ p.pushExtendedColorRGB() }> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 26 Action6 <- <{ p.setExtendedColorR(text) }> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 27 Action7 <- <{ p.setExtendedColorG(text) }> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 28 Action8 <- <{ p.setExtendedColorB(text) }> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
	}
	p.rules = _rules
	return nil
}
