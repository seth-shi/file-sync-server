## 文件传输服务端
[文件传输客户端](https://github.com/seth-shi/file-sync-client)

> 如果要让网络（同一网络）中的所有计算机都能收到这个数据包，就应该将这个数据包的接收者地址设置为这个网络中的最高的主机号。通常255.255.255.255就可以达到这个要求。所以我们如果要发送一次UDP广播报文，就可以试试如下实例代码：

```shell script
git clone https://github.com/seth-shi/file-sync-server
cd flash-sync-server
go build
```

## 为什么有此项目
* 刚学完`Golang`之前一直写`Web`,想尝试一下`GUI`
* 我电脑的[DAEMON SYNC](https://daemonsync.me/home)用不了了

## 编译选项
```
go build
# 不要黑窗口
go build -ldflags="-H windowsgui"
```

> CGo Optimizations
> The usual default message loop includes calls to win32 API functions, which incurs a decent amount of runtime overhead coming from Go. As an alternative to this, you may compile Walk using an optional C implementation of the main message loop, by passing the `walk_use_cgo` build tag:
```
go build -tags walk_use_cgo
```