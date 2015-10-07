package task

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(url, path string) Task {
	return &download{url, path}
}

type download struct {
	url  string
	path string
}

func (t *download) Name() string {
	return "Downloading\n" + t.url + "\nto\n" + t.path
}

func (t *download) Execute() error {
	resp, err := http.Get(t.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(t.path)
	if err != nil {
		return err
	}
	defer file.Close()

	r := NewProgressReportingReader(resp.Body, resp.ContentLength,
		func(progress float64) {
			fmt.Printf("%v%%\b\b\b\b", int(100*progress+0.5))
		})
	_, err = io.Copy(file, r)
	return err
}

// NewProgressReportingReader returns a reader that reports to the given
// observer how much of its content was already read.
func NewProgressReportingReader(r io.Reader, sizeInBytes int64, observer func(float64)) *ProgressReportingReader {
	return &ProgressReportingReader{
		ProgressReader{Reader: r, size: sizeInBytes},
		observer,
	}
}

type ProgressReportingReader struct {
	ProgressReader
	observer func(float64)
}

func (r *ProgressReportingReader) Read(b []byte) (n int, err error) {
	n, err = r.ProgressReader.Read(b)
	r.observer(r.ProgressReader.Progress())
	return
}

// NewProgressReader returns a reader that knows how much of its content was
// already read.
func NewProgressReader(r io.Reader, sizeInBytes int64) *ProgressReader {
	return &ProgressReader{Reader: r, size: sizeInBytes}
}

type ProgressReader struct {
	io.Reader
	size int64
	read int64
}

func (r *ProgressReader) Progress() float64 {
	return float64(r.read) / float64(r.size)
}

func (r *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = r.Reader.Read(b)
	r.read += int64(n)
	return
}
