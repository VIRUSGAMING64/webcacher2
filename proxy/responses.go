package proxy

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"webcacher2/cache"
	wdebug "webcacher2/debug"
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
	if resp.Request.Method != "GET" || resp.StatusCode != 200 {
		return resp
	}
	if resp.Header.Get("webcacher") == "true" {
		return resp
	}
	_, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp
	}
	//uri := urlutils.Parse(resp.Request.Method + resp.Request.URL.String())

	return resp
}

func OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	wdebug.Log("request from: ", req.URL.String())
	uri := urlutils.Parse(req)

	if req.Method != "GET" || uri == "" {
		return req, nil
	}

	data, err := cache.Global.Pop(uri)
	if err != nil {
		wdebug.Error("cache error", err)
		return req, nil
	}
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), req)
	resp.Header.Add("webcacher", "true")
	if err != nil {
		return req, nil
	}
	return req, resp
}
