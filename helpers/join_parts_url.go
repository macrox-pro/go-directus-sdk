package helpers

import "strings"

const (
	urlPathSep = "/"
)

type PartURL struct {
	Value string
	Skip  bool
}

func JoinPartsURL(parts ...PartURL) string {
	if len(parts) == 0 {
		return ""
	}

	elems := make([]string, 0, len(parts))
	for _, part := range parts {
		if part.Skip {
			continue
		}

		elems = append(
			elems,
			strings.Trim(
				part.Value,
				urlPathSep,
			),
		)
	}

	return strings.Join(elems, urlPathSep)
}
