package template

import (
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
)

func concatToRegoList(list []any) string {
	regoList := ""
	for i, item := range list {
		if i > 0 {
			regoList += ", "
		}
		regoList += `"` + fmt.Sprintf("%v", item) + `"`
	}
	return regoList
}

func newUUID(source ...string) string {
	if len(source) == 0 {
		return uuid.NewUUID()
	} else {
		return uuid.NewUUIDWithSource(source[0])
	}
}

func timestamp() string {
	t := time.Now()
	return t.Format("2006-01-02T15:04:05.999999999Z")
}
