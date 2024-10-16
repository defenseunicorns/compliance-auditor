package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/davecgh/go-spew/spew"
	"github.com/evertras/bubble-table/table"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
)

const (
	TabOffset     = 10
	DefaultWidth  = 200
	DefaultHeight = 60
)

var DumpFile *os.File

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
		return []key.Binding{ListKeys.Confirm, ListKeys.Help}
	}

	return d
}

func NewUnfocusedHighlightDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{ListKeys.Confirm, ListKeys.Help}
	}

	return d
}

func NewFocusedDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{ListKeys.Confirm, ListKeys.Help}
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

	return km
}

func UnfocusedPanelKeyMap() viewport.KeyMap {
	km := viewport.KeyMap{}

	return km
}

func FocusedTableKeyMap() table.KeyMap {
	km := table.DefaultKeyMap()
	km.PageUp = key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	)
	km.PageDown = key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdown", "page down"),
	)

	return km
}

func UnfocusedTableKeyMap() table.KeyMap {
	km := table.KeyMap{}

	return km
}

func FocusedTextAreaKeyMap() textarea.KeyMap {
	km := textarea.DefaultKeyMap

	km.InsertNewline = key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "insert newline"),
	)

	return km
}

func UnfocusedTextAreaKeyMap() textarea.KeyMap {
	km := textarea.KeyMap{}

	return km
}

func PrintToLog(format string, a ...any) {
	if DumpFile != nil {
		out := fmt.Sprintf(format, a...)
		spew.Fprintln(DumpFile, out)
	}
}

func DumpToLog(msg ...any) {
	if DumpFile != nil {
		spew.Fdump(DumpFile, msg)
	}
}

func ToYamlString(input interface{}) (string, error) {
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	return string(yamlData), nil
}

func DeepCopy(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
