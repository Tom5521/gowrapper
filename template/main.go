package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/gowrapper/util"
	"github.com/ncruces/zenity"
	"github.com/pierrec/lz4/v4"
)

//go:embed package.tar.lz4
var embeddedData []byte

var (
	AppName    string
	BinaryPath string
	Args       string
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			zenity.Error(fmt.Sprint(r))
		}
	}()

	tempDir, err := os.MkdirTemp(os.TempDir(), AppName)
	if err != nil {
		panic(err)
	}
	bytesReader := bytes.NewReader(embeddedData)
	lz4Reader := lz4.NewReader(bytesReader)
	err = util.DecompressTar(lz4Reader, tempDir)
	if err != nil {
		panic(err)
	}

	binary := filepath.Join(tempDir, BinaryPath)
	cmd := exec.Command(
		binary,
		slices.Concat(
			strings.Split(Args, " "),
			os.Args,
		)...,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		slog.Error(err.Error())
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		panic(err)
	}
}
