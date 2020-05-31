## 文件传输服务端
[文件传输客户端](https://github.com/seth-shi/file-sync-client)

## 流程图
```
+-----------------------------------------------+
|                                               |            2. connect to tcp,send auth code Verifies identity
|                                         +-----+------+  <-------------------------------------------------------+
|                                         |            |                                                          |
|                                         |            |     4. each file is transferred using a TCP connection   |
|                                         |            |  <----------------------------------------------------+  |
|               +-----------------------  |            |     The files to be uploaded include:                 |  |
|               |                         |            |     * File not marked as uploaded                     |  |
|               | 5. store file by        |   server   |     * The file mtime is longer than the last uptime
|               | filepath and filename   |            |     * custom filter                                   |  |
|               | receive eof             |            |                                                       |  |
|               | Verify file integrity   |            |     Send the steps                                    +  |
|               |                         |            |     send md5 filename, filepath, filesize
|               +---------------------->  |            |     send chunk file content loop            -------------+
|                                         |            |     send end of file                        |            |
|                                         +------------+                                             |            |
|                                                                                                    |            |
|                                         | | +       6. close tcp connect                           |            |
|                                         | | +----------------------------------------------------> |   client   |
v                                         | |                                                        |            |
+------------+                            | |         3. auth success, begin transfer files          |            |
|            |                            | +------------------------------------------------------> |            |
|   client   |                            |                                                          |            |
|            |                            |           1. send udp broad.content is tcp port          |            |
+------------+                            +--------------------------------------------------------> +------------+

```

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