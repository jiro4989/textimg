package main

import (
	"fmt"
	c "image/color"
	"os"
	"runtime"

	"github.com/jiro4989/textimg/v3/config"
	"github.com/jiro4989/textimg/v3/image"
	"github.com/jiro4989/textimg/v3/internal/global"
	"github.com/jiro4989/textimg/v3/parser"

	"github.com/spf13/cobra"
)

var (
	conf    config.Config
	envvars config.EnvVars
)

func init() {
	envvars = config.NewEnvVars()
	cobra.OnInitialize()

	RootCommand.Flags().SortFlags = false
	RootCommand.Flags().StringVarP(&conf.Foreground, "foreground", "g", "white", `foreground text color.
available color types are [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringVarP(&conf.Background, "background", "b", "black", `background text color.
color types are same as "foreground" option`)

	var font string
	envFontFile := envvars.FontFile
	if envFontFile != "" {
		font = envFontFile
	}
	RootCommand.Flags().StringVarP(&conf.FontFile, "fontfile", "f", font, `font file path.
You can change this default value with environment variables TEXTIMG_FONT_FILE`)
	RootCommand.Flags().IntVarP(&conf.FontIndex, "fontindex", "x", 0, "")
	conf.SetFontFileAndFontIndex(runtime.GOOS)

	envEmojiFontFile := envvars.EmojiFontFile
	RootCommand.Flags().StringVarP(&conf.EmojiFontFile, "emoji-fontfile", "e", envEmojiFontFile, "emoji font file")
	RootCommand.Flags().IntVarP(&conf.EmojiFontIndex, "emoji-fontindex", "X", 0, "")

	RootCommand.Flags().BoolVarP(&conf.UseEmojiFont, "use-emoji-font", "i", false, "use emoji font")
	RootCommand.Flags().BoolVarP(&conf.UseShellgeiEmojiFontfile, "shellgei-emoji-fontfile", "z", false, `emoji font file for shellgei-bot (path: "`+config.ShellgeiEmojiFontPath+`")`)

	RootCommand.Flags().IntVarP(&conf.FontSize, "fontsize", "F", 20, "font size")
	RootCommand.Flags().StringVarP(&conf.Outpath, "out", "o", "", `output image file path.
available image formats are [png | jpg | gif]`)
	RootCommand.Flags().BoolVarP(&conf.AddTimeStamp, "timestamp", "t", false, `add time stamp to output image file path.`)
	RootCommand.Flags().BoolVarP(&conf.SaveNumberedFile, "numbered", "n", false, `add number-suffix to filename when the output file was existed.
ex: t_2.png`)
	RootCommand.Flags().BoolVarP(&conf.UseShellgeiImagedir, "shellgei-imagedir", "s", false, `image directory path (path: "$HOME/Pictures/t.png" or "$TEXTIMG_OUTPUT_DIR/t.png")`)

	RootCommand.Flags().BoolVarP(&conf.UseAnimation, "animation", "a", false, "generate animation gif")
	RootCommand.Flags().IntVarP(&conf.Delay, "delay", "d", 20, "animation delay time")
	RootCommand.Flags().IntVarP(&conf.LineCount, "line-count", "l", 1, "animation input line count")
	RootCommand.Flags().BoolVarP(&conf.UseSlideAnimation, "slide", "S", false, "use slide animation")
	RootCommand.Flags().IntVarP(&conf.SlideWidth, "slide-width", "W", 1, "sliding animation width")
	RootCommand.Flags().BoolVarP(&conf.SlideForever, "forever", "E", false, "sliding forever")
	RootCommand.Flags().BoolVarP(&conf.PrintEnvironments, "environments", "", false, "print environment variables")
	RootCommand.Flags().BoolVarP(&conf.ToSlackIcon, "slack", "", false, "resize to slack icon size (128x128 px)")
	RootCommand.Flags().IntVarP(&conf.ResizeWidth, "resize-width", "", 0, "resize width")
	RootCommand.Flags().IntVarP(&conf.ResizeHeight, "resize-height", "", 0, "resize height")
}

var RootCommand = &cobra.Command{
	Use:     global.AppName,
	Short:   global.AppName + " is command to convert from colored text (ANSI or 256) to image.",
	Example: global.AppName + ` $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: global.Version,
	RunE:    runRootCommand,
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	if conf.PrintEnvironments {
		for _, envName := range global.EnvNames {
			text := fmt.Sprintf("%s=%s", envName, os.Getenv(envName))
			fmt.Println(text)
		}
		return nil
	}

	if err := conf.Adjust(args, envvars); err != nil {
		return err
	}
	defer conf.Writer.Close()

	ls := conf.Tokens.StringLines()
	param := &image.ImageParam{
		BaseWidth:          parser.StringWidth(ls),
		BaseHeight:         len(ls),
		ForegroundColor:    c.RGBA(conf.ForegroundColor),
		BackgroundColor:    c.RGBA(conf.BackgroundColor),
		FontFace:           conf.FontFace,
		EmojiFontFace:      conf.EmojiFontFace,
		EmojiDir:           conf.EmojiDir,
		FontSize:           conf.FontSize,
		Delay:              conf.Delay,
		AnimationLineCount: conf.LineCount,
		ResizeWidth:        conf.ResizeWidth,
		ResizeHeight:       conf.ResizeHeight,
	}
	img := image.NewImage(param)
	if err := img.Draw(conf.Tokens); err != nil {
		return err
	}
	if err := img.Encode(conf.Writer, conf.FileExtension); err != nil {
		return err
	}

	return nil
}
