package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Run inicia la interfaz de monitorización y devuelve cuando el usuario sale.
func Run() error {
	program := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("iniciar TUI: %w", err)
	}
	return nil
}
