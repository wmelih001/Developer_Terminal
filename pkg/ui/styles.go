package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPurple = lipgloss.Color("#bd93f9")
	ColorCyan   = lipgloss.Color("#8be9fd")
	ColorGrey   = lipgloss.Color("#44475a")
	ColorDark   = lipgloss.Color("#282a36")
	ColorRed    = lipgloss.Color("#ff5555")
	ColorGreen  = lipgloss.Color("#50fa7b")
	ColorYellow = lipgloss.Color("#f1fa8c")
	ColorWhite  = lipgloss.Color("#f8f8f2")
	ColorBlack  = lipgloss.Color("#282a36")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true).
			Padding(1, 2)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("252"))

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorPurple).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderLeftForeground(ColorPurple)

	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple)
)
