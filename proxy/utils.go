package proxy

import (
	"io/fs"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func ReadHTML(name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		return "No internet"
	}
	return string(data)
}

// {#685, 7} Esto esta hecho con copilot
func mustParseURL(raw string) *url.URL {
	parsed, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return parsed
}

var Internet bool = HasInternet()

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

func InternetChecker() {
	for {
		Internet = HasInternet()
		time.Sleep(time.Second * 3)
	}
}

func CacheSize(folder string) int64 {
	sz := int64(0)
	filepath.Walk(
		folder, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			sz += info.Size()
			return err
		},
	)
	return sz
}
