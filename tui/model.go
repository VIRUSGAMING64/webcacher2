package tui

import (
	"time"
	"webcacher2/queue"
	"webcacher2/tui/views"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

const (
	DASH int = iota
	QUEUE
	HIST
	CACHE
	SERVER
)

type model struct {
	view    int
	width   int
	height  int
	loading bool
}

func NewModel() model {
	return model{loading: true}
}

func (m model) Init() tea.Cmd {
	return tea.Batch()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "+":
			switch m.view {
			case DASH:
				queue.GQueue.Mtx.Lock()
				queue.GQueue.Workers += 1
				queue.GQueue.Mtx.Unlock()
			case QUEUE:
				views.Offset += 1
			}
		case "-":
			switch m.view {
			case DASH:
				queue.GQueue.Mtx.Lock()
				if queue.GQueue.Workers != 0 {
					queue.GQueue.Workers += -1
				}
				queue.GQueue.Mtx.Unlock()
			case QUEUE:
				views.Offset -= 1
				if views.Offset == 0 {
					views.Offset = 1
				}
			}
		case "r":
			m.loading = true
			return m, tea.Batch()
		case "1":
			m.view = DASH
		case "2":
			m.view = QUEUE
		case "3":
			m.view = HIST
		case "4":
			m.view = CACHE
		case "5":
			m.view = SERVER
		}
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
	case tickMsg:
		return m, tick()
	}
	return m, tick()
}

func (m model) View() string {
	return renderModel(m)
}
