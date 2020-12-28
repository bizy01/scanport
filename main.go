package main

import (
	"github.com/bizy01/scanport/scan"
	"github.com/bizy01/scanport/git"
	"github.com/bizy01/scanport/config"
	"log"
	"os"
	"io/ioutil"
	"flag"
	"fmt"
)

var (
	flagCfgPath = flag.String("c", "", "config path")
    flagCfgSimple = flag.Bool("init", false, "init config demo")
    flagVersion = flag.Bool("version", false, "version")

    flagTarget = flag.String("target", "", "target ip or domain or cidr")
    flagPort  = flag.String("port", "", "target port range")

    cfg = new(config.Config)
)

func main() {
	flag.Parse()

	if *flagTarget != "" {
		cfg.Target = *flagTarget
	}

	if *flagPort != "" {
		cfg.Port = *flagPort
	}

	// if *flagCfgPath != "" {
	// 	cfg = config.GetConfig(*flagCfgPath)
	// }

	if *flagCfgSimple  {
		log.Println("init config")

		context := config.InitConfig()
		if err := ioutil.WriteFile("./demo.toml", context, 0644); err != nil {
	        log.Fatalf("WriteFile failure, err=[%v]\n", err)
	    }

	    os.Exit(0)
	}

	if *flagVersion {
	fmt.Printf(`
       Version: %s
        Commit: %s
        Branch: %s
 Build At(UTC): %s
Golang Version: %s
      Uploader: %s
`, git.Version, git.Commit, git.Branch, git.BuildAt, git.Golang, git.Uploader)
			os.Exit(0)
	}

	cfg.Process = 100

	s :=scan.NewScan(cfg.Target, cfg.Port, cfg.Process)

	s.Run()

	s.Output()
}
