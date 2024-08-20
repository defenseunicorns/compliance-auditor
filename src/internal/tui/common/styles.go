package common

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30

	modalWidth  = 60
	modalHeight = 7
)

// Style definitions.
var (

	// Colors

	Subtle     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Highlight  = lipgloss.AdaptiveColor{Light: "#6d26fc", Dark: "#7D56F4"}
	Focused    = lipgloss.AdaptiveColor{Light: "#8378ab", Dark: "#bfb2eb"}
	Special    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	Background = lipgloss.AdaptiveColor{Light: "#c5c6c7", Dark: "#333436"}

	// Tabs

	ActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	TabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	Tab = lipgloss.NewStyle().
		Border(TabBorder, true).
		BorderForeground(Highlight).
		Padding(0, 1)

	ActiveTab = Tab.Border(ActiveTabBorder, true)

	TabGap = Tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// Widgets

	LabelStyle = lipgloss.NewStyle().
			Margin(1).
			Foreground(Highlight)

	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(Highlight)

	PanelStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder(), false, true, true, true).
			BorderForeground(Highlight)

	PanelTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).BorderForeground(Highlight).Height(1)
	}()

	HighlightedBorderStyle = lipgloss.NewStyle().Foreground(Highlight)

	OverlayStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(Focused).
			Padding(1, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(Highlight).
			Padding(1, 2).
			Margin(1).
			Width(30)
)

func HeaderView(titleText string, width int, focusColor lipgloss.AdaptiveColor) string {
	title := PanelTitleStyle.BorderForeground(focusColor).Render(titleText)
	lineWidth := max(0, width-lipgloss.Width(title)-2) // Adjust for spacing
	rightTitleBorder := lipgloss.NewStyle().Foreground(focusColor).Render("─" + strings.Repeat("─", lineWidth) + "╮")
	return lipgloss.JoinHorizontal(lipgloss.Center, title, rightTitleBorder)
}

func HelpStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(width - PanelStyle.GetHorizontalPadding() - PanelStyle.GetHorizontalMargins()).
		Height(1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
