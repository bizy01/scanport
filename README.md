### Info
scanport是一个使用go语言编写的端口扫描程序，内部实现协称池，支持各种灵活的配置

### Feature
- 支持tcp/upd端口
- 扫描目标支持灵活多变的配置
- 端口范围

### Usage
- help
![](images/usage.png)

- 结果
![](images/result.png)

- 容器用法
```
docker push bizy01/scanport:latest
docker run -it --rm scanport:latest scanport -target www.baidu.com -port "80,443"
```
