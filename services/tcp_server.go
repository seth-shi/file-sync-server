package services

import (
	"errors"
	. "flash-sync-server/global"
	"net"
	"strconv"

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

		msg := App.I18n.Format("received tcp connect message", plurr.Params{"address": conn.RemoteAddr(), "local": conn.LocalAddr()})
		LogInfo(msg)
	}
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
