package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"webcacher2/cache"
	"webcacher2/config"
	wdebug "webcacher2/debug"
	"webcacher2/queue"
	"webcacher2/urlutils"

	"github.com/elazarl/goproxy"
)

var ConnectHandler goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.MitmConnect, host
}

func OnResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Request == nil {
		return resp
	}
	if resp.Header.Get("webcacher") == "true" {
		Pstats.AddHint(resp)
		return resp
	}
	defer Pstats.AddBypass(resp)
	if resp.Request.Method != "GET" || resp.StatusCode != 200 {
		return resp
	}
	data, err := io.ReadAll(resp.Body)
	fmt.Println("Downloaded", len(data), "bytes from", resp.Request.URL.String())
	if err != nil {
		return resp
	}
	uri := urlutils.Parse(resp.Request)
	cache.Global.Push(uri, data)
	return resp
}

func OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	wdebug.Log("request to: ", req.URL.String())
	uri := urlutils.Parse(req)

	if req.Method != "GET" || uri == "" {
		return req, nil
	}

	data, err := cache.Global.Pop(uri)

	if err == nil && !config.Global.NoCache {
		resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), req)
		resp.Header.Add("webcacher", "true")
		if err != nil {
			return req, nil
		}
		return req, resp
	}

	if !Internet {
		flag := req.Header.Get("webcacher-queue") != "true"
		for _, elem := range config.Global.IgnoreQueue {
			if elem == req.Host {
				flag = false
				break
			}
		}
		if flag {
			queue.GQueue.Push(queue.NewObj(req))
		}
	}

	return req, nil

}

func RunProxy() {
	go MainWork()
	config.ParseArgs()
	Proxy := goproxy.NewProxyHttpServer()
	Proxy.OnRequest().HandleConnect(ConnectHandler)
	Proxy.OnRequest().DoFunc(OnRequest)
	Proxy.OnResponse().DoFunc(OnResponse)
	fmt.Println("Listening on 0.0.0.0:8092")
	fmt.Println(http.ListenAndServe(":8092", Proxy))
}
