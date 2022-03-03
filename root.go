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
	appconf config.Config
	envvars config.EnvVars
)

func init() {
	envvars = config.NewEnvVars()
	cobra.OnInitialize()

	RootCommand.Flags().SortFlags = false
	RootCommand.Flags().StringVarP(&appconf.Foreground, "foreground", "g", "white", `foreground text color.
available color types are [black|red|green|yellow|blue|magenta|cyan|white]
or (R,G,B,A(0~255))`)
	RootCommand.Flags().StringVarP(&appconf.Background, "background", "b", "black", `background text color.
color types are same as "foreground" option`)

	var font string
	envFontFile := envvars.FontFile
	if envFontFile != "" {
		font = envFontFile
	}
	RootCommand.Flags().StringVarP(&appconf.FontFile, "fontfile", "f", font, `font file path.
You can change this default value with environment variables TEXTIMG_FONT_FILE`)
	RootCommand.Flags().IntVarP(&appconf.FontIndex, "fontindex", "x", 0, "")
	appconf.SetFontFileAndFontIndex(runtime.GOOS)

	envEmojiFontFile := envvars.EmojiFontFile
	RootCommand.Flags().StringVarP(&appconf.EmojiFontFile, "emoji-fontfile", "e", envEmojiFontFile, "emoji font file")
	RootCommand.Flags().IntVarP(&appconf.EmojiFontIndex, "emoji-fontindex", "X", 0, "")

	RootCommand.Flags().BoolVarP(&appconf.UseEmojiFont, "use-emoji-font", "i", false, "use emoji font")
	RootCommand.Flags().BoolVarP(&appconf.UseShellgeiEmojiFontfile, "shellgei-emoji-fontfile", "z", false, `emoji font file for shellgei-bot (path: "`+config.ShellgeiEmojiFontPath+`")`)

	RootCommand.Flags().IntVarP(&appconf.FontSize, "fontsize", "F", 20, "font size")
	RootCommand.Flags().StringVarP(&appconf.Outpath, "out", "o", "", `output image file path.
available image formats are [png | jpg | gif]`)
	RootCommand.Flags().BoolVarP(&appconf.AddTimeStamp, "timestamp", "t", false, `add time stamp to output image file path.`)
	RootCommand.Flags().BoolVarP(&appconf.SaveNumberedFile, "numbered", "n", false, `add number-suffix to filename when the output file was existed.
ex: t_2.png`)
	RootCommand.Flags().BoolVarP(&appconf.UseShellgeiImagedir, "shellgei-imagedir", "s", false, `image directory path (path: "$HOME/Pictures/t.png" or "$TEXTIMG_OUTPUT_DIR/t.png")`)

	RootCommand.Flags().BoolVarP(&appconf.UseAnimation, "animation", "a", false, "generate animation gif")
	RootCommand.Flags().IntVarP(&appconf.Delay, "delay", "d", 20, "animation delay time")
	RootCommand.Flags().IntVarP(&appconf.LineCount, "line-count", "l", 1, "animation input line count")
	RootCommand.Flags().BoolVarP(&appconf.UseSlideAnimation, "slide", "S", false, "use slide animation")
	RootCommand.Flags().IntVarP(&appconf.SlideWidth, "slide-width", "W", 1, "sliding animation width")
	RootCommand.Flags().BoolVarP(&appconf.SlideForever, "forever", "E", false, "sliding forever")
	RootCommand.Flags().BoolVarP(&appconf.PrintEnvironments, "environments", "", false, "print environment variables")
	RootCommand.Flags().BoolVarP(&appconf.ToSlackIcon, "slack", "", false, "resize to slack icon size (128x128 px)")
	RootCommand.Flags().IntVarP(&appconf.ResizeWidth, "resize-width", "", 0, "resize width")
	RootCommand.Flags().IntVarP(&appconf.ResizeHeight, "resize-height", "", 0, "resize height")
}

var RootCommand = &cobra.Command{
	Use:     global.AppName,
	Short:   global.AppName + " is command to convert from colored text (ANSI or 256) to image.",
	Example: global.AppName + ` $'\x1b[31mRED\x1b[0m' -o out.png`,
	Version: global.Version,
	RunE:    runRootCommand,
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	if appconf.PrintEnvironments {
		for _, envName := range global.EnvNames {
			text := fmt.Sprintf("%s=%s", envName, os.Getenv(envName))
			fmt.Println(text)
		}
		return nil
	}

	if err := appconf.Adjust(args, envvars); err != nil {
		return err
	}
	defer appconf.Writer.Close()

	ls := appconf.Tokens.StringLines()
	param := &image.ImageParam{
		BaseWidth:          parser.StringWidth(ls),
		BaseHeight:         len(ls),
		ForegroundColor:    c.RGBA(appconf.ForegroundColor),
		BackgroundColor:    c.RGBA(appconf.BackgroundColor),
		FontFace:           appconf.FontFace,
		EmojiFontFace:      appconf.EmojiFontFace,
		EmojiDir:           appconf.EmojiDir,
		FontSize:           appconf.FontSize,
		Delay:              appconf.Delay,
		AnimationLineCount: appconf.LineCount,
		ResizeWidth:        appconf.ResizeWidth,
		ResizeHeight:       appconf.ResizeHeight,
	}
	img := image.NewImage(param)
	if err := img.Draw(appconf.Tokens); err != nil {
		return err
	}
	if err := img.Encode(appconf.Writer, appconf.FileExtension); err != nil {
		return err
	}

	return nil
}
