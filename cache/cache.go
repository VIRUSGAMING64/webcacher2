package cache

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"sync"
	"webcacher2/urlutils"
)

var ErrMiss = errors.New("cache: miss")

type MemoryCache struct {
	cache map[string][]byte
	mut   sync.Mutex
}

func (m MemoryCache) Push(name string, data []byte) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.cache[name] = data
}
func (m MemoryCache) Pop(name string) ([]byte, error) {
	m.mut.Lock()
	defer m.mut.Unlock()
	data := m.cache[name]
	var err error = nil
	if len(data) == 0 {
		err = ErrMiss
	}
	return data, err
}

type UrlCache struct {
	SavePath string            `json:"savepath"`
	Size     int64             `json:"size"`
	Counts   map[string]int64  `json:"count"`
	Hashes   map[string]string `json:"hashes"`
	MemKeys  map[string]string `json:"keys"`
	memcache MemoryCache
	mtx      sync.Mutex
}

var Global *UrlCache

func (c *UrlCache) init() {
	c.memcache.cache = make(map[string][]byte)
	if c.SavePath == "" {
		c.SavePath = "cache.json"
	}
	if c.Counts == nil {
		c.Counts = make(map[string]int64)
	}
	if c.Hashes == nil {
		c.Hashes = make(map[string]string)
	}
	if c.MemKeys == nil {
		c.MemKeys = make(map[string]string)
	}
}

func (c *UrlCache) Load() error {
	data, err := os.ReadFile(c.SavePath)
	defer c.init()

	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

// * <sha256 de METHOD+url>[.ext].phttp dentro de la carpeta del cache
func (c *UrlCache) Push(req *http.Request, resp *http.Response) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if req == nil || resp == nil {
		return false
	}
	key := urlutils.Key(req)
	if key == "" {
		return false
	}

	ext := urlutils.Extension(req)
	file := c.CachePath(key, ext)
	body, err := httputil.DumpResponse(resp, true)

	if err != nil {
		return false
	}

	curr := int64(0)
	inf, err := os.Stat(file)
	if err == nil {
		curr = inf.Size()
	}
	if err = os.MkdirAll(".cache", 0755); err != nil {
		return false
	}
	if err = os.WriteFile(file, body, 0644); err != nil {
		return false
	}

	if c.Counts[key] == 0 {
		curr = 0
	}
	c.Size += int64(len(body)) - curr
	c.Hashes[key] = Hash(key)
	c.MemKeys[file] = key
	c.memcache.Push(file, body)
	return true
}

// * Pop busca el archivo probando las claves que pudo haber usado webcacher1
func (c *UrlCache) Pop(req *http.Request) (*http.Response, error) {
	if req == nil || urlutils.Parse(req) == "" {
		return nil, ErrMiss
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()
	key := urlutils.Parse(req)
	ext := urlutils.Extension(req)
	file := c.CachePath(key, ext)
	data, err := c.memcache.Pop(file)
	mem := true
	if err != nil {
		mem = false
		data, err = os.ReadFile(file)
	} else {
		fmt.Println("memory hint", req.URL.String())
	}
	if err != nil {
		return nil, ErrMiss
	}
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), req)
	if err != nil {
		return nil, ErrMiss
	}
	if mem == false {
		c.memcache.Push(file, data)
	}
	c.Counts[key] += 1
	c.MemKeys[file] = key
	return resp, nil

}

func NewUrlCache() *UrlCache {
	c := UrlCache{}
	c.init()
	return &c
}

func (c *UrlCache) Save() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	data := c.Json()
	fd, err := os.OpenFile(c.SavePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = fd.Write(data)
	if err != nil {
		return err
	}
	return fd.Close()
}

func (c *UrlCache) Json() []byte {
	data, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return []byte{}
	}
	return data
}

func Hash(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (c *UrlCache) CachePath(key string, ext string) string {
	fil := Hash(key)
	if ext != "" {
		fil += "." + ext
	}
	fil += ".phttp"

	return filepath.Join(".cache", fil)
}
