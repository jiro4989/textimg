package main

import (
	"os"
	"testing"
)

const (
	inDir        = "testdata/in"
	outDir       = "testdata/out"
	outDockerDir = "testdata/out_docker"
)

func TestMain(m *testing.M) {
	testBefore()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func testBefore() {
	if err := os.RemoveAll(outDir); err != nil {
		panic(err)
	}
	// nolint
	os.Mkdir(outDir, os.ModePerm)
	// nolint
	os.Mkdir(outDockerDir, os.ModePerm)
}
