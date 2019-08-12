package main

import (
	"os"
)

func main() {
	if err := RootCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}
