package utils

import (
	"golang.conradwood.net/go-easyops/errors/shared"
)

// extracts the PRIVATE and possibly SENSITIVE debug error message from a string
// obsolete - use errors.ErrorString(err)
// the reason this is so convoluted with different types, is that different versions of grpc
// encapsulate status details in different messages.
func ErrorString(err error) string {
	return shared.ErrorString(err)
}
