package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func main() {
	var cmd *exec.Cmd

	if ok, realWindowsPath := getRealPathFromSymlink(os.Args[2]); ok {
		os.Args[2] = getUnixLikePath(realWindowsPath)
	}

	if isSandboxieSupported && isSandboxedPlugin(os.Args[2]) {
		box, fakePath := getSandboxieContainerPath(os.Args[2])

		os.Args[0] = getExecPath(os.Args[0])
		os.Args[2] = getUnixLikePath(fakePath)

		cmdline := []string{unixLikeStart, fmt.Sprintf(`/box:%s`, box), `/wait`, `/silent`, `/nosbiectrl`}
		cmdline = append(cmdline, os.Args[0:]...)

		cmd = exec.Command(cmdline[0], cmdline[1:]...)
	} else {
		os.Args[0] = getUnixLikePath(os.Args[0])
		os.Args[2] = getUnixLikePath(os.Args[2])

		cmd = exec.Command(getExecPath(os.Args[0]), os.Args[1:]...)
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cmd.Wait()
}
