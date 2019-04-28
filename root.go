package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().StringP("background", "b", "black", "background color")
	RootCommand.Flags().StringP("scale", "s", "auto", "scale size")
	RootCommand.Flags().StringP("out", "o", "", "output path")
	RootCommand.Flags().StringP("font", "f", "", "font file path")
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

		fontpath, err := f.GetString("font")
		if err != nil {
			panic(err)
		}

		fmt.Println(background, scale, outpath, fontpath)
	},
}
