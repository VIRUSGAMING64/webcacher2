package proxy

import (
	"net"
	"net/url"
	"time"
)

// {#685, 7} Esto esta hecho con copilot
func mustParseURL(raw string) *url.URL {
	parsed, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return parsed
}

func HasInternet() bool {
	var Dns []string = []string{"1.1.1.1:53", "8.8.8.8:83"}

	for _, dns := range Dns {
		_, err := net.DialTimeout("tcp", dns, time.Second*1)
		if err == nil {
			return true
		}
	}

	return false

}
