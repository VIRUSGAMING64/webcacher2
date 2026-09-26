package urlutils

import (
	"fmt"
	"net/http"
	"strings"
	"webcacher2/config"
)

func Parse(req *http.Request) string {
	//* Aqui es donde se aplican las config
	uri := req.URL.String()
	ext := Extension(uri)
	for _, elem := range config.Global.NoCacheSites {
		fmt.Println(req.Host, elem)
		if elem == req.Host {
			return ""
		}
	}
	for _, elem := range config.Global.NoCacheExt {
		if elem == ext {
			return ""
		}
	}
	for _, elem := range config.Global.NoArgs {
		if elem == req.Host {
			return req.Method + CutArgs(uri)
		}
	}

	return req.Method + uri
}

func Extension(uri string) string {
	strs := strings.Split(CutArgs(uri), ".")
	if len(strs) <= 1 {
		return ""
	}
	ext := strs[len(strs)-1]
	if len(ext) >= 5 {
		ext = ""
	}
	return ext
}

func GetArgs(uri string) map[string]string {
	Arrargs := strings.Split(uri, "?")
	ma := make(map[string]string)
	if len(Arrargs) <= 1 {
		return ma
	}
	a := strings.Split(Arrargs[len(Arrargs)-1], "&")
	for _, elem := range a {
		args := strings.Split(elem, "=")
		if len(args) <= 1 {
			continue
		}
		ma[args[0]] = args[1]
	}
	return ma
}

func CutArgs(uri string) string {
	Arrargs := strings.Split(uri, "?")
	if len(Arrargs) <= 1 {
		return uri
	}
	return Arrargs[0]
}
