package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path"
	"sync"
	"webcacher2/urlutils"
)

type UrlCache struct {
	SavePath string            `json:"savepath"`
	Folder   string            `json:"folder"`
	Size     int64             `json:"size"`
	Counts   map[string]int64  `json:"count"`
	Hashes   map[string]string `json:"hashes"`
	MemKeys  map[string]string `json:"keys"`
	mtx      sync.Mutex
}

var Global *UrlCache

func (c *UrlCache) Load() error {
	data, err := os.ReadFile(c.SavePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c)
}

func (c *UrlCache) Push(url string, body []byte) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var err error = nil

	path := c.CachePath(url)
	inf, err := os.Stat(path)
	curr := int64(0)
	if err == nil {
		curr = inf.Size()
	}
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		return false
	}
	if c.Counts[url] == 0 {
		curr = 0
	}
	c.Size += int64(len(body)) - curr
	c.MemKeys[path] = url

	return true
}

func (c *UrlCache) Pop(url string) ([]byte, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var err error = nil

	pth := c.CachePath(url)
	data, err := os.ReadFile(pth)
	if err != nil {
		return nil, err
	}
	c.Counts[url] += 1

	return data, err
}

func NewUrlCache() *UrlCache {
	c := UrlCache{}
	c.MemKeys = make(map[string]string)
	c.Counts = make(map[string]int64)
	c.Hashes = make(map[string]string)
	c.Folder = ".cache"
	c.SavePath = "cache.json"
	return &c
}

func (c *UrlCache) Save() error {
	data := c.Json()
	fd, err := os.OpenFile(c.SavePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = fd.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *UrlCache) Json() []byte {
	data, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return []byte{}
	}
	return data
}

func (c *UrlCache) CachePath(url string) string {
	ext := urlutils.Extension(url)
	hash := sha256.Sum256([]byte(url))
	fil := hex.EncodeToString(hash[:])

	if ext != "" {
		fil += "." + ext
	}
	fil += ".phttp"

	return path.Join(c.Folder, fil)
}
