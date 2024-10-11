package assessmentresults

import (
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type focus int

const (
	noFocus focus = iota
	focusResultSelection
	focusCompareSelection
	focusFindings
	focusObservations
)

var maxFocus = focusObservations

func (m *Model) updateKeyBindings() {
	m.outOfFocus()
	m.updateFocusHelpKeys()

	switch m.focus {
	case focusFindings:
		m.findingsTable = m.findingsTable.WithKeyMap(common.FocusedTableKeyMap())
		m.findingsTable = m.findingsTable.Focused(true)
	case focusObservations:
		m.observationsTable = m.observationsTable.WithKeyMap(common.FocusedTableKeyMap())
		m.observationsTable = m.observationsTable.Focused(true)
	}
}

func (m *Model) outOfFocus() {
	focusMinusOne := m.focus - 1
	focusPlusOne := m.focus + 1

	if m.focus == 0 {
		focusMinusOne = maxFocus
	}
	if m.focus == maxFocus {
		focusPlusOne = 0
	}

	for _, f := range []focus{focusMinusOne, focusPlusOne} {
		switch f {
		case focusFindings:
			m.findingsTable = m.findingsTable.WithKeyMap(common.UnfocusedTableKeyMap())
			m.findingsTable = m.findingsTable.Focused(false)
		case focusObservations:
			m.observationsTable = m.observationsTable.WithKeyMap(common.UnfocusedTableKeyMap())
			m.observationsTable = m.observationsTable.Focused(false)
		}
	}
}

func (m *Model) updateFocusHelpKeys() {
	switch m.focus {
	case focusFindings:
		m.help.ShortHelp = common.ShortHelpTable
		m.help.FullHelpOneLine = common.FullHelpTableOneLine
		m.help.FullHelp = common.FullHelpTable
	case focusObservations:
		m.help.ShortHelp = common.ShortHelpTable
		m.help.FullHelpOneLine = common.FullHelpTableOneLine
		m.help.FullHelp = common.FullHelpTable
	default:
		m.help.ShortHelp = shortHelpNoFocus
		m.help.FullHelpOneLine = fullHelpNoFocusOneLine
		m.help.FullHelp = fullHelpNoFocus
	}
}
