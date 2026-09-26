package proxy

import (
	"net"
	"time"
)

func HasInternet() bool {
	var Dns []string = []string{"1.1.1.1:53", "8.8.8.8:83"}

	for _, dns := range Dns {
		_, err := net.DialTimeout("tcp", dns, time.Second*10)
		if err == nil {
			return true
		}
	}

	return false

}
