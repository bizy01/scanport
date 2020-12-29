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
<生成默认配置文件>
  usage: -init

<以加载配置文件方式运行>
  usage: -c

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
    flagProcess   = flag.Uint64("process", 100, processUsage)
    flagTimeout   = flag.Int("timeout", 100, timeoutUsage)

    cfg = config.Config{}
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

func Init() {
	cfg.Protocol = *flagProtocol
	cfg.Target = *flagTarget
	cfg.Port = *flagPort
	cfg.Process = *flagProcess
	cfg.Timeout = *flagTimeout
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(*flagCfgPath) != 0 {
		cfg = config.GetConfig(*flagCfgPath)
	} else {
		Init()
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

	s, err :=scan.NewScan(cfg)
	if err != nil {
		log.Printf("new create groutinue pool error %v", err)
		os.Exit(-1)
	}

	s.Run()
	s.Output()
}
