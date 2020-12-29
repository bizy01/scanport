package scan

import (
	"time"
	"fmt"
	"net"
	"strings"
	"github.com/bizy01/scanport/config"
	"github.com/bizy01/scanport/cliutils"
	"github.com/bizy01/scanport/pool"
	"os"
	"log"
	"sync"
)

var wg = sync.WaitGroup{}

type Scanport struct{
    config.Config
	IPs []string
	Ports []int
	pool *pool.Pool
	ResChan chan string
}

func NewScan(cfg config.Config) (*Scanport, error) {
	scan := &Scanport{
		ResChan: make(chan string, cfg.Process),
	}

	scan.Config = cfg

	// new pool
	var err error

	scan.pool, err = pool.NewPool(scan.Process)
	if err != nil {
		return scan, err
	}

  return scan, nil
}

var res []string

// 获取ip列表
func (s *Scanport) getAllIP() {
	res := []string{}

	targets := []string{"127.0.0.1"}


	if len(s.Target) != 0 {
		//处理 ","号 如 80,81,88 或 80,88-100
		targets = strings.Split(strings.Trim(s.Target, ","), ",")
	}

	for _, target := range targets {
		target = strings.TrimSpace(target)
		if cliutils.IsIP(target) {
			res = append(res, target)
		} else if cliutils.IsDNS(target) {
			// 添加ip
			ip, _ := cliutils.ParseDNS(target)
			res = append(res, ip)
		} else if cliutils.IsCIDR(target) {
			// crdr解析
			cidrIp, _ := cliutils.ParseCIDR(target)
			res = append(res, cidrIp...)
		}
	}

	s.IPs = res
}

// 获取端口
func (s *Scanport) getAllPort()  {
	var ports []int

	portArr := []string{"80"}

	if len(s.Port) != 0 {
		//处理 ","号 如 80,81,88 或 80,88-100
		portArr = strings.Split(strings.Trim(s.Port, ","), ",")
	}

	for _, v := range portArr {
		portArr2 := strings.Split(strings.Trim(v, "-"), "-")
		startPort, err := cliutils.FilterPort(portArr2[0])
		if err != nil {
			continue
		}
		//第一个端口先添加
		ports = append(ports, startPort)
		if len(portArr2) > 1 {
			//添加第一个后面的所有端口
			endPort, _ := cliutils.FilterPort(portArr2[1])
			if endPort > startPort {
				for i := 1; i <= endPort-startPort; i++ {
					ports = append(ports, startPort+i)
				}
			}
		}
	}
	//去重复
	ports = cliutils.ArrayUnique(ports)

	s.Ports = ports
}

func (s *Scanport) Run() {
	s.getAllIP()
	s.getAllPort()
	s.scan()
}

func (s *Scanport) scan() {
	proto := strings.Split(s.Protocol, ",")

	for _, protocol := range proto {
		for _, ip := range s.IPs {
			for _, port := range s.Ports {
				// 协称池
				wg.Add(1)
				task := &pool.Task{
					Handler: func(v ...interface{}) {
						defer wg.Done()
						if isOpen(v[0].(string), v[1].(string), v[2].(int)) {
							s.ResChan <- fmt.Sprintf("%v:%v:%v",v[0], v[1], v[2])
						}
					},
					Params: []interface{}{strings.Trim(protocol, ""), ip, port},
				}

				s.pool.Put(task)
			}
		}
	}


	wg.Wait()
	close(s.ResChan)
}

// open test
func isOpen(proto string, ip string, port int) bool {
	conn, err := net.DialTimeout(proto, fmt.Sprintf("%s:%d", ip, port), 100*time.Millisecond)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			log.Println("too many open files" + err.Error())
			os.Exit(1)
		}
		return false
	}

	conn.Close()

	return true
}

// 结果输出
func (s *Scanport) Output() {
	for item := range  s.ResChan {
		res := strings.Split(item, ":")
		fmt.Printf("协议类型: %v  扫描目标: %-15v  开放端口: %-5v\n", res[0], res[1], res[2])
	}
}

