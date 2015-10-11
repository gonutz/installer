package task

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipFile, path string) Task {
	return &unzip{zipFile, path}
}

type unzip struct {
	zipFileName, path string
}

func (t *unzip) Name() string {
	return "Unzipping\n" + t.zipFileName + "\nto\n" + t.path
}

func (t *unzip) Execute() error {
	zipFile, err := os.Open(t.zipFileName)
	if err != nil {
		return makeError("opening file '"+t.zipFileName+"'", err)
	}
	defer zipFile.Close()

	info, err := zipFile.Stat()
	if err != nil {
		return makeError("reading file info of '"+zipFile.Name()+"'", err)
	}

	zipReader, err := zip.NewReader(zipFile, info.Size())
	if err != nil {
		return makeError("creating zip file reader for '"+zipFile.Name()+"'", err)
	}

	for _, f := range zipReader.File {
		destPath := filepath.Join(t.path, filepath.FromSlash(f.Name))
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0666); err != nil {
				return makeError("creating path to folder '"+destPath+"'", err)
			}
		} else {
			if err := copyFile(f, destPath); err != nil {
				return makeError(
					"copying zip file '"+f.Name+"' data to file '"+destPath+"'",
					err)
			}
		}
	}

	return nil
}

func copyFile(zipFile *zip.File, destPath string) error {
	src, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}
