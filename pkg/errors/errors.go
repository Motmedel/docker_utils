package errors

import (
	"strconv"
)

type BuildError struct {
	Message string
	Code    int
	Raw     []byte
	Cause   error
}

func (buildError *BuildError) Error() string {
	return buildError.Message
}

func (buildError *BuildError) GetCode() string {
	return strconv.Itoa(buildError.Code)
}

func (buildError *BuildError) Unwrap() error {
	return buildError.Cause
}

func (buildError *BuildError) GetCause() error {
	return buildError.Cause
}
