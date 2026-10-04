package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"webcacher2/cache"
	"webcacher2/config"
	wdebug "webcacher2/debug"
	"webcacher2/queue"
	"webcacher2/urlutils"

	"github.com/elazarl/goproxy"
)

var Proxy *goproxy.ProxyHttpServer = nil
var Parents []*http.Transport = make([]*http.Transport, 0)

var ConnectHandler goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.MitmConnect, host
}

func BuildProxyTransport(host string) (*http.Transport, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
	}

	if strings.TrimSpace(host) == "" {
		return transport, nil
	}

	parentProxy, err := parseParentProxy(host)
	if err != nil {
		return nil, err
	}

	transport.Proxy = http.ProxyURL(parentProxy)
	return transport, nil
}

func parseParentProxy(rawValue string) (*url.URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawValue))
	if err != nil {
		return nil, fmt.Errorf("invalid Pproxy value: %w", err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid Pproxy value: missing scheme or host")
	}

	return parsed, nil
}

var Pmtx sync.Mutex
var CurrParent int = 0

func IterateParentProxy() {
	Pmtx.Lock()
	defer Pmtx.Unlock()
	if len(Parents) == 0 {
		return
	}
	Proxy.Tr = Parents[CurrParent]
	CurrParent = (1 + CurrParent) % len(Parents)
}

func OnResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Request == nil {
		return resp
	}
	if resp.Header.Get("webcacher") == "true" {
		//	AddPayload(resp)
		queue.UGQueue.Push(queue.NewObj(resp.Request))
		Pstats.AddHint(resp)
		return resp
	}
	if resp.Request.Method != "GET" || resp.StatusCode != 200 {
		return resp
	}

	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	//* hay que devolver el body al response, si no el dump queda vacio
	resp.Body = io.NopCloser(bytes.NewReader(data))
	cache.Global.Push(resp.Request, resp)
	Pstats.AddBypass(resp)
	return resp
}

func OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	go IterateParentProxy()
	wdebug.Log("request to: ", req.URL.String())

	if req.Method != "GET" || urlutils.Parse(req) == "" {
		return req, nil
	} else if req.Header.Get("webcacher-update") == "true" {
		//	return req, nil
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
	Proxy = goproxy.NewProxyHttpServer()
	Proxy.OnRequest().HandleConnect(ConnectHandler)
	Proxy.OnRequest().DoFunc(OnRequest)
	Proxy.OnResponse().DoFunc(OnResponse)
	fmt.Println(http.ListenAndServe(":8092", Proxy))
}
