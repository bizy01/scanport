package core

import (
	"fmt"
	"os"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	/* Use:
		go tool dist list
	to get current os/arch list */

	OSArches = []string{ // supported os/arch list
		`linux/386`,
		`linux/amd64`,
		`linux/arm`,
		`linux/arm64`,
		`darwin/amd64`,
		`windows/amd64`,
		`windows/386`,
	}
)

type Compiler struct {
	AppName string
	AppBin  string
	BuildDir string
	Version  string
	Release  string
	Archs    string
	PubDir   string
	OSSPath  string
	MainEntry string
	Build    bool
	Pub      bool
}

func runEnv(args, env []string) ([]byte, error) {
	cmd := exec.Command(args[0], args[1:]...)
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}

	return cmd.CombinedOutput()
}

func (c *Compiler) Compile() {
	start := time.Now()

	os.RemoveAll(c.BuildDir)
	// _ = os.MkdirAll(c.BuildDir, os.ModePerm)

	var archs []string

	switch c.Archs {
	case "all":
		archs = OSArches
	case "local":
		archs = []string{runtime.GOOS + "/" + runtime.GOARCH}
	default:
		archs = strings.Split(c.Archs, "|")
	}

	for idx, _ := range archs {
		parts := strings.Split(archs[idx], "/")
		if len(parts) != 2 {
			log.Printf("invalid arch %q\n", parts)
		}

		goos, goarch := parts[0], parts[1]

		dir := fmt.Sprintf("%s/%s-%s-%s", c.BuildDir, c.AppName, goos, goarch)

		dir, err := filepath.Abs(dir)
		if err != nil {
			log.Println(err)
		}

		compileArch(c.MainEntry, c.AppBin, goos, goarch, dir, c.Version)
	}

	log.Printf("Done!(elapsed %v)", time.Since(start))
}

func compileArch(mainEntry, bin, goos, goarch, dir, version string) {
	output := filepath.Join(dir, bin)
	args := []string{
		"go", "build",
		"-ldflags",
		fmt.Sprintf("-w -s -X main.Version=%s", version),
		"-o", output,
		mainEntry,
	}

	env := []string{
		"GOOS=" + goos,
		"GOARCH=" + goarch,
	}

	log.Printf("building %s", fmt.Sprintf("%s-%s/%s", goos, goarch, bin))
	msg, err := runEnv(args, env)
	if err != nil {
		log.Printf("failed to run %v, envs: %v: %v, msg: %s", args, env, err, string(msg))
	}
}





