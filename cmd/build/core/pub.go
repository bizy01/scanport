package core

import (
	"fmt"
	"os"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bizy01/scanport/cliutils"
	"github.com/dustin/go-humanize"
)

func (c *Compiler) tarFiles(goos, goarch string) string {
	os.RemoveAll(c.PubDir)
	_ = os.MkdirAll(c.PubDir, os.ModePerm)

	gz := filepath.Join(c.PubDir, fmt.Sprintf("%s-%s-%s-%s.tar.gz",
		c.AppName, c.Version, goos, goarch))

	args := []string{
		`czf`,
		gz,
		`-C`,
		// the whole buildDir/datakit-<goos>-<goarch> dir
		filepath.Join(c.BuildDir, fmt.Sprintf("%s-%s-%s", c.AppName, goos, goarch)), `.`,
	}

	cmd := exec.Command("tar", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	return gz
}

func (c *Compiler) PubOSS() {
	start := time.Now()

	if _, err := os.Stat(c.BuildDir); err != nil {
		c.Compile()
		log.Println(err)
	}

	var ak, sk, bucket, ossHost string

	// 在你本地设置好这些 oss-key 环境变量
	switch c.Release {
	case `release`:
		tag := strings.ToUpper(c.Release)
		ak = os.Getenv(tag + "_OSS_ACCESS_KEY")
		sk = os.Getenv(tag + "_OSS_SECRET_KEY")
		bucket = os.Getenv(tag + "_OSS_BUCKET")
		ossHost = os.Getenv(tag + "_OSS_HOST")
	default:
		log.Printf("unknown release type: %s", c.Release)
	}

	if ak == "" || sk == "" {
		log.Printf("oss access key or secret key missing, tag=%s", strings.ToUpper(c.Release))
	}

	oc := &cliutils.OssCli{
		Host:       ossHost,
		PartSize:   512 * 1024 * 1024,
		AccessKey:  ak,
		SecretKey:  sk,
		BucketName: bucket,
		WorkDir:    c.OSSPath,
	}

	if err := oc.Init(); err != nil {
		log.Println(err)
	}
	// upload all build archs
	archs := []string{runtime.GOOS + "/" + runtime.GOARCH}

	ossfiles := map[string]string{}

	// tar files and collect OSS upload/backup info
	for _, arch := range archs {
		parts := strings.Split(arch, "/")
		if len(parts) != 2 {
			log.Printf("invalid arch %q", parts)
		}
		goos, goarch := parts[0], parts[1]

		gzName := c.tarFiles(goos, goarch)

		ossfiles[gzName] = path.Join(c.OSSPath, gzName)
	}

	// test if all file ok before uploading
	for k, _ := range ossfiles {
		if _, err := os.Stat(k); err != nil {
			log.Println(err)
		}
	}

	for k, v := range ossfiles {
		fi, _ := os.Stat(k)
		log.Printf("upload %s(%s)...", k, humanize.Bytes(uint64(fi.Size())))

		if err := oc.Upload(k, v); err != nil {
			log.Println(err)
		}
	}

	log.Printf("Done!(elapsed: %v)", time.Since(start))
}