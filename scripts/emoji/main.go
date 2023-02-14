package main

import (
	"fmt"
	"strings"
)

func main() {
	const startEmojiCodePoint rune = 0x1f600
	const endEmojiCodePoint = startEmojiCodePoint + 1000
	var sb strings.Builder
	counter := 1
	for i := startEmojiCodePoint; i < endEmojiCodePoint; i++ {
		r := []rune{i}
		s := string(r)
		sb.WriteString(s)
		if counter%50 == 0 {
			sb.WriteString("\n")
		}
		counter++
	}
	fmt.Println(sb.String())
}
