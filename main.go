package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
)

func main() {
	if err := RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func maxStringWidth(s []string) (max int) {
	for _, v := range s {
		text := classifyString(v).OnlyText()
		width := runewidth.StringWidth(text)
		if max < width {
			max = width
		}
	}
	return
}
