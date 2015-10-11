package main

import (
	"fmt"
	"github.com/gonutz/installer/task"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	for {
		tasks := []task.Task{
			installGo,
			installMinGW,
			installGit,
			installMercurial,
			installLiteIDE,
			installSDL2,
			installEverything,
		}

		for i, task := range tasks {
			fmt.Printf("(%v) %v\n", i+1, task.Name())
		}

		var choice int
		fmt.Scanf("%d\n", &choice)
		if choice >= 1 && choice <= len(tasks) {
			task := tasks[choice-1]
			if err := task.Execute(); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}

var (
	installGo = task.FailOnFirstError("Install Go", []task.Task{
		task.Check(
			func() bool { return exec.Command("cmd", "/C", "go", "version").Run() != nil },
			"Go is already installed",
		),
		task.Conditional(
			is32BitSystem,
			task.Download(
				`https://storage.googleapis.com/golang/go1.5.1.windows-386.msi`,
				userPath("Downloads", "go_installer.msi"),
			),
			task.Download(
				`https://storage.googleapis.com/golang/go1.5.1.windows-amd64.msi`,
				userPath("Downloads", "go_installer.msi"),
			),
		),
		task.RunProgram(userPath("Downloads", "go_installer.msi")),
		task.CreateFolder(userPath("Documents", "gocode")),
		task.SetEnv("GOPATH", userPath("Documents", "gocode")),
		task.AddToPathEnv(userPath("Documents", "gocode", "bin")),
	})

	installMinGW = task.FailOnFirstError("Install MinGW", []task.Task{
		task.Check(
			func() bool { return exec.Command("gcc", "-v").Run() != nil },
			"GCC is already installed",
		),
		task.WriteFile(
			userPath("Downloads", "mingw32setup.exe"),
			mingwGetSetupExe,
		),
		task.Inform(`Please install into the default direcotry (C:\MinGW).
On the second page, uncheck support for the graphical user interface.`),
		task.RunProgram(userPath("Downloads", "mingw32setup.exe")),
		task.AddToPathEnv(`C:\MinGW\msys\1.0\bin`),
		task.Conditional(
			is32BitSystem,
			task.FailOnFirstError("Installing GCC (32 Bit)", []task.Task{
				task.AddToPathEnv(`C:\MinGW\bin`),
				task.RunProgram(`C:\MinGW\bin\mingw-get.exe`, "install", "gcc"),
			}),
			task.FailOnFirstError("Installing MinGW (64 Bit)", []task.Task{
				task.WriteFile(
					userPath("Downloads", "mingw64setup.exe"),
					mingwW64InstallExe,
				),
				task.Inform("On the second page choose Architecture=x86_64"),
				task.RunProgram(userPath("Downloads", "mingw64setup.exe")),
				task.AddToPathEnv(mingw64binFolder()),
			}),
		),
	})

	installGit = task.FailOnFirstError("Install Git", []task.Task{
		task.Check(
			func() bool { return exec.Command("git", "version").Run() != nil },
			"Git is already installed",
		),
		task.Conditional(
			is32BitSystem,
			task.Download(
				`https://github.com/git-for-windows/git/releases/download/v2.6.1.windows.1/Git-2.6.1-32-bit.exe`,
				userPath("Downloads", "git_installer.exe"),
			),
			task.Download(
				`https://github.com/git-for-windows/git/releases/download/v2.6.1.windows.1/Git-2.6.1-64-bit.exe`,
				userPath("Downloads", "git_installer.exe"),
			),
		),
		task.Inform(`Choose
    "Use Git from the Windows Command Prompt"
and leave all other options on default.`),
		task.RunProgram(userPath("Downloads", "git_installer.exe")),
	})

	installMercurial = task.Inform("TODO install Mercurial")

	installLiteIDE = task.FailOnFirstError("Install LiteIDE", []task.Task{
		task.Download(
			`downloads.sourceforge.net/project/liteide/X23.2/liteidex23.2.windows.zip?r=http%3A%2F%2Fsourceforge.net%2Fprojects%2Fliteide%2Ffiles%2FX23.2%2F&ts=1444238477&use_mirror=skylink`,
			userPath("Downloads", "liteide23.2.zip"),
		),
		task.Unzip(
			userPath("Downloads", "liteide23.2.zip"),
			userPath("Downloads"),
		),
		task.CreateShortcut(
			userPath("Downloads", "liteidex23.2.windows", "liteide", "bin", "liteide.exe"),
			userPath("Desktop", "LiteIDE.lnk"),
		),
	})

	installSDL2 = task.Inform("TODO install SDL2")

	installEverything = task.ContinueAfterError("Install Everything", []task.Task{
		installGo,
		installMinGW,
		installGit,
		installMercurial,
		installLiteIDE,
		installSDL2,
	})
)

func userPath(path ...string) string {
	all := []string{os.Getenv("userprofile")}
	for _, p := range path {
		all = append(all, p)
	}
	return filepath.Join(all...)
}

func is32BitSystem() bool {
	return runtime.GOARCH == "386" || runtime.GOARCH == "arm"
}

func mingw64binFolder() (bin string) {
	root := `C:\Program Files\mingw-w64`
	filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
		if path == root {
			return nil
		}
		bin = path
		return filepath.SkipDir
	})
	bin = filepath.Join(bin, "mingw64", "bin")
	return
}

