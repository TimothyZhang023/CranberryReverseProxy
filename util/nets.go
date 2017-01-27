package util

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetRemoteAddress(r *http.Request) string {
	ip := strings.ToLower(r.Header.Get("x-forwarded-for"))
	iplen := len(ip)
	localIp := strings.ToLower("127.0.0.1")
	unKnown := "unknown"
	if ip == "" || iplen == 0 || ip == localIp || unKnown == ip {
		ip = r.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || iplen == 0 || ip == localIp || unKnown == ip {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || iplen == 0 || ip == localIp || unKnown == ip {
		ip = r.RemoteAddr
	}
	return ip

}

func GetLocalAddressByName(name string) (string, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, inter := range interfaces {

		if strings.HasPrefix(inter.Name, name) {
			fmt.Printf("found network card : %s \n", inter.Name)

			addrs, _ := inter.Addrs()

			for _, value := range addrs {
				ips := strings.Split(value.String(), "/")

				for _, item := range ips {

					ip := net.IP{}

					ip.UnmarshalText([]byte(item))

					res := ip.To4()

					if res != nil {
						return item, nil
					}

				}

			}

		}

	}

	return "", errors.New("no available addr found")

}

func GetLocalAddress() string {

	var addr string
	var err error

	netCards := []string{"bond", "eth", "em", "Internet", "WLAN", ""}

	for _, netCard := range netCards {
		addr, err = GetLocalAddressByName(netCard)

		if err == nil {
			break
		}
	}

	return addr

}

func ValidIpV4Addr(addr string) bool {
	ip := net.IP{}

	ip.UnmarshalText([]byte(addr))

	res := ip.To4()

	if res != nil {
		return true
	}
	return false
}
