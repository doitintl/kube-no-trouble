package printer

import (
	"fmt"
	"strings"
)

func mapToCommaSeparatedString(m map[string]interface{}) string {
	var sb strings.Builder
	for k, v := range m {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s:%v", k, v))
	}
	return sb.String()
}
