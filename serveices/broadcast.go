package serveices

import (
	. "flash-sync-server/global"
	"flash-sync-server/models"
	"github.com/iafan/Plurr/go/plurr"
	"net"
	"strconv"
	"time"
)

func SendConnectUdpPack(ticker *time.Ticker) {

	udpPort, tcpPort := App.Config.Udp.Port, App.Config.Tcp.Port

	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: udpPort}

	broadcast, err := net.ListenUDP("udp", srcAddr)
	if err != nil {

		panic(err)
	}

	App.LogChan <- models.InfoLog(App.I18n.Format("start udp broadcast, udp port: {port}", plurr.Params{"port": tcpPort}))


	for _ = range ticker.C {

		// 广播自己的 tcp 端口
		_, err := broadcast.WriteToUDP([]byte(strconv.Itoa(tcpPort)), dstAddr)

		if err != nil {

			App.LogChan <- models.ErrorLog(err.Error())
		}
	}
}
