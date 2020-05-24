package main

import (
	"fmt"
	"net"
)

func main() {

	// 创建连接
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, 104),
		Port: 8888,
	})

	if err != nil {
		fmt.Println("连接失败!", err)
		return
	}
	defer socket.Close()

	// 发送数据
	sendData := []byte("hello server!")
	fmt.Println(sendData)
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("发送数据失败!", err)
		return
	}
	for {
		// 接收数据
		data := make([]byte, 4096)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("读取数据失败!", err)
			return
		}
		fmt.Println(read, remoteAddr)
		fmt.Printf("%s\n", data)
	}
}
