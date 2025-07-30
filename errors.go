package directus

import "strings"

type ErrorDetailsPart struct {
	Message string `json:"message"`
}

func (e ErrorDetailsPart) Error() string {
	return e.Message
}

type ErrorDetails []ErrorDetailsPart

func (e ErrorDetails) Error() string {
	switch {
	case len(e) == 0:
		return ""
	case len(e) == 1:
		return e[0].Error()
	default:
		builder := strings.Builder{}
		for _, err := range e {
			if builder.Len() > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(err.Error())
		}
		return builder.String()
	}
}

type Error struct {
	Status  int
	Details error
}

func (e Error) Error() string {
	return e.Details.Error()
}

type ErrorResponse struct {
	Errors ErrorDetails `json:"errors"`
}
