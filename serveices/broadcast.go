package serveices

import (
	. "flash-sync-server/global"
	"fmt"
	"log"
	"net"
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

	log.Printf("start udp broadcast, udp port: %d", udpPort)

	for _ = range ticker.C {

		// 广播自己的 tcp 端口
		msg := fmt.Sprintf("hello ! my tcp port=[%d]", tcpPort)
		_, err := broadcast.WriteToUDP([]byte(msg), dstAddr)
		log.Printf(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
