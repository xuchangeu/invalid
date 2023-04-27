package invalid

import "strings"

func deepFieldWithDot(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	return strings.Join(keys, ConstraintKeyDepthIndicator)
}

const (
	ConstraintKeyDepthIndicator = "."
)
