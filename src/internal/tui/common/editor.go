package common

import "time"

type Edit struct {
	timestamp time.Time
	value     any
}

func NewEdit(value any) *Edit {
	return &Edit{
		timestamp: time.Now(),
		value:     value,
	}
}
