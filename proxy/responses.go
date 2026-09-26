package proxy

import (
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
)

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

	return resp
}

func OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	if req.Method != "GET" {
		return req, nil
	}

	return req, nil
}
