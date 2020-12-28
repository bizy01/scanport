package config

import (
	"github.com/BurntSushi/toml"
	"sync"
	"path/filepath"
	"time"
)

const simpleCfg  = `
[scanport]
   # 扫描的目标ip
   target = "127.0.0.1"
   port = "8000-30000"
   count = 100
`

type Config struct {
	Target string        `toml:"target"`
	Port   string	     `toml:"port"`
	Process  uint64         `toml:"process"`
	Timeout time.Duration `toml:"timeout"`
}


var (
	cfg *Config
	once sync.Once
)

func GetConfig(path string) *Config {
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