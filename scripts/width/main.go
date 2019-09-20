// width は文字のrunewidthが返す文字幅を確認するためのツール。

/*

使い方

cd tools/width
go build .
./width あいうえお■漢字abcde😲

*/

package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
)

func main() {
	args := os.Args
	fmt.Println("Char CodePoint Width")
	for _, c := range []rune(args[1]) {
		text := fmt.Sprintf("%v %d %d", string(c), c, runewidth.RuneWidth(c))
		fmt.Println(text)
	}
}
