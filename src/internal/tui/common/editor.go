package common

import "time"

type EditType int

const (
	EditTypeAdd EditType = iota
	EditTypeUpdate
	EditTypeDelete
)

type Editor struct {
	EditsByPath map[string][]*Edit
}

type Edit struct {
	EditType
	Path      string
	Value     map[string]interface{} // could this handle any type?
	Timestamp time.Time
}

func NewEditor() *Editor {
	return &Editor{
		EditsByPath: make(map[string][]*Edit),
	}
}

// ResetEditor resets the editor - runs after data is saved
func (e *Editor) ResetEditor() {
	e.EditsByPath = make(map[string][]*Edit)
}

// IsEmpty checks if any edits have been made
func (e *Editor) IsEmpty() bool {
	return len(e.EditsByPath) == 0
}

func (e *Editor) AddUpdateEdit(path string, value map[string]interface{}) {
	e.EditsByPath[path] = append(e.EditsByPath[path],
		&Edit{
			EditType:  EditTypeUpdate,
			Path:      path,
			Value:     value,
			Timestamp: time.Now(),
		})
}
