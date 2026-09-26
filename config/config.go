package config

import "encoding/json"

type Config struct {
	IgnoreQueue   []string
	NoCacheSites  []string
	NoArgs        []string
	NoCacheExt    []string
	Syncs         []string
	CacheArgs     []string
	Pproxy        []string
	ProxyAddr     string
	DashboardAddr string
	CacheDir      string
	MetricsFile   string
	LogLevel      string
	args          bool
	NoQueue       bool
	NoMemoryCache bool
}

var Global *Config

func (c *Config) Json() []byte {
	data, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return []byte{}
	}
	return data
}
