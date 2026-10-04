package views

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"webcacher2/config"
	"webcacher2/proxy"
	"webcacher2/queue"
)

func metric(name, value string) string {
	return fmt.Sprintf("%-22s %s", name, ValueStyle.Render(value))
}

func RenderDashboard() string {
	proxy.Pmtx.Lock()
	defer proxy.Pmtx.Unlock()

	s := proxy.Pstats.Copy()
	w := queue.GQueue.Workers

	rows := []string{
		metric("Descargado", formatBytes(s.Downloaded)),
		metric("Online", strconv.FormatBool(proxy.Internet)),
		metric("Bypass", ValueStyle.Render(strconv.Itoa(s.Bypass))),
		metric("Hits", strconv.Itoa(s.Hints)),
		metric("Workers", strconv.Itoa(w)),
		metric("Requests por minuto", strconv.Itoa(len(s.History))),
	}
	if len(proxy.Parents) > 0 {
		rows = append(rows, metric("Current Parent", ValueStyle.Render(config.Global.Pproxy[proxy.CurrParent])))
	}
	rows = append(rows, metric("Last update", snapshotTime(time.Now())))
	return PanelStyle.Render(strings.Join(rows, "\n"))
}
