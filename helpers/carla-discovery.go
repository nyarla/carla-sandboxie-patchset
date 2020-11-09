package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func init() {
	sandboxiePrefix = normalizePath(os.Getenv(`CARLA_SANDBOXIE_PREFIX`))
	sandboxieStart = normalizePath(os.Getenv(`CARLA_SANDBOXIE_START`))

	if sandboxiePrefix != "" && sandboxieStart != "" {
		isSandboxSupported = true
	}
}

func main() {
	var cmd *exec.Cmd

	inSandbox = strings.HasPrefix(os.Args[2], sandboxiePrefix)

	if isSandboxSupported && inSandbox {
		box, fakePath := getSandboxieContainerInfo(os.Args[2])

		os.Args[0] = executable(os.Args[0])
		os.Args[2] = fakePath

		commandLine := []string{sandboxieStart, fmt.Sprintf(`/box:%s`, box), `/wait`}
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

	fp := fmt.Sprintf("%s/carla-discovery.out", os.Getenv(`TEMP`))

	fmt.Println(fp)

	if out, err := ioutil.ReadFile(fp); err == nil {
		os.Stdout.Write(out)
	} else {
		fmt.Fprint(os.Stderr, "%s\n", err)
	}

	os.Remove(fp)
	os.Exit(0)
}
