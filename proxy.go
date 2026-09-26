package main

import (
	"net/http"
	"webcacher2/cache"
	"webcacher2/config"
	"webcacher2/proxy"

	"github.com/elazarl/goproxy"
)

func Run() {
	Proxy := goproxy.NewProxyHttpServer()
	Proxy.OnRequest().HandleConnect(proxy.ConnectHandler)
	Proxy.OnRequest().DoFunc(proxy.OnRequest)
	Proxy.OnResponse().DoFunc(proxy.OnResponse)
	http.ListenAndServe(":8092", Proxy)
}

func main() {
	conf, err := config.ReadConfig("webcacher.conf")
	if err != nil {
		panic(err)
	}
	config.Global = conf
	cache.Global = cache.NewUrlCache()
	Run()
}
