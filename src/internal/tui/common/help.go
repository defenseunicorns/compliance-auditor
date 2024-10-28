package common

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Note: This package has been adapted from the original help package in github.com/charmbracelet/bubbles

type HelpType int

const (
	HelpTypeMain HelpType = iota
	HelpTypeEdit
	HelpTypeList
	HelpTypeSelect
	HelpTypePanel
	HelpTypeQuit
)

// Styles is a set of available style definitions for the Help bubble.
type HelpStyles struct {
	Ellipsis lipgloss.Style

	// Styling for the short help
	ShortKey       lipgloss.Style
	ShortDesc      lipgloss.Style
	ShortSeparator lipgloss.Style

	// Styling for the full help
	FullKey       lipgloss.Style
	FullDesc      lipgloss.Style
	FullSeparator lipgloss.Style
}

func NewStyles(active bool) HelpStyles {
	if active {
		return HelpStyles{
			ShortKey:       ActiveKeyStyle,
			ShortDesc:      ActiveDescStyle,
			ShortSeparator: ActiveSepStyle,
			Ellipsis:       ActiveSepStyle,
			FullKey:        ActiveKeyStyle,
			FullDesc:       ActiveDescStyle,
			FullSeparator:  ActiveSepStyle,
		}
	}
	return HelpStyles{
		ShortKey:       KeyStyle,
		ShortDesc:      DescStyle,
		ShortSeparator: SepStyle,
		Ellipsis:       SepStyle,
		FullKey:        KeyStyle,
		FullDesc:       DescStyle,
		FullSeparator:  SepStyle,
	}
}

// HelpModel contains the state of the help view.
type HelpModel struct {
	Width int

	ShowAll bool // if true, render the "full" help menu
	OneLine bool

	ShortHelp       []key.Binding
	FullHelpOneLine []key.Binding
	FullHelp        [][]key.Binding

	ShortSeparator string
	FullSeparator  string

	// The symbol we use in the short help when help items have been truncated
	// due to width. Periods of ellipsis by default.
	Ellipsis string

	Styles HelpStyles
}

// New creates a new help view with some useful defaults.
func NewHelpModel(active bool) HelpModel {
	return HelpModel{
		ShortSeparator: " • ",
		FullSeparator:  "    ",
		Ellipsis:       "…",
		Styles:         NewStyles(active),
	}
}

// Update helps satisfy the Bubble Tea Model interface. It's a no-op.
func (m HelpModel) Update(_ tea.Msg) (HelpModel, tea.Cmd) {
	return m, nil
}

// View renders the help view's current state.
func (m HelpModel) View() string {
	if m.ShowAll {
		if m.OneLine {
			return m.SingleLineHelpView(m.FullHelpOneLine)
		}
		return m.FullHelpView(m.FullHelp)
	}
	return m.SingleLineHelpView(m.ShortHelp)
}

// SingleLineHelpView renders a single line help view from a slice of keybindings.
// If the line is longer than the maximum width it will be gracefully
// truncated, showing only as many help items as possible.
func (m HelpModel) SingleLineHelpView(bindings []key.Binding) string {
	if len(bindings) == 0 {
		return ""
	}

	var b strings.Builder
	var totalWidth int
	separator := m.Styles.ShortSeparator.Inline(true).Render(m.ShortSeparator)

	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		var sep string
		if totalWidth > 0 && i < len(bindings) {
			sep = separator
		}

		str := sep +
			m.Styles.ShortKey.Inline(true).Render(kb.Help().Key) + " " +
			m.Styles.ShortDesc.Inline(true).Render(kb.Help().Desc)

		w := lipgloss.Width(str)

		// If adding this help item would go over the available width, stop
		// drawing.
		if m.Width > 0 && totalWidth+w > m.Width {
			// Although if there's room for an ellipsis, print that.
			tail := " " + m.Styles.Ellipsis.Inline(true).Render(m.Ellipsis)
			tailWidth := lipgloss.Width(tail)

			if totalWidth+tailWidth < m.Width {
				b.WriteString(tail)
			}

			break
		}

		totalWidth += w
		b.WriteString(str)
	}

	return b.String()
}

// FullHelpView renders help columns from a slice of key binding slices. Each
// top level slice entry renders into a column.
func (m HelpModel) FullHelpView(groups [][]key.Binding) string {
	if len(groups) == 0 {
		return ""
	}

	// Linter note: at this time we don't think it's worth the additional
	// code complexity involved in preallocating this slice.
	//nolint:prealloc
	var (
		out []string

		totalWidth int
		sep        = m.Styles.FullSeparator.Render(m.FullSeparator)
		sepWidth   = lipgloss.Width(sep)
	)

	// Iterate over groups to build columns
	for i, group := range groups {
		if group == nil || !shouldRenderColumn(group) {
			continue
		}

		var (
			keys         []string
			descriptions []string
		)

		// Separate keys and descriptions into different slices
		for _, kb := range group {
			if !kb.Enabled() {
				continue
			}
			keys = append(keys, kb.Help().Key)
			descriptions = append(descriptions, kb.Help().Desc)
		}

		col := lipgloss.JoinHorizontal(lipgloss.Top,
			m.Styles.FullKey.Render(strings.Join(keys, "\n")),
			m.Styles.FullKey.Render(" "),
			m.Styles.FullDesc.Render(strings.Join(descriptions, "\n")),
		)

		// Column
		totalWidth += lipgloss.Width(col)
		if m.Width > 0 && totalWidth > m.Width {
			break
		}

		out = append(out, col)

		// Separator
		if i < len(group)-1 {
			totalWidth += sepWidth
			if m.Width > 0 && totalWidth > m.Width {
				break
			}
			out = append(out, sep)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}

func shouldRenderColumn(b []key.Binding) (ok bool) {
	for _, v := range b {
		if v.Enabled() {
			return true
		}
	}
	return false
}
