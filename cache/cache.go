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
	Get(key string) (string, error)
	Put(key, directory string) error
}

func NewCache(backend string) Cache {
	switch backend {
	case "local":
		return NewLocalFsCache()
	case "s3":
		return NewS3Cache()
	default:
		return NewLocalFsCache()
	}
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

	deflateDir, err = filepath.Abs(deflateDir)

	if err != nil {
		return err
	}

	walkDir(deflateDir, deflateDir, tarGzWriter)

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

	inflateDir, err = filepath.Abs(inflateDir)

	if err != nil {
		return err
	}

	if _, err = os.Stat(inflateDir); err != nil {
		os.MkdirAll(inflateDir, 0755)
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
		outPath := inflateDir + string(filepath.Separator) + header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err = os.Stat(outPath); err != nil {
				os.MkdirAll(outPath, 0755)
			}

		case tar.TypeReg, tar.TypeRegA:
			if _, err = io.Copy(buffer, reader); err != nil {
				return err
			}

			if err = ioutil.WriteFile(outPath, buffer.Bytes(), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkDir(walkingDir, baseDir string, tarGzWriter *tar.Writer) error {
	dir, err := os.Open(walkingDir)

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
			walkDir(filePath, baseDir, tarGzWriter)
		} else {
			writeTarGz(filePath, baseDir, tarGzWriter, fileInfo)
		}
	}

	return nil
}

func writeTarGz(filePath, baseDir string, tarGzWriter *tar.Writer, fileInfo os.FileInfo) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	relativePath, err := filepath.Rel(baseDir, filePath)

	if err != nil {
		return err
	}

	header := new(tar.Header)
	header.Name = relativePath
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

func cacheFilePath(cacheDir, key string) string {
	return cacheDir + string(filepath.Separator) + key + ".tar.gz"
}

func inflateDirPath(cacheDir, key string) string {
	return cacheDir + string(filepath.Separator) + key
}
