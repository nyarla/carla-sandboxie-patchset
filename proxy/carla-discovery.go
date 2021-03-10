package main

import (
	"debug/pe"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	sandboxiePrefixPath string
	sandboxieStartPath  string
	isSandboxSupported  bool
)

const (
	CARLA_SANDBOXIE_PREFIX = `CARLA_SANDBOXIE_PREFIX`
	CARLA_SANDBOXIE_START  = `CARLA_SANDBOXIE_START`
)

func normalizePath(src string) string {
	return strings.ReplaceAll(src, `/`, `\`)
}

func getSandboxDirs(src string) (box string, pluginPath string) {
	src = strings.TrimPrefix(normalizePath(src), sandboxiePrefixPath)
	dirs := strings.Split(src, `\`)

	for idx, dir := range dirs {
		if dir == "drive" {
			drive := dirs[idx+1]
			path := strings.Join(dirs[idx+2:], `\`)

			box = dirs[idx-1]
			pluginPath = fmt.Sprintf(`%s:\%s`, drive, path)

			return
		}
	}

	return
}

func getRealExecutable(src string) string {
	src = normalizePath(src)
	dir := filepath.Dir(src)
	fn := fmt.Sprintf(`_%s`, filepath.Base(src))

	return strings.Join([]string{dir, fn}, `\`)
}

func isSupportedArch(is64bit bool, fn string) bool {
	file, err := pe.Open(fn)
	if err != nil {
		return false
	}

	if is64bit {
		return file.FileHeader.Machine == pe.IMAGE_FILE_MACHINE_AMD64
	}

	return file.FileHeader.Machine == pe.IMAGE_FILE_MACHINE_I386
}

func getDiscoveryOutPath(src string, is64bit bool) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	src = re.ReplaceAllString(src, `_`)

	outBit := `64`
	if !is64bit {
		outBit = `32`
	}

	return fmt.Sprintf(`%s\carla-discovery_%s_%s`, os.Getenv(`TEMP`), outBit, src)
}

func init() {
	sandboxiePrefixPath = normalizePath(os.Getenv(CARLA_SANDBOXIE_PREFIX))
	sandboxieStartPath = normalizePath(os.Getenv(CARLA_SANDBOXIE_START))

	if sandboxiePrefixPath != "" && sandboxieStartPath != "" {
		isSandboxSupported = true
	}
}

func main() {
	is64bit := runtime.GOARCH == "amd64"
	pluginPath := normalizePath(os.Args[2])

	if !isSupportedArch(is64bit, pluginPath) {
		os.Exit(1)
	}

	hasSandboxPrefix := strings.HasPrefix(pluginPath, sandboxiePrefixPath)

	var cmd *exec.Cmd
	if isSandboxSupported && hasSandboxPrefix {
		box, fakePluginPath := getSandboxDirs(pluginPath)

		os.Args[0] = getRealExecutable(os.Args[0])
		os.Args[2] = fakePluginPath
		pluginPath = fakePluginPath

		cmdline := []string{sandboxieStartPath, fmt.Sprintf(`/box:%s`, box), `/silent`, `/nosbiectrl`}
		cmdline = append(cmdline, os.Args[0:]...)

		cmd = exec.Command(cmdline[0], cmdline[1:]...)
	} else {
		cmd = exec.Command(getRealExecutable(os.Args[0]), os.Args[1:]...)
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	outPath := getDiscoveryOutPath(pluginPath, is64bit)

	cmd.Start()
	cmd.Wait()

	_, err := os.Stat(outPath)
	for os.IsNotExist(err) {
		_, err = os.Stat(outPath)
		time.Sleep(1 * time.Second)
	}

	out, err := ioutil.ReadFile(outPath)
	if err != nil {
		os.Exit(1)
	}

	os.Stdout.Write(out)
	os.Remove(outPath)
	os.Exit(0)
}
