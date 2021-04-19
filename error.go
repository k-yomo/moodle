package moodle

import "fmt"

type APIError struct {
	Reason           string  `json:"error"`
	ErrorCode        string  `json:"errorcode"`
	StackTrace       *string `json:"stacktrace,omitempty"`
	DebugInfo        *string `json:"debuginfo,omitempty"`
	ReproductionLink *string `json:"reproductionlink,omitempty"`
}

func (a *APIError) Error() string {
	return fmt.Sprintf(
		"Reason: %s, ErrorCode: %s, DebugInfo: %v, ReproductionLink: %v, StackTrace: %v",
		a.Reason,
		a.ErrorCode,
		a.DebugInfo,
		a.ReproductionLink,
		a.StackTrace,
	)
}

func Code(err error) string {
	if apiError, ok := err.(*APIError); ok {
		return apiError.ErrorCode
	}
	return "unknown"
}
