package util

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func DecompressTar(source io.Reader, dest string) error {
	tarReader := tar.NewReader(source)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(path, header.FileInfo().Mode())
			if err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(
				path,
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				header.FileInfo().Mode(),
			)
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}

			err = file.Close()
			if err != nil {
				return err
			}

		}
	}

	return nil
}
