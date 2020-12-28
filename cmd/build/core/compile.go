package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"path"
	"runtime"
	"strings"
	"time"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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

var (
	l = logrus.New()
)

func (c *Compiler) Compile() {
	start := time.Now()

	os.RemoveAll(c.BuildDir)
	_ = os.MkdirAll(c.BuildDir, os.ModePerm)

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
			l.Fatalf("invalid arch %q", parts)
		}

		goos, goarch := parts[0], parts[1]

		dir := fmt.Sprintf("%s/%s-%s-%s", c.BuildDir, c.AppName, goos, goarch)

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			l.Fatalf("failed to mkdir: %v", err)
		}

		dir, err = filepath.Abs(dir)
		if err != nil {
			l.Fatal(err)
		}

		compileArch(c.MainEntry, c.AppBin, goos, goarch, dir, c.Version)
	}

	l.Infof("Done!(elapsed %v)", time.Since(start))
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

	l.Debugf("building %s", fmt.Sprintf("%s-%s/%s", goos, goarch, bin))
	msg, err := runEnv(args, env)
	if err != nil {
		l.Fatalf("failed to run %v, envs: %v: %v, msg: %s", args, env, err, string(msg))
	}
}

func copyFile(source, dist string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
	   return err
	}

	err = ioutil.WriteFile(dist, input, 0644)
	if err != nil {
	   return err
	}

	return nil
}

// Dir copies a whole directory recursively
func copyDir(src string, dst string) error {
    var err error
    var fds []os.FileInfo
    var srcinfo os.FileInfo

    if srcinfo, err = os.Stat(src); err != nil {
        return err
    }

    if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
        return err
    }

    if fds, err = ioutil.ReadDir(src); err != nil {
        return err
    }
    for _, fd := range fds {
        srcfp := path.Join(src, fd.Name())
        dstfp := path.Join(dst, fd.Name())

        if fd.IsDir() {
            if err = copyDir(srcfp, dstfp); err != nil {
                l.Error(err)
            }
        } else {
            if err = copyFile(srcfp, dstfp); err != nil {
                l.Error(err)
            }
        }
    }
    return nil
}



