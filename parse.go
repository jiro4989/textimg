package main

import (
	"strings"
)

// parseText はテキストを解析しAnsiiEscapeSequenceにマッチした箇所と色を返す
// マッチするものが存在しなかった場合は、次のAnsiiEscapeSequenceが出現する場所ま
// での文字列を返す。
// エスケープ文字が全く出てこなければ、全部をmatchedとして返す。
func parseText(s string) (string, string, string) {
	col := getPrefixColor(s)
	// エスケープ文字自体は返す文字列に含めないため削除する
	headPos := 0
	if col != colorNone {
		headPos = len(col)
	}
	s = s[headPos:]

	// 次のエスケープシーケンスが見つかるまでをmatchedとする
	// 何も見つからなければ全部を返す
	idx := strings.Index(s, "\x1b[0m")
	if idx == -1 {
		return col, s, ""
	}
	return col, s[:idx], s[idx:]
}

func getPrefixColor(s string) string {
	for _, v := range colors {
		if strings.HasPrefix(s, v) {
			return v
		}
	}
	return colorNone
}
