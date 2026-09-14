package global

import "runtime/debug"

const (
	AppName = "textimg"
)

var (
	Version = "dev"
)

func init() {
	const copyright = `
Copyright (c) 2019 jiro4989
Released under the MIT License.
https://github.com/jiro4989/textimg`

	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		Version = v + copyright
	}
}
