package util

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DecompressAll(folderPath string, filePrefix string, fileExtension string) error {
	fi, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, f := range fi {
		if strings.HasPrefix(f.Name(), filePrefix) && strings.HasSuffix(f.Name(), fileExtension) {
			err := Decompress(folderPath, f.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func Decompress(folderPath, fileName string) error {
	in, out, err := openFiles(folderPath, fileName)
	if err != nil {
		return err
	}

	defer in.Close()
	defer out.Close()

	r, err := gzip.NewReader(in)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return r.Close()
}

func openFiles(folderPath, inputFileName string) (in *os.File, out *os.File, err error) {
	in, err = os.Open(filepath.Join(folderPath, inputFileName))
	if err != nil {
		return nil, nil, err
	}

	out, err = os.Create(filepath.Join(folderPath, inputFileName[:len(inputFileName)-3]))
	if err != nil {
		in.Close()
		return nil, nil, err
	}

	return in, out, nil
}
