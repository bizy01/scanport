package scan

import (
	"time"
	"fmt"
	"net"
	"strings"
	// "github.com/bizy01/scanport/config"
	"github.com/bizy01/scanport/util"
	"github.com/bizy01/scanport/pool"
	"os"
	"log"
	"sync"
)

var wg = sync.WaitGroup{}

type Scanport struct{
	Target string
	Port   string
	IPs []string
	Ports []int
	Result   []string
	Timeout time.Duration
	MaxProcess int
	debug bool
	pool *pool.Pool
	ResChan chan string
}

func NewScan(target, port string) *Scanport {
	scan := &Scanport{
		ResChan: make(chan string, 100),
	}

	scan.Target = target
	scan.Port = port

  return scan
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
		if util.IsIP(target) {
			res = append(res, target)
		} else if util.IsDNS(target) {
			// 添加ip
			ip, _ := util.ParseDNS(target)
			res = append(res, ip)
		} else if util.IsCIDR(target) {
			// crdr解析
			cidrIp, _ := util.ParseCIDR(target)
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
		startPort, err := util.FilterPort(portArr2[0])
		if err != nil {
			continue
		}
		//第一个端口先添加
		ports = append(ports, startPort)
		if len(portArr2) > 1 {
			//添加第一个后面的所有端口
			endPort, _ := util.FilterPort(portArr2[1])
			if endPort > startPort {
				for i := 1; i <= endPort-startPort; i++ {
					ports = append(ports, startPort+i)
				}
			}
		}
	}
	//去重复
	ports = util.ArrayUnique(ports)

	s.Ports = ports
}

func (s *Scanport) Run() {
	s.pool, _ = pool.NewPool(100)
	s.getAllIP()
	s.getAllPort()
	s.scan()
}

func (s *Scanport) scan() {
	for _, ip := range s.IPs {
		for _, port := range s.Ports {
			// 协称池
			wg.Add(1)
			task := &pool.Task{
				Handler: func(v ...interface{}) {
					defer wg.Done()
					if isOpen("tcp", v[0].(string), v[1].(int)) {
						s.ResChan <- fmt.Sprintf("%v:%v", v[0], v[1])
					}
				},
				Params: []interface{}{ip, port},
			}

			s.pool.Put(task)
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
		fmt.Println("开放端口:", item)
	}
}

