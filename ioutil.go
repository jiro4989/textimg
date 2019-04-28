package main

import (
	"bufio"
	"os"
)

// readStdin は標準入力を文字列の配列として返す。
func readStdin() (ret []string) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		ret = append(ret, line)
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	return
}
