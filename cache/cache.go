package cache

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Cache interface {
}

func DeflateTarGz(tarGzPath, deflateDir string) error {
	tarFile, err := os.Create(tarGzPath)

	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarGzWriter := tar.NewWriter(gzipWriter)
	defer tarGzWriter.Close()

	walkDir(deflateDir, tarGzWriter)

	return nil
}

func InflateTarGz(tarGzPath, inflateDir string) error {
	file, err := os.Open(tarGzPath)

	if err != nil {
		return err
	}
	defer file.Close()

	gzfile, err := gzip.NewReader(file)

	if err != nil {
		return err
	}

	reader := tar.NewReader(gzfile)

	for {
		header, err := reader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		buffer := new(bytes.Buffer)
		outPath := inflateDir + "/" + header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err = os.Stat(outPath); err != nil {
				os.MkdirAll(outPath, 0755)
			}

		case tar.TypeReg:
			if _, err = io.Copy(buffer, reader); err != nil {
				return err
			}

			if err = ioutil.WriteFile(header.Name, buffer.Bytes(), 0755); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkDir(baseDir string, tarGzWriter *tar.Writer) error {
	dir, err := os.Open(baseDir)

	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)

	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		filePath := dir.Name() + string(filepath.Separator) + fileInfo.Name()

		if fileInfo.IsDir() {
			walkDir(filePath, tarGzWriter)
		} else {
			writeTarGz(filePath, tarGzWriter, fileInfo)
		}
	}

	return nil
}

func writeTarGz(filePath string, tarGzWriter *tar.Writer, fileInfo os.FileInfo) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	header := new(tar.Header)
	header.Name = filePath
	header.Size = fileInfo.Size()
	header.Mode = int64(fileInfo.Mode())
	header.ModTime = fileInfo.ModTime()

	err = tarGzWriter.WriteHeader(header)

	if err != nil {
		return err
	}

	_, err = io.Copy(tarGzWriter, file)

	if err != nil {
		return err
	}

	return nil
}
