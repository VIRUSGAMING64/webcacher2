package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
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
		queue.UGQueue.Push(queue.NewObj(resp.Request))
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
	//* hay que devolver el body al response, si no el dump queda vacio
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(data))
	cache.Global.Push(resp.Request, resp)
	return resp
}

func OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	wdebug.Log("request to: ", req.URL.String())

	if req.Method != "GET" || urlutils.Parse(req) == "" {
		return req, nil
	} else if req.Header.Get("webcacher-update") == "true" {
		return req, nil
	}

	resp, err := cache.Global.Pop(req)

	if err == nil && !config.Global.NoCache {

		if err != nil {
			return req, nil
		}
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
		data := ReadHTML("public/nointernet.html")
		return req, goproxy.NewResponse(req, goproxy.ContentTypeHtml, http.StatusServiceUnavailable, data)

	}

	return req, nil

}

func RunProxy() {
	fd, err := os.OpenFile("/tmp/webcacher-stderr.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	os.Stderr = fd
	if err != nil {
		fmt.Println(err)
	}
	go MainWork()
	config.ParseArgs()
	Proxy := goproxy.NewProxyHttpServer()
	Proxy.OnRequest().HandleConnect(ConnectHandler)
	Proxy.OnRequest().DoFunc(OnRequest)
	Proxy.OnResponse().DoFunc(OnResponse)
	fmt.Println("Listening on 0.0.0.0:8092")
	fmt.Println(http.ListenAndServe(":8092", Proxy))
}
