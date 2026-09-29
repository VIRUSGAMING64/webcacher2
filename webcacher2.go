package main

import (
	"fmt"
	"os"
	"os/signal"
	"webcacher2/cache"
	"webcacher2/config"
	"webcacher2/proxy"
	"webcacher2/queue"
	"webcacher2/tui"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)

	go func() {
		for s := range c {
			queue.GQueue.Save("queue.json")
			cache.Global.Save()

			if s == os.Kill || s == os.Interrupt {
				os.Exit(0)
			}
		}
	}()

	go proxy.InternetChecker()
	config.Global.Watcher()
	cache.Global = cache.NewUrlCache()
	cache.Global.Load()

	proxy.Pstats.Load("stats.json")
	fmt.Println("Calculating size")
	proxy.Pstats.Total = proxy.CacheSize(".cache/")
	fmt.Println("Size:", proxy.Pstats.Total)
	queue.GQueue.Load("queue.json")
	fmt.Println("Loaded queue with size: [", queue.GQueue.Length(), "]")

	go func() {
		if err := tui.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		queue.GQueue.Save("queue.json")
		cache.Global.Save()
		proxy.Pstats.Save("stats.json")
		os.Exit(0)
	}()

	proxy.RunProxy()
}
