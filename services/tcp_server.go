package services

import (
	. "flash-sync-server/global"
	"github.com/iafan/Plurr/go/plurr"
	"net"
	"strconv"
)


func StartTcpServer()  {

	tcpPort := strconv.Itoa(App.Config.Tcp.Port)


	server, err := net.Listen("tcp", ":" + tcpPort)
	if err != nil {
		LogErrorHandle(err)
		return
	}
	defer server.Close()


	ip, err := getIp()
	if err != nil {
		LogError(App.I18n.Tr("get ip fail"))
		return
	}

	LogInfo(App.I18n.Format("start tcp server", plurr.Params{"address": ip + ":" + tcpPort}))

	for {
		conn, err := server.Accept()
		if err != nil {
			LogErrorHandle(err)
			continue
		}

		msg := App.I18n.Format("received tcp connect message", plurr.Params{"address": conn.RemoteAddr(), "local": conn.LocalAddr()})
		LogInfo(msg)
	}
}

func getIp() (string, error) {

	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return "", nil
	}

	for _, address := range addresses {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				return ipnet.IP.String(), nil
			}

		}
	}

	return "", nil
}