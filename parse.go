package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// classEscape はテキストに付与されているエスケープシーケンス
	classEscape ClassString = iota
	// classText はただのテキスト
	classText
)

var reANSIColorEscapeSequence *regexp.Regexp

type (
	ClassString      int
	ClassifiedString struct {
		class ClassString
		text  string
	}
	ClassifiedStrings []ClassifiedString
)

func init() {
	reANSIColorEscapeSequence = regexp.MustCompile(`[34][0-7]`)
}

// onlyText はcs.textのみを結合して返却する。
func (cs ClassifiedStrings) onlyText() (ret string) {
	for _, v := range cs {
		if v.class == classText {
			ret += v.text
		}
	}
	return
}

// parseText はテキストを解析しエスケープシーケンスにマッチした箇所と色を返す
// マッチするものが存在しなかった場合は、次のエスケープシーケンスが出現する場所ま
// での文字列を返す。
// エスケープ文字が全く出てこなければ、全部をmatchedとして返す。
//
// 前提として、色コードとは全く関係のない文字列は削除しておく必要がある。
// See also: removeNotColorEscapeSequences
func parseText(s string) (string, string, string) {
	col := getOnlyColorEscapeSequence(s)
	// エスケープ文字自体は返す文字列に含めないため削除する
	headPos := 0
	if col != colorEscapeSequenceNone {
		headPos = len(col)
	}
	s = s[headPos:]

	// 次のエスケープシーケンスが見つかるまでをmatchedとする
	// 何も見つからなければ全部を返す
	// _, idx := getColorPos(s)
	for _, searchWord := range []string{"\x1b[3", "\x1b[4", "\x1b[0"} {
		idx := strings.Index(s, searchWord)
		if idx != -1 {
			return col, s[:idx], s[idx:]
		}
	}
	return col, s, ""
}

// getOnlyColorEscapeSequence はエスケープ文字のうち、色に関連のあるエスケープ文
// 字と通常のテキストのみを残して返却する。
func getOnlyColorEscapeSequence(s string) string {
	const pref = "\x1b["
	if !strings.HasPrefix(s, pref) {
		return ""
	}

	var esc string
	for _, v := range s[len(pref):] {
		if v == 'm' {
			break
		}
		esc += string(v)
	}

	for _, v := range strings.Split(esc, ";") {
		if reANSIColorEscapeSequence.MatchString(v) {
			return fmt.Sprintf("\x1b[%sm", v)
		}
	}

	for _, v := range strings.Split(esc, ";") {
		if v == "0" || v == "01" {
			return colorEscapeSequenceReset
		}
	}

	return ""
}

// 色エスケープシーケンス以外のエスケープシーケンスは不要なので削除して返す
func removeNotColorEscapeSequences(s string) (ret string) {
	// エスケースシーケンス部とテキスト部のスライスに分割する
	// 例: ["\x1b[01;31m", "Red", "\x1b[0m", "\x1b[0Km", "Bold"]
	cs := classifyString(s)

	// 不要な色コード以外のエスケープシーケンスを削除する
	// 例: ["\x1b[31m", "Red", "\x1b[0m", "", "Bold"]
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		if c.class == classEscape {
			fixed := getOnlyColorEscapeSequence(c.text)
			cs[i].text = fixed
		}
	}

	// 文字列をすべて結合しreturn
	for _, v := range cs {
		ret += v.text
	}
	return
}

// 不要な色コード以外のエスケープシーケンスを削除する
// 前: ["\x1b[01;31m", "Red", "\x1b[0m", "\x1b[0Km", "Bold"]
// 後: ["\x1b[31m", "Red", "\x1b[0m", "", "Bold"]
func classifyString(s string) (ret ClassifiedStrings) {
	var matchCnt int
	var text string
	for _, v := range s {
		if matchCnt == 0 && v == '\x1b' {
			if text != "" {
				ret = append(ret, ClassifiedString{class: classText, text: text})
				text = ""
			}
			matchCnt += 1
			text += string(v)
			continue
		}
		if matchCnt == 1 && v == '[' {
			matchCnt += 1
			text += string(v)
			continue
		}
		if matchCnt == 2 && v == 'm' {
			text += string(v)
			ret = append(ret, ClassifiedString{class: classEscape, text: text})
			text = ""
			matchCnt = 0
			continue
		}
		text += string(v)
	}
	if text != "" {
		ret = append(ret, ClassifiedString{class: classText, text: text})
	}
	return
}
