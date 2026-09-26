package config

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
