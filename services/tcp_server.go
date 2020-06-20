package services

import (
	"encoding/json"
	"errors"
	"flash-sync-server/enums"
	. "flash-sync-server/global"
	"net"
	"strconv"
	"strings"

	"github.com/iafan/Plurr/go/plurr"
)

func StartTcpServer() {

	tcpPort := strconv.Itoa(App.Config.Tcp.Port)

	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		LogErrorHandle(err)
		return
	}
	defer server.Close()

	ip, err := externalIP()
	if err != nil {

		LogError(err.Error())
		LogError(App.I18n.Tr("get ip fail"))
		return
	}

	LogInfo(App.I18n.Format("start tcp server", plurr.Params{"address": ip.String() + ":" + tcpPort}))

	for {
		conn, err := server.Accept()
		if err != nil {
			LogErrorHandle(err)
			continue
		}

		go handleFileSync(&TcpClient{conn, 1204})
	}
}

func handleFileSync(client *TcpClient) {

	msg := App.I18n.Format("received tcp connect message", plurr.Params{"address": client.conn.RemoteAddr()})
	LogInfo(msg)

	// 首次连接, 客户端会发送设备 id 过来
	deviceId := string(client.readContent())
	// 查看设备是否已经连接验证过
	_, exists := App.ClientDevices[deviceId]
	if !exists {

		// 进行身份验证
		_ = client.writeContents(enums.LinkCode)

		for {

			linkCode := strings.Trim(string(client.readContent()), " ")
			if linkCode == App.LinkCode {
				_ = client.writeContents(enums.LinkSuccess)

				// 验证成功加入到设备号当中
				App.ClientDevices[deviceId] = deviceId
				if err := App.DeviceDb.Put([]byte("devices-"+deviceId), []byte(deviceId), nil); err != nil {
					LogErrorHandle(err)
				}
				break
			}
			_ = client.writeContents(enums.LinkFail)
		}
	}

	var fileInfo FileInfo

	for {

		// 验证身份成功
		// 先发送文件名,文件大小.
		_ = client.writeContents(enums.FileInfo)
		err := json.Unmarshal(client.readContent(), &fileInfo)
		if err == nil {
			break
		}

		LogErrorHandle(err)
	}

	// 开始切片内容
}

func (client *TcpClient) readContent() []byte {

	//读取客户端发送的内容
	buf := make([]byte, client.bufSize)
	n, err := client.conn.Read(buf)

	if err != nil {
		LogErrorHandle(err)
		return []byte{}
	}

	return buf[:n]
}

func (client *TcpClient) writeContents(msg string) error {

	//读取客户端发送的内容
	_, err := client.conn.Write([]byte(msg))
	if err != nil {

		LogErrorHandle(err)
	}

	return err
}

func externalIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, face := range interfaces {
		if face.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if face.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addresses, err := face.Addrs()
		if err != nil {
			return nil, err
		}

		var ip net.IP
		for _, addr := range addresses {

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip, nil
		}
	}
	return nil, errors.New(App.I18n.Tr("connected to the network?"))
}

type TcpClient struct {
	conn    net.Conn
	bufSize int
}

type FileInfo struct {
	FileName string `json:"file_name"`
	FileSize int    `json:"file_size"`
}
