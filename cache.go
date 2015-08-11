package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func inflateTarGz(tarGzPath, outDir string) {
	file, err := os.Open(tarGzPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	gzfile, err := gzip.NewReader(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := tar.NewReader(gzfile)

	for {
		header, err := reader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		buffer := new(bytes.Buffer)
		outPath := outDir + "/" + header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err = os.Stat(outPath); err != nil {
				os.MkdirAll(outPath, 0755)
			}

		case tar.TypeReg:
			if _, err = io.Copy(buffer, reader); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err = ioutil.WriteFile(header.Name, buffer.Bytes(), 0755); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
