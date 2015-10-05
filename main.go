package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	go64URL    = `https://storage.googleapis.com/golang/go1.5.1.windows-amd64.msi`
	go32URL    = `https://storage.googleapis.com/golang/go1.5.1.windows-386.msi`
	mingw32URL = `http://downloads.sourceforge.net/project/mingw/Installer/mingw-get-setup.exe?r=http%3A%2F%2Fsourceforge.net%2Fprojects%2Fmingw%2Ffiles%2F&ts=1444041541&use_mirror=heanet`
	mingw64URL = `http://downloads.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe?r=http%3A%2F%2Fsourceforge.net%2Fprojects%2Fmingw-w64%2F&ts=1444041784&use_mirror=netcologne`
	msysBin    = `C:\MinGW\msys\1.0\bin`
	mingw32bin = `C:\MinGW\bin`
	mingw64bin = `C:\Program Files\mingw-w64\x86_64-4.9.0-win32-seh-rt_v3-rev2\mingw64\bin`
)

func main() {
	for {
		fmt.Println("(1) install Go  (2) install MinGW")
		var choice int
		fmt.Scanf("%d\n", &choice)

		if choice == 1 {
			if err := installGo(); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}

func download(
	url, destFilePath string,
	progressObserver ProgressObserver,
	stopSignal <-chan struct{}) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	r := NewStoppableReader(NewProgressReportingReader(resp.Body, resp.ContentLength,
		func(progress float64) {
			progressObserver(false, progress)
		}))
	done := make(chan bool, 1)
	cancelled := false
	go func() {
		select {
		case <-done:
			return
		case <-stopSignal:
			cancelled = true
			r.stop()
		}
	}()
	_, err = io.Copy(file, r)
	done <- true
	progressObserver(true, 1)
	if cancelled {
		return errors.New("Operation was cancelled by the user.")
	}
	return err
}

type ProgressObserver func(done bool, progress float64)

func NewStoppableReader(r io.Reader) *readStopper {
	return &readStopper{r, false}
}

type readStopper struct {
	io.Reader
	stopped bool
}

func (r *readStopper) stop() {
	r.stopped = true
}

func (r *readStopper) Read(b []byte) (n int, err error) {
	if r.stopped {
		return 0, io.EOF
	}
	return r.Reader.Read(b)
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

func runProgram(path string) error {
	return exec.Command("cmd", "/C", path).Run()
}

func createFolder(path string) error {
	info, err := os.Stat(path)

	if err.(*os.PathError).Err == os.ErrNotExist {
		// create the folder if there is no such file
		return os.Mkdir(path, 0666)
	} else if err != nil {
		// report other errors because the file may already exists in this case
		return err
	}

	// at this point err is nil so we can check the file info
	if info.IsDir() {
		// if the folder already exists, we are done
		return nil
	}
	return errors.New("The path already exists but is not a directory.")
}

// setEnvVariable sets a system wide, permanent environment variable for the
// current user
func setEnvVariable(name, value string) error {
	return exec.Command("setx", name, value).Run()
}

func is32BitSystem() bool {
	return runtime.GOARCH == "386" || runtime.GOARCH == "arm"
}

func userPath(path ...string) string {
	all := []string{os.Getenv("userprofile")}
	for _, p := range path {
		all = append(all, p)
	}
	return filepath.Join(all...)
}

// installEverything does the following:
// download Go, 32 or 64 bit
// install Go
// create gocode folder
// set GOPATH
// modify PATH to include GOPATH/bin
// download MinGW, 32 or 64 bit
// install MinGW
// modify PATH to include msys/bin and mingw/bin
// download SDL2_*
// unzip SDL2_*
// copy SDL2_* headers to mingw/include/SDL2
// copy SDL2_* lib to mingw/lib
// copy SDL2_* dlls to C:\Windows\System
// download git
// install git
// run go get github.com/gonutz/prototype/draw
// maybe run a sample to be sure
func installEverything() error {
	defer stopProgress()

	if err := installGo(); err != nil {
		return err
	}

	if err := installMinGW(); err != nil {
		return err
	}

	return nil
}

func installGo() error {
	defer stopProgress()

	goURL := go64URL
	if is32BitSystem() {
		goURL = go32URL
	}
	goInstaller := userPath("Downloads", "go1.5.1_installer.msi")
	stop := startProgress("Downloading Go from\n" + goURL + "\nto\n" + goInstaller)
	if err := download(goURL, goInstaller, observeProgress, stop); err != nil {
		return err
	}
	stopProgress()

	startProgress("Installing Go")
	if err := runProgram(goInstaller); err != nil {
		return err
	}
	stopProgress()

	gopath := userPath("Documents", "gocode")
	startProgress("Createing GOPATH in\n" + gopath)
	if err := createFolder(gopath); err != nil {
		return err
	}
	if err := setEnvVariable("GOPATH", gopath); err != nil {
		return err
	}
	path := os.Getenv("PATH")
	binPath := filepath.Join(gopath, "bin")
	if !strings.Contains(path, binPath) {
		newPath := path + ";" + binPath
		return setEnvVariable("PATH", newPath)
	}

	return nil
}

func installMinGW() error {
	defer stopProgress()

	mingw32 := userPath("Downloads", "mingw32setup.exe")
	startProgress("Copying MinGW Setup (32 Bit) to\n" + mingw32)
	if err := ioutil.WriteFile(mingw32, mingwGetSetupExe[:], 0666); err != nil {
		return err
	}

	inform("Installing MinGW (32 Bit)", `Please install into the default direcotry (C:\MinGW).
On the second page, uncheck support for the graphical user interface.`)
	startProgress("Installing MinGW (32 Bit)")
	if err := runProgram(mingw32); err != nil {
		return err
	}
	stopProgress()

	if !is32BitSystem() {
		mingw64 := userPath("Downloads", "mingw64setup.exe")
		startProgress("Copying MinGW Setup (64 Bit) to\n" + mingw64)
		if err := ioutil.WriteFile(mingw64, mingwW64InstallExe[:], 0666); err != nil {
			return err
		}
		stopProgress()

		startProgress("Installing MinGW (64 Bit)")
		if err := runProgram(mingw64); err != nil {
			return err
		}
		stopProgress()
	}

	return nil
}

func startProgress(desc string) (stopSignal <-chan struct{}) {
	fmt.Println(desc)
	return nil
}

func observeProgress(done bool, p float64) {
	if done {
		stopProgress()
	} else {
		progressTo(p)
	}
}

func progressTo(progress float64) {
	fmt.Printf("%v%%\b\b\b\b\b\b", int(100*progress+0.5))
}

func stopProgress() {
	fmt.Printf("\b\b\b\b\b\b\b\b\n")
}

func inform(about, msg string) {
	fmt.Println("***", about, "***")
	fmt.Println(msg)
}
