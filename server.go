package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	// 这里设置接收者的IP地址为广播地址
	conn, err := net.DialUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8888,
	}, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 8888,
	})

	if err != nil {
		println(err.Error())
		return
	}

	for {
		_, err := conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println(err)
		}
		log.Println("server send")
		time.Sleep(time.Second)
	}
}
