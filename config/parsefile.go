package config

import (
	"os"
	"strings"
)

func ReadConfig(file string) (*Config, error) {

	tdata, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config = Config{}
	var data = strings.Split(string(tdata), "\n")

	for i := 0; i < len(data); i++ {
		if len(data[i]) == 0 {
			continue
		}
		var action = data[i]
		if strings.HasPrefix(action, "#") {
			continue
		}
		if strings.HasSuffix(action, ":") {
			var temp_arr []string
			for j := i + 1; j < len(data); j++ {
				if len(data[j]) == 0 {
					continue
				}
				if strings.HasPrefix(data[j], "#") {
					continue
				}
				if strings.HasSuffix(data[j], ":") {
					break
				}
				temp_arr = append(temp_arr, data[j])
				i++
			}
			if action == "NoCacheSite:" {
				config.NoCacheSites = temp_arr
			} else if action == "NoArgs:" {
				config.NoArgs = temp_arr
			} else if action == "NoCacheExt:" {
				config.NoCacheExt = temp_arr
			} else if action == "Pproxy:" && len(temp_arr) > 0 {
				config.Pproxy = temp_arr
			} else if action == "Syncs:" {
				config.Syncs = temp_arr
			} else if action == "CacheArgs:" {
				config.CacheArgs = temp_arr
			} else if action == "ProxyAddr:" && len(temp_arr) > 0 {
				config.ProxyAddr = temp_arr[0]
			} else if action == "DashboardAddr:" && len(temp_arr) > 0 {
				config.DashboardAddr = temp_arr[0]
			} else if action == "CacheDir:" && len(temp_arr) > 0 {
				config.CacheDir = temp_arr[0]
			} else if action == "MetricsFile:" && len(temp_arr) > 0 {
				config.MetricsFile = temp_arr[0]
			} else if action == "LogLevel:" && len(temp_arr) > 0 {
				config.LogLevel = temp_arr[0]
			} else if action == "IgnoreQueue:" && len(temp_arr) > 0 {
				config.IgnoreQueue = temp_arr
			}
		}
	}
	return &config, nil
}
