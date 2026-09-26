package proxy

import (
	"io"
	"net/http"
	"webcacher2/cache"
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
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp
	}
	uri := urlutils.Parse(resp.Request.URL)
	cache.Global.Push(uri, bodyBytes)

	return resp
}

func OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	Log("request from: ", req.URL.String())
	if req.Method != "GET" {
		return req, nil
	}

	return req, nil
}
