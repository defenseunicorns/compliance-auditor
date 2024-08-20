package common

import (
	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/mattn/go-runewidth"
)

const TabOffset = 10

func TruncateText(text string, width int) string {
	if runewidth.StringWidth(text) <= width {
		return text
	}

	ellipsis := "…"
	trimmedWidth := width - runewidth.StringWidth(ellipsis)
	trimmedText := runewidth.Truncate(text, trimmedWidth, "")

	return trimmedText + ellipsis
}

func NewUnfocusedDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.NormalTitle
	d.Styles.SelectedDesc = d.Styles.NormalDesc

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{ListHotkeys.Confirm, ListHotkeys.Help}
	}

	return d
}

func NewUnfocusedHighlightDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{ListHotkeys.Confirm, ListHotkeys.Help}
	}

	return d
}

func NewFocusedDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{ListHotkeys.Confirm, ListHotkeys.Help}
	}

	return d
}

func FocusedListKeyMap() blist.KeyMap {
	km := blist.DefaultKeyMap()
	km.NextPage.Unbind()
	km.PrevPage.Unbind()
	km.ForceQuit.Unbind()
	km.Quit.Unbind()

	return km
}

func UnfocusedListKeyMap() blist.KeyMap {
	km := blist.KeyMap{}

	return km
}

func FocusedPanelKeyMap() viewport.KeyMap {
	km := viewport.DefaultKeyMap()
	// km.Up.SetEnabled(true)
	// km.Down.SetEnabled(true)

	return km
}

func UnfocusedPanelKeyMap() viewport.KeyMap {
	km := viewport.KeyMap{}

	return km
}
