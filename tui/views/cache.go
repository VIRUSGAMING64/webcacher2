package views

import (
	"strings"
	"time"
	"webcacher2/cache"
	"webcacher2/proxy"
)

func RenderCache() string {

	s := proxy.Pstats.Copy()
	cache.Global.Memcache.Mtx.Lock()
	cache.Global.Memcache.Mtx.Unlock()
	rows := []string{
		metric("Total en Cache", formatBytes(s.Total)),
		metric("Caché en memoria", formatBytes(cache.Global.Memcache.Size)),
		metric("Uso de caché", formatBytes(s.CacheUse)),
		metric("Last update", snapshotTime(time.Now())),
	}
	return PanelStyle.Render(strings.Join(rows, "\n"))
}
