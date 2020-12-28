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

const (
	protocolUsage = `
<扫描的协议>
usage: tcp,udp`

	targetUsage   = `
<扫描的目标主机，支持ip, 域名，cidr>
default: 127.0.0.1
usage:
(1): 127.0.0.1
(2): 192.168.0.1, 192.168.0.2
(3): 192.168.0.1-20
(4): www.baidu.com
(5): 192.168.1.1/30`

	portUsage     = `
<端口值>
default: 80
usage:
(1): 3000, 8080, 3306
(2): 3000-10000
(3): 8080,3000-10000`

	processUsage  = `
<扫描并发数>
default: 1000`

	timeoutUsage  = `
<Dial timeout(unit Millisecond)>
default: 100 (Millisecond)`
)

var (
	flagCfgPath   = flag.String("c", "", "config path")
    flagCfgSimple = flag.Bool("init", false, "init simple config")
    flagVersion   = flag.Bool("version", false, "version")

    flagProtocol  = flag.String("protocol", "tcp", protocolUsage)
    flagTarget    = flag.String("target", "127.0.0.1", targetUsage)
    flagPort      = flag.String("port", "80", portUsage)
    flagProcess   = flag.Int("process", 100, processUsage)
    flagTimeout   = flag.Int("timeout", 100, timeoutUsage)

    cfg = new(config.Config)
)

func usage() {
    fmt.Printf("usage of %s\n", os.Args[0])
    fmt.Println(protocolUsage)
    fmt.Println(targetUsage)
    fmt.Println(portUsage)
    fmt.Println(processUsage)
    fmt.Println(timeoutUsage)
    os.Exit(0)
}


func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NFlag() == 0 {
        usage()
        os.Exit(-1)
    }

	if *flagTarget != "" {
		cfg.Target = *flagTarget
	}

	if *flagPort != "" {
		cfg.Port = *flagPort
	}

	if *flagCfgPath != "" {
		cfg = config.GetConfig(*flagCfgPath)
	}

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

	s :=scan.NewScan(cfg.Target, cfg.Port, cfg.Process)

	s.Run()

	s.Output()
}
