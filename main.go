package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Tom5521/gowrapper/util"
	"github.com/spf13/cobra"
)

//go:embed template/main.go
var template []byte

var root = cobra.Command{
	Use: "gowrapper",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if CompressionLevel > 9 {
			return errors.New("invalid compression level")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		slog.Info("Creating temporary directory...")
		tmpDir, err := os.MkdirTemp(
			os.TempDir(),
			"gowrapper-",
		)
		if err != nil {
			return err
		}

		err = os.WriteFile(
			filepath.Join(tmpDir, "main.go"),
			template,
			os.ModePerm,
		)
		if err != nil {
			return err
		}

		slog.Info("Compressing bundle...")

		tarBuffer := &bytes.Buffer{}
		err = util.PackageTar(BundlePath, tarBuffer)
		if err != nil {
			return err
		}
		// Reset the buffer counter.
		tarBuffer = bytes.NewBuffer(tarBuffer.Bytes())
		var compressedBuffer bytes.Buffer
		err = util.CompressLz4(tarBuffer, &compressedBuffer, CompressionLevel)
		if err != nil {
			return err
		}

		err = os.WriteFile(
			filepath.Join(tmpDir, "package.tar.lz4"),
			compressedBuffer.Bytes(),
			os.ModePerm,
		)
		if err != nil {
			return err
		}

		cmd := newCmd("go", "mod", "init", "template")
		cmd.Dir = tmpDir

		slog.Info("Configuring go modules...")
		err = cmd.Run()
		if err != nil {
			return err
		}
		cmd = newCmd("go", "mod", "tidy")
		cmd.Dir = tmpDir
		err = cmd.Run()
		if err != nil {
			return err
		}

		fullOut, err := filepath.Abs(Output)
		if err != nil {
			return err
		}

		var ldflags string
		if WindowsGUI {
			ldflags += "-H=windowsgui "
		}
		if AppName != "" {
			ldflags += fmt.Sprintf(
				"-X main.AppName=%s ",
				AppName,
			)
		}
		ldflags += fmt.Sprintf("-X main.BinaryPath=%s ",
			BinaryName,
		)
		ldflags += fmt.Sprintf("-X main.Args=%s",
			strings.Join(DefaultArgs, " "),
		)

		cmd = newCmd("go",
			"build",
			"-o",
			fullOut,
			"-ldflags",
			ldflags,
		)
		cmd.Dir = tmpDir

		if Verbose {
			cmd.Args = append(cmd.Args, "-v")
		}
		cmd.Args = append(cmd.Args, GoArgs...)
		cmd.Args = append(cmd.Args, "main.go")

		slog.Info("Building binary...")
		err = cmd.Run()
		if err != nil {
			return err
		}

		slog.Info("Deleting temporary directory...")
		err = os.RemoveAll(tmpDir)
		if err != nil {
			return err
		}

		return nil
	},
}

func newCmd(bin string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd
}

var (
	AppName          string
	DefaultArgs      []string
	BinaryName       string
	BundlePath       string
	Output           string
	Verbose          bool
	CompressionLevel int
	GoArgs           []string
	WindowsGUI       bool
)

func init() {
	flags := root.Flags()
	flags.BoolVar(&WindowsGUI, "windowsgui", false,
		"Adds the flag -H=windowsgui into ldflags",
	)
	flags.BoolVarP(&Verbose, "verbose", "v", false,
		"Set's the -v flag of the go compiler.",
	)
	flags.StringVarP(&AppName, "name", "n", "",
		"Set's the application name to work with.",
	)
	flags.StringVarP(&BundlePath, "bundle", "b", "",
		"Specifies the directory of the binary with all dependencies "+
			"and it's environment.",
	)
	flags.StringVar(&BinaryName, "bin", "",
		"Specifies the binary path inside the bundle.",
	)
	flags.IntVarP(&CompressionLevel, "compression-level", "c", 9,
		"Specifies the lz4 compression level for the bundle embedded into"+
			"the binary [0-9], 0 = fast",
	)
	flags.StringVarP(&Output, "out", "o", "",
		"Set's the output file for the binary.",
	)
	flags.StringSliceVarP(&DefaultArgs, "args", "a", nil,
		"Specifies the default arguments to run the bundled binary.",
	)
	flags.StringSliceVarP(&GoArgs, "go-args", "g", nil,
		"Specifies extra arguments for the compiler.",
	)

	root.MarkFlagRequired("bundle")
	root.MarkFlagRequired("bin")
	root.MarkFlagRequired("out")
}

func main() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
		return
	}
}
