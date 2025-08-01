package directus

import "strings"

type ErrorCode string

const (
	ForbiddenErrorCode            ErrorCode = "FORBIDDEN"
	InvalidIPErrorCode            ErrorCode = "INVALID_IP"
	InvalidOTPErrorCode           ErrorCode = "INVALID_OTP"
	InvalidQueryErrorCode         ErrorCode = "INVALID_QUERY"
	InvalidTokenErrorCode         ErrorCode = "INVALID_TOKEN"
	TokenExpiredErrorCode         ErrorCode = "TOKEN_EXPIRED"
	RouteNotFoundErrorCode        ErrorCode = "ROUTE_NOT_FOUND"
	InvalidPayloadErrorCode       ErrorCode = "INVALID_PAYLOAD"
	RequestExceededErrorCode      ErrorCode = "REQUESTS_EXCEEDED"
	FailedValidationErrorCode     ErrorCode = "FAILED_VALIDATION"
	ServiceUnavailableErrorCode   ErrorCode = "SERVICE_UNAVAILABLE"
	InvalidCredentialsErrorCode   ErrorCode = "INVALID_CREDENTIALS"
	UnprocessableContentErrorCode ErrorCode = "UNPROCESSABLE_CONTENT"
	UnsupportedMediaTypeErrorCode ErrorCode = "UNSUPPORTED_MEDIA_TYPE"
)

func (code ErrorCode) Status() int {
	switch code {
	case ForbiddenErrorCode,
		InvalidTokenErrorCode:
		return 403

	case FailedValidationErrorCode,
		InvalidPayloadErrorCode,
		InvalidQueryErrorCode:
		return 400

	case InvalidCredentialsErrorCode,
		TokenExpiredErrorCode,
		InvalidOTPErrorCode,
		InvalidIPErrorCode:
		return 401

	case RouteNotFoundErrorCode:
		return 404

	case UnsupportedMediaTypeErrorCode:
		return 415

	case RequestExceededErrorCode:
		return 429

	case ServiceUnavailableErrorCode:
		return 503

	case UnprocessableContentErrorCode:
		return 422

	default:
		return 500
	}
}

type ErrorExtensions struct {
	Code ErrorCode `json:"code,omitempty"`
}

type Error struct {
	Extensions ErrorExtensions `json:"extensions,omitempty"`
	Message    string          `json:"message"`
}

func (err Error) Status() int {
	return err.Extensions.Code.Status()
}

func (err Error) Error() string {
	return err.Message
}

type Errors []Error

func (errs Errors) Status() int {
	var maxStatus int
	for _, err := range errs {
		status := err.Extensions.Code.Status()
		if maxStatus > status {
			maxStatus = status
		}
	}

	if maxStatus == 0 {
		return 500
	}

	return maxStatus
}

func (errs Errors) Error() string {
	switch {
	case len(errs) == 0:
		return ""

	case len(errs) == 1:
		return errs[0].Message

	default:
		builder := strings.Builder{}
		for _, err := range errs {
			if builder.Len() > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(err.Message)
		}

		return builder.String()
	}
}

type ErrorsPayload struct {
	Errors Errors `json:"errors"`
}
