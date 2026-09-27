package urlutils

import (
	"net/http"
	"path/filepath"
	"strings"
	"webcacher2/config"
)

const maxExtLen = 10

func Key(req *http.Request) string {
	if req == nil {
		return ""
	}
	if config.Global == nil || !config.Global.NoArgsMode() {
		return req.Method + req.URL.String()
	}

	cacheargs := true
	for _, elem := range config.Global.NoArgs {
		if elem == req.Host {
			cacheargs = false
		}
	}

	uri := req.URL.String()
	parts := strings.Split(uri, "?")
	path := parts[0]

	if len(parts) > 1 && len(config.Global.CacheArgs) >= 1 {
		sep := "?"
		for _, arg := range strings.Split(parts[1], "&") {
			kv := strings.Split(arg, "=")
			if len(kv) <= 1 {
				continue
			}
			if isIn(kv[0], config.Global.CacheArgs) {
				if sep != "?" {
					sep += "&"
				}
				sep = sep + kv[0] + "=" + kv[1] //!TODO Aqui es donde se debe cambiar los valores de los arguments
			}
		}
		path = path + sep
	}

	if cacheargs || len(parts) <= 1 {
		return req.Method + req.URL.String()
	}

	return req.Method + path
}

func Parse(req *http.Request) string {
	if req == nil {
		return ""
	}
	if config.Global != nil {
		ext := Extension(req)
		for _, elem := range config.Global.NoCacheSites {
			if elem == req.Host {
				return ""
			}
		}
		for _, elem := range config.Global.NoCacheExt {
			if elem == ext {
				return ""
			}
		}
	}
	return Key(req)
}

func Extension(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	ext := filepath.Ext(req.URL.Path)
	if ext == "" {
		return ""
	}
	ext = strings.TrimPrefix(ext, ".")
	if len(ext) > maxExtLen {
		return ""
	}
	return ext
}

func validExt(ext string) bool {
	if ext == "" {
		return true
	}
	return !strings.ContainsAny(ext, `/\`)
}

func isIn(elem string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if elem == arr[i] {
			return true
		}
	}
	return false
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
