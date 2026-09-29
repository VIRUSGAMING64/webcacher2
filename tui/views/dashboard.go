package views

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"webcacher2/proxy"
	"webcacher2/queue"
)

func metric(name, value string) string {
	return fmt.Sprintf("%-22s %s", name, ValueStyle.Render(value))
}

func RenderDashboard() string {

	s := proxy.Pstats.Copy()
	w := queue.GQueue.Workers

	rows := []string{
		metric("Descargado", formatBytes(s.Downloaded)),
		metric("Bypass", strconv.Itoa(s.Bypass)),
		metric("Hits", strconv.Itoa(s.Hints)),
		metric("Workers", strconv.Itoa(w)),
		metric("Requests por minuto", strconv.Itoa(len(s.History))),
		metric("Last update", snapshotTime(time.Now())),
	}
	return PanelStyle.Render(strings.Join(rows, "\n"))
}
