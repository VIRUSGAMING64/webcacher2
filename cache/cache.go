package cache

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type UrlCache struct {
	SavePath string            `json:"savepath"`
	Folder   string            `json:"folder"`
	Size     int64             `json:"size"`
	Counts   map[string]int64  `json:"count"`
	Hashes   map[string]string `json:"hashes"`
	MemKeys  map[string]string `json:"keys"`
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
	var err error = nil
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
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

	fmt.Println(curr)
	c.MemKeys[path] = url

	return true
}

func (c *UrlCache) Pop(url string) ([]byte, error) {
	var err error = nil
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
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
	hash := sha512.Sum512([]byte(url))
	fil := hex.EncodeToString(hash[:])
	return path.Join(c.Folder, fil)
}
