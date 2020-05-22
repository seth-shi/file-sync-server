## 文件传输服务端
[文件传输客户端](https://github.com/seth-shi/file-sync-client)

```shell script
git clone https://github.com/seth-shi/file-sync-server
cd flash-sync-server
go build
```

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