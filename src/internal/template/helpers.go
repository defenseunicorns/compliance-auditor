package template

import "fmt"

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
