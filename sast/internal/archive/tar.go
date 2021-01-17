package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func Unpack(r io.Reader, destDir string) error {
	z, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(z)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path.Join(destDir, header.Name), 0755); err != nil {
				return fmt.Errorf("failed to create dir: %s", err)
			}
		case tar.TypeReg:
			data := make([]byte, header.Size)
			if _, err := io.ReadFull(tarReader, data); err != nil {
				return fmt.Errorf("failed to read tar: %s", err)
			}
			if err := ioutil.WriteFile(path.Join(destDir, header.Name), data, 0644); err != nil {
				return fmt.Errorf("failed to write file: %s", err)
			}
		default:
			return fmt.Errorf("invalid tar type: %c for file %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
