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
				userPath("go_installer.msi"),
			),
			task.Download(
				`https://storage.googleapis.com/golang/go1.5.1.windows-amd64.msi`,
				userPath("go_installer.msi"),
			),
		),
		task.RunProgram(userPath("go_installer.msi")),
		task.CreateFolder(userPath("gocode")),
		task.SetEnv("GOPATH", userPath("gocode")),
		task.AddToPathEnv(userPath("gocode", "bin")),
	})

	installMinGW = task.FailOnFirstError("Install MinGW", []task.Task{
		task.Check(
			func() bool { return exec.Command("gcc", "-v").Run() != nil },
			"GCC is already installed",
		),
		task.WriteFile(
			userPath("mingw32setup.exe"),
			mingwGetSetupExe,
		),
		task.Inform(`Please install into the default direcotry (C:\MinGW).
On the second page, uncheck support for the graphical user interface.`),
		task.RunProgram(userPath("mingw32setup.exe")),
		task.AddToPathEnv(`C:\MinGW\msys\1.0\bin`),
		task.Conditional(
			is32BitSystem,
			task.FailOnFirstError("Installing GCC (32 Bit)", []task.Task{
				task.AddToPathEnv(`C:\MinGW\bin`),
				task.RunProgram(`C:\MinGW\bin\mingw-get.exe`, "install", "gcc"),
			}),
			task.FailOnFirstError("Installing MinGW (64 Bit)", []task.Task{
				task.WriteFile(
					userPath("mingw64setup.exe"),
					mingwW64InstallExe,
				),
				task.Inform("On the second page choose Architecture=x86_64"),
				task.RunProgram(userPath("mingw64setup.exe")),
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
				userPath("git_installer.exe"),
			),
			task.Download(
				`https://github.com/git-for-windows/git/releases/download/v2.6.1.windows.1/Git-2.6.1-64-bit.exe`,
				userPath("git_installer.exe"),
			),
		),
		task.Inform(`In the installer choose
    "Use Git from the Windows Command Prompt"
and leave all other options on default.`),
		task.RunProgram(userPath("git_installer.exe")),
	})

	installMercurial = task.Inform("TODO install Mercurial")

	installLiteIDE = task.FailOnFirstError("Install LiteIDE", []task.Task{
		task.Download(
			`http://downloads.sourceforge.net/project/liteide/X23.2/liteidex23.2.windows.zip?r=&ts=1444591804&use_mirror=netassist`,
			userPath("liteide23.2.zip"),
		),
		task.Unzip(
			userPath("liteide23.2.zip"),
			userPath(),
		),
		task.CreateShortcut(
			userPath("liteidex23.2.windows", "liteide", "bin", "liteide.exe"),
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
