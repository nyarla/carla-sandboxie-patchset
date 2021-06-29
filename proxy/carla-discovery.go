package main

import (
	"debug/pe"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	unixLikePrefix       string
	unixLikeStart        string
	isSandboxieSupported bool
)

func getWindowsPath(src string) string {
	return strings.ReplaceAll(src, `/`, `\`)
}

func getUnixLikePath(src string) string {
	return strings.ReplaceAll(src, `\`, `/`)
}

func getSandboxieContainerPath(src string) (box string, fakePath string) {
	dirs := strings.Split(getUnixLikePath(src), `/`)

	for idx, dir := range dirs {
		if dir == `drive` {
			drive := dirs[idx+1]
			path := strings.Join(dirs[idx+2:], `/`)

			box = dirs[idx-1]
			fakePath = fmt.Sprintf(`%s:/%s`, drive, path)

			return
		}
	}

	return
}

func getRealPathFromSymlink(src string) (ok bool, realPath string) {
	path := getWindowsPath(src)

	if dest, err := os.Readlink(path); err != nil {
		ok = false
		realPath = src
		return
	} else {
		ok = true
		realPath = dest
	}

	return
}

func isSandboxedPlugin(src string) bool {
	return strings.HasPrefix(getUnixLikePath(src), unixLikePrefix)
}

func getExecPath(src string) string {
	dir := filepath.Dir(src)
	fn := filepath.Base(src)

	return getUnixLikePath(filepath.Join(dir, fmt.Sprintf(`_%s`, fn)))
}

func init() {
	unixLikePrefix = getUnixLikePath(os.Getenv(`CARLA_SANDBOXIE_PREFIX`))
	unixLikeStart = getUnixLikePath(os.Getenv(`CARLA_SANDBOXIE_START`))
	isSandboxieSupported = unixLikePrefix != `` && unixLikeStart != ``
}

func isSupportedArch(is64bit bool, src string) bool {
	path := getWindowsPath(src)
	file, err := pe.Open(path)

	if err != nil {
		return false
	}

	if is64bit {
		return file.FileHeader.Machine == pe.IMAGE_FILE_MACHINE_AMD64
	}

	return file.FileHeader.Machine == pe.IMAGE_FILE_MACHINE_I386
}

func getDiscoveryOutPath(src string, is64bit bool) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	path := re.ReplaceAllString(getWindowsPath(src), `_`)
	bit := `64`
	if !is64bit {
		bit = `32`
	}

	dest := fmt.Sprintf(`%s\carla-discovery_%s_%s`, os.Getenv(`TEMP`), bit, path)

	return getWindowsPath(dest)
}

func main() {
	is64bit := runtime.GOARCH == "amd64"

	if ok, realWindowsPath := getRealPathFromSymlink(os.Args[2]); ok {
		os.Args[2] = getUnixLikePath(realWindowsPath)
	}

	pluginPath := os.Args[2]
	if !isSupportedArch(is64bit, pluginPath) {
		os.Exit(0)
	}

	var (
		cmd, terminate *exec.Cmd
		isSandboxied   bool = isSandboxedPlugin(os.Args[2])
	)

	if isSandboxieSupported && isSandboxied {
		box, fakePath := getSandboxieContainerPath(os.Args[2])

		os.Args[0] = getExecPath(os.Args[0])
		os.Args[2] = getUnixLikePath(fakePath)
		pluginPath = getUnixLikePath(fakePath)

		cmdline := []string{unixLikeStart, fmt.Sprintf(`/box:%s`, box), `/wait`, `/silent`, `/nosbiectrl`}
		cmdline = append(cmdline, os.Args[0:]...)

		cmd = exec.Command(cmdline[0], cmdline[1:]...)
		terminate = exec.Command(unixLikeStart, fmt.Sprintf(`/box:%s`, box), `/terminate`, `/wait`)
	} else {
		cmd = exec.Command(getExecPath(os.Args[0]), os.Args[1:]...)
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	outPath := getDiscoveryOutPath(pluginPath, is64bit)

	if isSandboxieSupported && isSandboxied {
		cmd.Start()
	} else {
		cmd.Run()
	}

	count := 0
	_, err := os.Stat(outPath)
	for os.IsNotExist(err) {
		_, err = os.Stat(outPath)
		count++
		time.Sleep(1 * time.Second)
		if count > 60 {
			os.Exit(1)
		}
	}

	out, err := os.ReadFile(outPath)
	if err == nil {
		os.Stdout.Write(out)
	}

	os.Remove(outPath)
	os.Remove(outPath + `.tmp`)

	if isSandboxieSupported && isSandboxied {
		terminate.Run()
	}

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
