package util

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"

	"github.com/pierrec/lz4/v4"
)

var Lz4Levels = map[int]lz4.CompressionLevel{
	0: lz4.Fast,
	1: lz4.Level1,
	2: lz4.Level2,
	3: lz4.Level3,
	4: lz4.Level4,
	5: lz4.Level5,
	6: lz4.Level6,
	7: lz4.Level7,
	8: lz4.Level8,
	9: lz4.Level9,
}

func CompressLz4(source io.Reader, out io.Writer, level int) error {
	writer := lz4.NewWriter(out)
	err := writer.Apply(
		lz4.CompressionLevelOption(Lz4Levels[level]),
	)
	if err != nil {
		return err
	}

	_, err = writer.ReadFrom(source)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func PackageTarFS(source fs.FS, writer io.Writer) error {
	tw := tar.NewWriter(writer)
	if err := tw.AddFS(source); err != nil {
		return err
	}
	return tw.Close()
}

func PackageTar(source string, writer io.Writer) error {
	tw := tar.NewWriter(writer)
	if err := tw.AddFS(os.DirFS(source)); err != nil {
		return err
	}
	return tw.Close()
}
