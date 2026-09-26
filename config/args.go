package config

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("WebCacher - HTTP proxy for offline caching")
	fmt.Println("")
	fmt.Println("Usage: webcacher [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --no-all            Only cache files with recognized extensions")
	fmt.Println("  --no-args           Cache files excluding url arguments")
	fmt.Println("  --no-cache          Bypass cache for all requests")
	fmt.Println("  --no-queue          Disable offline request queue")
	fmt.Println("  --proxy-addr ADDR   Proxy listen address")
	fmt.Println("  --dashboard-addr ADDR Dashboard listen address")
	fmt.Println("  --cache-dir DIR     Cache directory")
	fmt.Println("  --metrics-file FILE Metrics file path")
	fmt.Println("  --log-level LEVEL   Log level")
	fmt.Println("  --help              Show this message")
	fmt.Println("  --no-mem-cache      Only use cache in files")
}

func ParseArgs() {
	ParseArgsFrom(os.Args[1:])
}

func ParseArgsFrom(args []string) {
	supportedFlags := map[string]bool{
		"--no-all":         true,
		"--no-args":        true,
		"--help":           true,
		"--no-cache":       true,
		"--no-queue":       true,
		"--proxy-addr":     true,
		"--dashboard-addr": true,
		"--cache-dir":      true,
		"--metrics-file":   true,
		"--log-level":      true,
		"--no-mem-cache":   true,
	}

	for i := 0; i < len(args); i++ {
		flag := args[i]
		if !supportedFlags[flag] {
			help()
			os.Exit(1)
		}

		switch flag {
		case "--no-args":
			Global.args = true
		case "--no-all":
			Global.NoAll = true
		case "--no-cache":
			Global.NoCache = true
		case "--no-mem-cache":
			Global.NoMemoryCache = true
		case "--no-queue":
			Global.NoQueue = true
		case "--help":
			help()
			os.Exit(0)
		}
	}
}
