package moodle

import "fmt"

type APIError struct {
	Err              *string `json:"error,omitempty"`
	Message          *string `json:"message,omitempty"`
	ErrorCode        string  `json:"errorcode"`
	StackTrace       *string `json:"stacktrace,omitempty"`
	Exception        *string `json:"exception,omitempty"`
	DebugInfo        *string `json:"debuginfo,omitempty"`
	ReproductionLink *string `json:"reproductionlink,omitempty"`
}

func (a *APIError) Error() string {
	return fmt.Sprintf(
		"Reason: %v, Message: %v, ErrorCode: %s, DebugInfo: %v, ReproductionLink: %v, StackTrace: %v",
		a.Err,
		a.Message,
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
