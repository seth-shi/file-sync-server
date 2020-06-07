package services

import (
	"errors"
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

	msg := App.I18n.Format("received tcp connect message", plurr.Params{"address": conn.RemoteAddr()})
	LogInfo(msg)

	// 首次连接, 客户端会发送设备 id 过来
	deviceId := string(client.readContent())
	// 查看设备是否已经连接验证过
	_, exists := App.ClientDevices[deviceId]
	if !exists {

		// 进行身份验证
		client.conn.Write([]byte("link_code"))

		for {

			linkCode := strings.Trim(string(client.readContent()), " ")
			if linkCode == App.LinkCode {
				client.conn.Write([]byte("link_success"))
				App.ClientDevices[deviceId] = deviceId
				App.DeviceDb.Put([]byte("devices-"+deviceId), []byte(deviceId), nil)
				break
			}

			client.conn.Write([]byte("link_fail"))
		}
	}

	// 验证身份成功
	// 开始发送文件名和文件内容
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

type TcpClient struct {
	conn    net.Conn
	bufSize int
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
