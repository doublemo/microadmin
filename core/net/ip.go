package net

import (
	"errors"
	"net"
)

// LocalIP 获取本地IP地址
func LocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ip := ipNet.IP.To4()
				if ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 {
					return ip, nil
				}
			}
		}
	}

	return nil, errors.New("Failed to get ip address")
}
