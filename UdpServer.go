package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func main() {
	// 这里设置接收者的IP地址为广播地址
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 8888,
	})
	if err != nil {
		println(err.Error())
		return
	}

	for {
		conn.Write([]byte("say hello"))
		fmt.Println("广播包")
		time.Sleep(3 * time.Second)
	}
	conn.Close()
}

func GetIntranetIp() (*net.IPNet, error) {

	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return nil, errors.New("cannot get interface address")
	}

	for _, address := range addresses {

		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet, nil
			}

		}
	}

	return nil, errors.New("cannot get address")
}
