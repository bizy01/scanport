### Info
scanport是一个使用go语言编写的端口扫描程序，内部实现协称池，支持各种灵活的配置

### Feature
- 支持tcp/upd端口
- 扫描目标支持灵活多变的配置
- 端口范围

### Usage
scanport -h

<生成默认配置文件>
  usage: -init

<以加载配置文件方式运行>
  usage: -c

<扫描的协议>
  usage: tcp,udp

<扫描的目标主机，支持ip, 域名，cidr>
  default: 127.0.0.1
  usage:
  (1): 127.0.0.1
  (2): 192.168.0.1,192.168.0.2
  (3): 192.168.0.1-20
  (4): www.baidu.com
  (5): 192.168.1.1/30

<端口值>
  default: 80
  usage:
  (1): 3000, 8080, 3306
  (2): 3000-10000
  (3): 8080,3000-10000

<扫描并发数>
  default: 1000

<Dial timeout(unit Millisecond)>
  default: 100 (Millisecond)
