package common

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultPopupWidth  = 40
	defaultPopupHeight = 10
)

// Common Style definitions.
var (

	// Colors
	Text       = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}
	Subtle     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Subtle2    = lipgloss.AdaptiveColor{Light: "#706f6f", Dark: "#989797"}
	Highlight  = lipgloss.AdaptiveColor{Light: "#6d26fc", Dark: "#7D56F4"}
	Highlight2 = lipgloss.AdaptiveColor{Light: "#8f58fc", Dark: "#8f6ef0"}
	Focused    = lipgloss.AdaptiveColor{Light: "#8378ab", Dark: "#bfb2eb"}
	Special    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	Background = lipgloss.AdaptiveColor{Light: "#c5c6c7", Dark: "#333436"}
	Warning    = lipgloss.AdaptiveColor{Light: "#FFA100", Dark: "#F9A431"}

	HelpKey        = lipgloss.AdaptiveColor{Light: "#909090", Dark: "#626262"}
	HelpDesc       = lipgloss.AdaptiveColor{Light: "#B2B2B2", Dark: "#4A4A4A"}
	HelpSep        = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	ActiveHelpKey  = Highlight2
	ActiveHelpDesc = Highlight
	ActiveHelpSep  = Highlight

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
			Padding(1, 1, 0, 1)

	OverlayWarnStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder(), true).
				BorderForeground(Warning).
				Padding(1, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(Highlight).
			Padding(1, 2).
			Margin(1).
			Width(30)

	SummaryTextStyle = lipgloss.NewStyle().
				Foreground(Text).Margin(0, 1)

	// Table Styles
	TableStyles = table.Styles{
		Header:   lipgloss.NewStyle().Foreground(Highlight).Bold(true),
		Cell:     lipgloss.NewStyle().Foreground(Subtle),
		Selected: lipgloss.NewStyle().Foreground(Highlight),
	}

	UnfocusedTableStyles = table.Styles{
		Header:   lipgloss.NewStyle().Foreground(Highlight).Bold(true),
		Cell:     lipgloss.NewStyle().Foreground(Subtle),
		Selected: lipgloss.NewStyle().Foreground(Subtle),
	}

	TableStyleBase = lipgloss.NewStyle().
			Foreground(Text).
			BorderForeground(Subtle2).
			Align(lipgloss.Left)

	TableStyleActive = lipgloss.NewStyle().
				Foreground(Text).
				BorderForeground(Highlight).
				Align(lipgloss.Left)

	// Help
	KeyStyle        = lipgloss.NewStyle().Foreground(HelpKey)
	DescStyle       = lipgloss.NewStyle().Foreground(HelpDesc)
	SepStyle        = lipgloss.NewStyle().Foreground(HelpSep)
	ActiveKeyStyle  = lipgloss.NewStyle().Foreground(ActiveHelpKey)
	ActiveDescStyle = lipgloss.NewStyle().Foreground(ActiveHelpDesc)
	ActiveSepStyle  = lipgloss.NewStyle().Foreground(ActiveHelpSep)
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
