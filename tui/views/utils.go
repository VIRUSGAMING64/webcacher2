package views

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	MutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ValueStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	PanelStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	SelectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
)

func formatBytes(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%d B", value)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	amount := float64(value)
	for _, unit := range units {
		amount /= 1024
		if amount < 1024 {
			return fmt.Sprintf("%.1f %s", amount, unit)
		}
	}
	return fmt.Sprintf("%.1f PB", amount/1024)
}

func snapshotTime(value time.Time) string {
	if value.IsZero() {
		return "nunca"
	}
	return value.Local().Format("15:04:05")
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max-3] + "..."
}
