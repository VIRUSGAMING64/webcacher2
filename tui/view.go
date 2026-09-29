package tui

import (
	"strings"
	"webcacher2/tui/views"
)

func renderModel(m model) string {
	var body string
	switch m.view {
	case CACHE:
		body = views.RenderCache()
	case SERVER:
		body = views.RenderServer()
	case QUEUE:
		body = views.RenderQueue()
	default:
		body = views.RenderDashboard()
	}

	header := views.TitleStyle.Render("WebCacher TUI") + "  " + views.MutedStyle.Render("r refrescar  q salir")
	navigation := strings.Join([]string{
		navItem("1 Dashboard", m.view == DASH),
		navItem("2 Cola", m.view == QUEUE),
		navItem("3 Historial", m.view == HIST),
		navItem("4 Cache", m.view == CACHE),
		navItem("5 Server", m.view == SERVER),
	}, "  ")
	return views.PanelStyle.Render(header + "\n" + navigation + "\n\n" + body + "\n")
}

func navItem(label string, selected bool) string {
	if selected {
		return views.SelectedStyle.Render("[" + label + "]")
	}
	return views.MutedStyle.Render("[" + label + "]")
}