/*import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// TODO run mingw-get to install GCC
// simply run
// mingw-get install gcc
// after the MinGW32 installer has finished

const (
	go64URL            = `https://storage.googleapis.com/golang/go1.5.1.windows-amd64.msi`
	go32URL            = `https://storage.googleapis.com/golang/go1.5.1.windows-386.msi`
	mingw32URL         = `http://downloads.sourceforge.net/project/mingw/Installer/mingw-get-setup.exe?r=http%3A%2F%2Fsourceforge.net%2Fprojects%2Fmingw%2Ffiles%2F&ts=1444041541&use_mirror=heanet`
	mingw64URL         = `http://downloads.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe?r=http%3A%2F%2Fsourceforge.net%2Fprojects%2Fmingw-w64%2F&ts=1444041784&use_mirror=netcologne`
	msysBin            = `C:\MinGW\msys\1.0\bin`
	mingw32bin         = `C:\MinGW\bin`
	mingw64bin         = `C:\Program Files\mingw-w64\x86_64-4.9.0-win32-seh-rt_v3-rev2\mingw64\bin`
	git32URL           = `https://github.com/git-for-windows/git/releases/download/v2.6.1.windows.1/Git-2.6.1-32-bit.exe`
	git64URL           = `https://github.com/git-for-windows/git/releases/download/v2.6.1.windows.1/Git-2.6.1-64-bit.exe`
	sdl2URL            = `https://www.libsdl.org/release/SDL2-devel-2.0.3-mingw.tar.gz`
	mingw32includePath = `C:\MinGW\include`
	mingw32libPath     = `C:\MinGW\lib`
	mingw64includePath = `C:\Program Files\mingw-w64\x86_64-4.9.0-win32-seh-rt_v3-rev2\mingw64\x86_64-w64-mingw32\include`
	mingw64libPath     = `C:\Program Files\mingw-w64\x86_64-4.9.0-win32-seh-rt_v3-rev2\mingw64\x86_64-w64-mingw32\lib`
	sdl2imageURL       = `https://www.libsdl.org/projects/SDL_image/release/SDL2_image-devel-2.0.0-mingw.tar.gz`
	sdl2ttfURL         = `https://www.libsdl.org/projects/SDL_ttf/release/SDL2_ttf-devel-2.0.12-mingw.tar.gz`
	sdl2mixerURL       = `https://www.libsdl.org/projects/SDL_mixer/release/SDL2_mixer-devel-2.0.0-mingw.tar.gz`
)

func main() {
	for {
		fmt.Println("(1) install Go  (2) install MinGW  (3) install Git  (4) install SDL2")
		fmt.Println("(5) * install everything *")
		var choice int
		fmt.Scanf("%d\n", &choice)

		installFuncs := []func() error{
			nil,
			installGo,
			installMinGW,
			installGit,
			installSDL2,
			installEverything,
		}

		if choice >= 1 && choice < len(installFuncs) {
			install := installFuncs[choice]
			if err := install(); err != nil {
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

func runProgram(path string, args ...string) error {
	params := []string{"/C", path}
	for _, arg := range args {
		params = append(params, arg)
	}
	return exec.Command("cmd", params...).Run()
}

func createFolder(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		// create the folder if there is no such file
		return os.Mkdir(path, 0666)
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
//
// download git
// install git with PATH modification
//
// download SDL2_*
// unzip SDL2_*
// copy SDL2_* headers to mingw/include/SDL2
// copy SDL2_* lib to mingw/lib
// copy SDL2_* dlls to C:\Windows\System
// run go get github.com/gonutz/prototype/draw
// maybe run a sample to be sure
func installEverything() error {
	installFuncs := []func() error{
		installGo,
		installMinGW,
		installGit,
		installSDL2,
	}

	for _, install := range installFuncs {
		if err := install(); err != nil {
			if isAlreadyInstalledError(err) {
				inform("NOTE: skipping this component", err.Error())
			} else {
				return err
			}
		}
	}

	return nil
}

func installGo() error {
	defer stopProgress()

	goIsAlreadyInstalled := exec.Command("go", "version").Run() == nil
	if goIsAlreadyInstalled {
		return ErrAlreadyInstalled("Go")
	}

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
	startProgress("Creating GOPATH in\n" + gopath)
	if err := createFolder(gopath); err != nil {
		return err
	}
	if err := setEnvVariable("GOPATH", gopath); err != nil {
		return err
	}
	return addToPath(filepath.Join(gopath, "bin"))
}

func addToPath(add string) error {
	path := os.Getenv("PATH")
	alreadyInPath := strings.Contains(path, add+";") ||
		strings.Contains(path, `"`+add+`";`) ||
		strings.HasSuffix(path, add) ||
		strings.HasSuffix(path, `"`+add+`"`)
	if !alreadyInPath {
		newPath := path + ";" + add
		return setEnvVariable("PATH", newPath)
	}
	return nil
}

func installMinGW() error {
	defer stopProgress()

	mingwBin := mingw64bin
	if is32BitSystem() {
		mingwBin = mingw32bin
	}
	if folderExists(msysBin) && folderExists(mingwBin) {
		return ErrAlreadyInstalled("MinGW and msys")
	}

	mingw32 := userPath("Downloads", "mingw32setup.exe")
	startProgress("Copying MinGW Setup (32 Bit) to\n" + mingw32)
	if err := ioutil.WriteFile(mingw32, mingwGetSetupExe, 0666); err != nil {
		return err
	}

	inform("Installing MinGW (32 Bit)", `Please install into the default direcotry (C:\MinGW).
On the second page, uncheck support for the graphical user interface.`)
	startProgress("Installing MinGW (32 Bit)")
	if err := runProgram(mingw32); err != nil {
		return err
	}
	stopProgress()

	startProgress("Installing GCC C-Compiler from MinGW")
	// TODO test if this actually works
	if err := runProgram("mingw-get", "install", "gcc"); err != nil {
		return err
	}
	stopProgress()

	startProgress("Modifying PATH to include msys")
	if err := addToPath(msysBin); err != nil {
		return err
	}
	stopProgress()

	if is32BitSystem() {
		startProgress("Modifying PATH to include MinGW")
		if err := addToPath(mingw32bin); err != nil {
			return err
		}
		stopProgress()
	} else {
		mingw64 := userPath("Downloads", "mingw64setup.exe")
		startProgress("Copying MinGW Setup (64 Bit) to\n" + mingw64)
		if err := ioutil.WriteFile(mingw64, mingwW64InstallExe, 0666); err != nil {
			return err
		}
		stopProgress()

		startProgress("Installing MinGW (64 Bit)")
		if err := runProgram(mingw64); err != nil {
			return err
		}
		stopProgress()

		startProgress("Modifying PATH to include MinGW")
		if err := addToPath(mingw64bin); err != nil {
			return err
		}
		stopProgress()
	}

	return nil
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func installGit() error {
	defer stopProgress()

	gitIsAlreadyInstalled := exec.Command("git", "version").Run() == nil
	if gitIsAlreadyInstalled {
		return ErrAlreadyInstalled("Git")
	}

	gitURL := git64URL
	if is32BitSystem() {
		gitURL = git32URL
	}
	gitInstaller := userPath("Downloads", "git_installer.exe")
	stop := startProgress("Downloading Git from\n" + gitURL + "\nto\n" + gitInstaller)
	if err := download(gitURL, gitInstaller, observeProgress, stop); err != nil {
		return err
	}
	stopProgress()

	inform("Installing Git", `In the setup please choose the option
"Use Git from the Windows Command Prompt"
and leave all other settings on default.`)

	startProgress("Installing Git")
	return runProgram(gitInstaller)
}

func installSDL2() error {
	defer stopProgress()

	sdl2 := userPath("Downloads", "sdl2.tar.gz")
	stop := startProgress("Downloading SDL2 library from\n" + sdl2URL + "\nto\n" + sdl2)
	if err := download(sdl2URL, sdl2, observeProgress, stop); err != nil {
		return err
	}

	file, err := os.Open(sdl2)
	if err != nil {
		return err
	}
	defer file.Close()

	unzip, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	includePath := mingw64includePath
	if is32BitSystem() {
		includePath = mingw32includePath
	}
	includePath = filepath.Join(includePath, "SDL2")
	if err := createFolder(includePath); err != nil {
		return err
	}

	libPath := mingw64libPath
	if is32BitSystem() {
		libPath = mingw32libPath
	}

	systemPath := `C:\Windows\System32`

	tarReader := tar.NewReader(unzip)
	copyFile := func(folder, filename string) error {
		file, err := os.Create(filepath.Join(folder, filename))
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tarReader)
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if !header.FileInfo().IsDir() &&
			path.Dir(header.Name) == `SDL2-2.0.3/include` {
			if err := copyFile(includePath, path.Base(header.Name)); err != nil {
				return err
			}
		} else if !header.FileInfo().IsDir() &&
			(path.Dir(header.Name) == `SDL2-2.0.3/lib/x86` && is32BitSystem()) ||
			(path.Dir(header.Name) == `SDL2-2.0.3/lib/x64` && !is32BitSystem()) {
			if err := copyFile(libPath, path.Base(header.Name)); err != nil {
				return err
			}
		} else if path.Base(header.Name) == "SDL2.dll" {
			if is32BitSystem() &&
				path.Dir(header.Name) == `SDL2-2.0.3/i686-w64-mingw32/bin` {
				if err := copyFile(systemPath, path.Base(header.Name)); err != nil {
					return err
				}
			} else if !is32BitSystem() &&
				path.Dir(header.Name) == `SDL2-2.0.3/x86_64-w64-mingw32/bin` {
				if err := copyFile(systemPath, path.Base(header.Name)); err != nil {
					return err
				}
			}
		}
	}

	libs := []sdl2lib{
		{sdl2imageURL, "image"},
		{sdl2ttfURL, "ttf"},
		{sdl2mixerURL, "mixer"},
	}
	for _, lib := range libs {
		if err := installSDL2lib(lib); err != nil {
			return err
		}
	}

	return nil
}

type sdl2lib struct {
	url  string
	name string
}

func installSDL2lib(lib sdl2lib) error {
	defer stopProgress()

	archivePath := userPath("Downloads", "sdl2"+lib.name+".tar.gz")
	stop := startProgress("Downloading SDL2_" + lib.name + " library from\n" +
		lib.url + "\nto\n" + archivePath)
	if err := download(lib.url, archivePath, observeProgress, stop); err != nil {
		return err
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	unzip, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	includePath := mingw64includePath
	if is32BitSystem() {
		includePath = mingw32includePath
	}
	includePath = filepath.Join(includePath, "SDL2")

	libPath := mingw64libPath
	if is32BitSystem() {
		libPath = mingw32libPath
	}

	systemPath := `C:\Windows\System32`

	tarReader := tar.NewReader(unzip)

	var topLevelFolderName string
	var expectedPathStart string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if topLevelFolderName == "" {
			topLevelFolderName = header.Name
			expectedPathStart = topLevelFolderName + "x86_64-w64-mingw32"
			if is32BitSystem() {
				expectedPathStart = topLevelFolderName + "i686-w64-mingw32"
			}
		}

		if strings.HasPrefix(path.Dir(header.Name), expectedPathStart) {
			if path.Base(header.Name) == "SDL_"+lib.name+".h" {
				if err := writeToFile(tarReader, includePath,
					path.Base(header.Name)); err != nil {
					return err
				}
			}
			if path.Base(header.Name) == "libSDL2_"+lib.name+".a" ||
				path.Base(header.Name) == "libSDL2_"+lib.name+".dll.a" ||
				path.Base(header.Name) == "libSDL2_"+lib.name+".la" {
				if err := writeToFile(tarReader, libPath,
					path.Base(header.Name)); err != nil {
					return err
				}
			}
			if path.Dir(header.Name) == expectedPathStart+"/bin" &&
				strings.HasSuffix(header.Name, ".dll") {
				if err := writeToFileIfNotThere(tarReader, systemPath,
					path.Base(header.Name)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeToFile(r io.Reader, folder, filename string) error {
	file, err := os.Create(filepath.Join(folder, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func writeToFileIfNotThere(r io.Reader, folder, filename string) error {
	path := filepath.Join(folder, filename)
	if _, err := os.Stat(path); err == nil {
		// no error, so the file info was retrieved, so the file is there
		return nil
	}
	return writeToFile(r, folder, filename)
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
	fmt.Println("Press ENTER to continue...")
	fmt.Scanln()
}

type ErrAlreadyInstalled string

func (e ErrAlreadyInstalled) Error() string {
	return string(e) + " is already installed"
}

func isAlreadyInstalledError(err error) bool {
	_, ok := err.(ErrAlreadyInstalled)
	return ok
}
*/
// 650 lines in original version
