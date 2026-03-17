package ips

import (
	"fmt"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	RemoteIP := c.RemoteIP()
	if RemoteIP != "127.0.0.1" {
		ip = RemoteIP
	}
	ClientIP := c.ClientIP()
	if ClientIP != "127.0.0.1" {
		ip = ClientIP
	}
	return ip
}

// 获取局域网ip地址
func GetLocalHost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}
