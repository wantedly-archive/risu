package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func deflateTarGz(tarGzPath, deflateDir string) {
	tarFile, err := os.Create(tarGzPath)

	if err != nil {
		log.Fatal(err)
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarGzWriter := tar.NewWriter(gzipWriter)
	defer tarGzWriter.Close()

	walkDir(deflateDir, tarGzWriter)
}

func inflateTarGz(tarGzPath, inflateDir string) {
	file, err := os.Open(tarGzPath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gzfile, err := gzip.NewReader(file)

	if err != nil {
		log.Fatal(err)
	}

	reader := tar.NewReader(gzfile)

	for {
		header, err := reader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}

			if err = ioutil.WriteFile(header.Name, buffer.Bytes(), 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func walkDir(baseDir string, tarGzWriter *tar.Writer) {
	dir, err := os.Open(baseDir)

	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0)

	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range files {
		filePath := dir.Name() + string(filepath.Separator) + fileInfo.Name()

		if fileInfo.IsDir() {
			walkDir(filePath, tarGzWriter)
		} else {
			writeTarGz(filePath, tarGzWriter, fileInfo)
		}
	}
}

func writeTarGz(filePath string, tarGzWriter *tar.Writer, fileInfo os.FileInfo) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	header := new(tar.Header)
	header.Name = filePath
	header.Size = fileInfo.Size()
	header.Mode = int64(fileInfo.Mode())
	header.ModTime = fileInfo.ModTime()

	err = tarGzWriter.WriteHeader(header)

	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(tarGzWriter, file)

	if err != nil {
		log.Fatal(err)
	}
}
