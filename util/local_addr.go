package util

import (
	"net"
	"os"
	"sync"
)

var (
	once      sync.Once
	localAddr string
)

func GetLocalAddr() string {
	once.Do(func() {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			os.Exit(1)
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			ipv4 := ipnet.IP.To4()
			if ipv4 != nil {
				localAddr = ipv4.String()
			}
		}
	})
	return localAddr
}
