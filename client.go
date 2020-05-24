package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	// 这里设置接收者的IP地址为广播地址
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8888,
	})

	if err != nil {
		println(err.Error())
		return
	}

	buffSize := 4096
	err = conn.SetReadBuffer(buffSize)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("staring listen")
	for {
		// 接收数据
		data := make([]byte, buffSize)
		read, address, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
		}

		log.Println("接收到消息", read, address, string(data))
		time.Sleep(time.Second)
	}
}
