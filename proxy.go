package main

import (
	"fmt"
	"webcacher2/cache"
	"webcacher2/config"
	"webcacher2/proxy"
	"webcacher2/queue"
)

func main() {
	conf, err := config.ReadConfig("webcacher.conf")
	if err != nil {
		panic(err)
	}
	proxy.Pstats.Load("stats.json")
	fmt.Println("Calculating size")
	proxy.Pstats.Total = proxy.CacheSize(".cache/")
	fmt.Println("Size:", proxy.Pstats.Total)
	queue.GQueue.Load("queue.json")
	fmt.Println("Loaded queue with size: [", queue.GQueue.Length(), "]")
	config.Global = conf
	cache.Global = cache.NewUrlCache()
	proxy.RunProxy()
}
