package directus

import "strings"

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

type Errors []Error

func (errs Errors) Error() string {
	switch {
	case len(errs) == 0:
		return ""
	case len(errs) == 1:
		return errs[0].Error()
	default:
		builder := strings.Builder{}
		for _, err := range errs {
			if builder.Len() > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(err.Error())
		}
		return builder.String()
	}
}

type ErrorsResponse struct {
	Errors Errors `json:"errors"`
}
