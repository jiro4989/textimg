package main

import (
	"fmt"
	"strings"
)

func main() {
	// alphabet
	const startAlphabetCodePoint rune = 'A'
	const endAlphabetCodePoint rune = 'z'
	var sb strings.Builder
	appendRangeStrings(&sb, startAlphabetCodePoint, endAlphabetCodePoint)

	// emoji
	const startEmojiCodePoint rune = 0x1f600
	const endEmojiCodePoint = startEmojiCodePoint + 1000
	appendRangeStrings(&sb, startEmojiCodePoint, endEmojiCodePoint)
	fmt.Println(sb.String())
}

func appendRangeStrings(sb *strings.Builder, s, e rune) {
	counter := 1
	for i := s; i <= e; i++ {
		r := []rune{i}
		s := string(r)
		sb.WriteString(s)
		if counter%50 == 0 {
			sb.WriteString("\n")
		}
		counter++
	}
	sb.WriteString("\n")
}
