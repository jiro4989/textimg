package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().StringP("background", "b", "black", "background color")
	RootCommand.Flags().StringP("scale", "s", "auto", "scale size")
	RootCommand.Flags().StringP("out", "o", "", "output path")
	RootCommand.Flags().StringP("fontfile", "f", "", "font file path")
	RootCommand.Flags().IntP("fontsize", "F", 64, "font size")
}

var RootCommand = &cobra.Command{
	Use:     "coltoi",
	Short:   "coltoi is command to convert from ANSI colored text to image.",
	Example: `coltoi $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flags()

		background, err := f.GetString("background")
		if err != nil {
			panic(err)
		}

		scale, err := f.GetString("scale")
		if err != nil {
			panic(err)
		}

		outpath, err := f.GetString("out")
		if err != nil {
			panic(err)
		}

		fontpath, err := f.GetString("fontfile")
		if err != nil {
			panic(err)
		}

		fmt.Println(background, scale, outpath, fontpath)

		// 引数にテキストの指定がなければ標準入力を使用する
		var texts []string
		if len(args) < 1 {
			texts = readStdin()
		} else {
			texts = args
		}

		// 出力先画像の指定がなければ標準出力を出力先にする
		var w *os.File
		if outpath == "" {
			w = os.Stdout
		} else {
			var err error
			w, err = os.Create(outpath)
			if err != nil {
				panic(err)
			}
			defer w.Close()
		}

		writeImage(texts, w)
	},
}
