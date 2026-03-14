package mux

import "github.com/charmbracelet/lipgloss"

var (
	// Pane headers
	focusHeaderStyle   = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("15")).Padding(0, 1)
	sidebarHeaderStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("250")).Padding(0, 1)

	// Status dots (rendered strings, not styles)
	dotRunning = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("●")
	dotReady   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("●")
	dotExited  = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("●")
	dotError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("●")

	// Status bar (bottom of screen)
	statusBarStyle = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("15")).Padding(0, 1)

	// Focus border
	focusBorderStyle = lipgloss.NewStyle().BorderRight(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("57"))

	// Control pane
	controlHeaderStyle = lipgloss.NewStyle().Background(lipgloss.Color("17")).Foreground(lipgloss.Color("75")).Padding(0, 1)
	controlLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Width(12)
	controlValueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

	// Muted text
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
)
