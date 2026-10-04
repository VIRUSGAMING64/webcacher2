package views

import (
	"fmt"
	"strings"
	"time"
	"webcacher2/proxy"
)

func RenderHist() string {
	s := proxy.Pstats.Copy()
	rows := []string{TitleStyle.Render(fmt.Sprintf("Historial reciente [%d]", len(s.History)))}
	if len(s.History) == 0 {
		rows = append(rows, MutedStyle.Render("No hay peticiones recientes."))
		return PanelStyle.Render(strings.Join(rows, "\n"))
	}

	const maxRows = 15
	start := len(s.History) - maxRows
	if start < 0 {
		start = 0
	}
	for i := len(s.History) - 1; i >= start; i-- {
		item := s.History[i]
		when := time.Unix(0, item.Time).Format("15:04:05")
		rows = append(rows, fmt.Sprintf("%s  %-6s  %s", when, item.Method, truncateURL(item.Url, 80)))
	}
	if start > 0 {
		rows = append(rows, MutedStyle.Render(fmt.Sprintf("... %d peticiones anteriores", start)))
	}
	return PanelStyle.Render(strings.Join(rows, "\n"))
}

func truncateURL(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit-3] + "..."
}
