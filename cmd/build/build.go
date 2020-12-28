package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/bizy01/scanport/cmd/build/core"
)

var (
	// build
	flagBuild         = flag.Bool(`build`, true, `build program to build dir`)

	// pub
	flagPub           = flag.Bool(`pub`,   false, `publish binaries to OSS`)
	flagBinary       = flag.String("binary", "echoServer", "binary name to build")
	flagName         = flag.String("name", *flagBinary, "same as -binary")
	flagBuildDir     = flag.String("build-dir", "dist", "output of build files")
	flagDownloadAddr = flag.String("download-addr", "echoServer", "oss path")
	flagArchs        = flag.String("archs", "", "archs")
	flagMain         = flag.String("main", "", "main path")
	flagPubDir       = flag.String("pub-dir", "pub", "")
	flagVersion      = flag.String("version", "", "version")
	flagEnv          = flag.String(`env`, ``, `build for dev/release`)

	l = logrus.New()
)

func main() {
	flag.Parse()

	c := &core.Compiler{
		AppName: *flagBinary,
		AppBin:  *flagName,
		BuildDir: *flagBuildDir,
		Release:  *flagEnv,
		PubDir:   *flagPubDir,
		Build: *flagBuild,
		Archs: *flagArchs,
		Pub: *flagPub,
		MainEntry: *flagMain,
		Version: *flagVersion,
		OSSPath: "echoServer",
	}

	// 编译
	if *flagBuild {
		c.Compile()
	}

	// 发布oss
	if *flagPub {
		c.PubOSS()
	}
}
