package main

import (
	"os"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	if err := RootCommand.Execute(); err != nil {
		return -1
	}
	return 0
}
