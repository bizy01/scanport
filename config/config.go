package config

import (
	"github.com/BurntSushi/toml"
	"sync"
	"path/filepath"
)

const simpleCfg  = `
[scanport]
	## 【扫描的协议】
	# default:
	# tcp
	# usage:
	# tcp,udp
   	protocol = "tcp"
   	## 【扫描的目标主机，支持ip, 域名，cidr】
   	# default:
   	# 127.0.0.1
   	# usage:
   	# (1): 127.0.0.1
    # (2): 192.168.0.1, 192.168.0.2
    # (3): 192.168.0.1-20
    # (4): www.baidu.com
    # (5): 192.168.1.1/30
   	target = "127.0.0.1"
   	## 【端口值】
   	# default:
   	# 80
   	# usage:
   	# (1): 3000, 8080, 3306
   	# (2): 3000-10000
   	# (3): 8080,3000-10000
   	port = "8000-30000"
   	## 【扫描并发数】
   	# default
   	# 1000
   	process = 1000
   	## 【Dial timeout(unit Millisecond)】
   	# default
   	# 100 (Millisecond)
   	timeout = 100
`

type Config struct {
	Protocol string	      `toml:"protocol"`
	Target   string       `toml:"target"`
	Port     string	      `toml:"port"`
	Process  uint64       `toml:"process"`
	Timeout  int          `toml:"timeout"`
}

var (
	cfg Config
	once sync.Once
)

func GetConfig(path string) Config {
	once.Do(func() {
		filePath, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		if _ , err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func InitConfig() []byte {
   return []byte(simpleCfg)
}