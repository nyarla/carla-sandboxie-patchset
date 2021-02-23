package main

import (
	"debug/pe"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	sandboxiePrefix    string
	sandboxieStart     string
	isSandboxSupported bool
	inSandbox          bool
)

func normalizePath(src string) string {
	return strings.ReplaceAll(src, `\`, `/`)
}

func getSandboxieContainerInfo(src string) (box string, fakePath string) {
	dirs := strings.Split(strings.TrimPrefix(normalizePath(src), sandboxiePrefix), `/`)
	for idx, seg := range dirs {
		if seg == "drive" {
			drive := dirs[idx+1]
			path := strings.Join(dirs[idx+2:], `/`)

			box = dirs[idx-1]
			fakePath = fmt.Sprintf(`%s:/%s`, drive, path)

			return
		}
	}

	return
}

func executable(src string) string {
	return normalizePath(filepath.Join(filepath.Dir(src), fmt.Sprintf(`_%s`, filepath.Base(src))))
}

func checkArch(is64bit bool, fp string) bool {
	file, err := pe.Open(fp)
	if err != nil {
		return false
	}

	if is64bit {
		return file.FileHeader.Machine == pe.IMAGE_FILE_MACHINE_AMD64
	} else {
		return file.FileHeader.Machine == pe.IMAGE_FILE_MACHINE_I386
	}
}

func discoveryOutPath(src string) string {
	return strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(src, `_`))
}

func init() {
	sandboxiePrefix = normalizePath(os.Getenv(`CARLA_SANDBOXIE_PREFIX`))
	sandboxieStart = normalizePath(os.Getenv(`CARLA_SANDBOXIE_START`))

	if sandboxiePrefix != "" && sandboxieStart != "" {
		isSandboxSupported = true
	}
}

func main() {
	if !checkArch(strings.HasSuffix(os.Args[0], `64.exe`), os.Args[2]) {
		os.Exit(1)
	}

	var cmd *exec.Cmd
	inSandbox = strings.HasPrefix(normalizePath(os.Args[2]), normalizePath(sandboxiePrefix))

	if isSandboxSupported && inSandbox {
		box, fakePath := getSandboxieContainerInfo(os.Args[2])

		os.Args[0] = executable(os.Args[0])
		os.Args[2] = fakePath

		commandLine := []string{sandboxieStart, fmt.Sprintf(`/box:%s`, box), `/silent`, `/nosbiectrl`}
		commandLine = append(commandLine, os.Args[0:]...)
		cmd = exec.Command(commandLine[0], commandLine[1:]...)
	} else {
		cmd = exec.Command(executable(os.Args[0]), os.Args[1:]...)
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Errorf("%s", err)
	}

	cmd.Wait()
	fp := fmt.Sprintf("%s/carla-discovery_%s", os.Getenv(`TEMP`), discoveryOutPath(os.Args[2]))
	defer os.Remove(fp)

	for count := 0; count < 5; count++ {
		if _, err := os.Stat(fp); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if out, err := ioutil.ReadFile(fp); err == nil {
		os.Stdout.Write(out)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	os.Exit(0)
}
