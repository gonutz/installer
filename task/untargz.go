package task

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func UnTarGZ(archive, extractTo string) Task {
	return &untargz{archive, extractTo}
}

type untargz struct {
	archive, path string
}

func (t *untargz) Name() string {
	return "Extracting\n" + t.archive + "\nto\n" + t.path
}

func (t *untargz) Execute() error {
	file, err := os.Open(t.archive)
	if err != nil {
		return makeError("opening archive '"+t.archive+"'", err)
	}
	defer file.Close()

	zipReader, err := gzip.NewReader(file)
	if err != nil {
		return makeError("creating gzip decompressor for file '"+file.Name()+"'", err)
	}
	defer zipReader.Close() // ignore the possible error here

	tarReader := tar.NewReader(zipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return makeError("reading tar ball", err)
		}

		path := filepath.Join(t.path, filepath.FromSlash(header.Name))
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0666); err != nil {
				return makeError("creating folders to '"+path+"'", err)
			}
		} else {
			if err := createFile(path, tarReader); err != nil {
				return err
			}
		}
	}

	return nil
}

func createFile(path string, r io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return makeError("creating file '"+path+"'", err)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return makeError("copying data to '"+path+"'", err)
}
