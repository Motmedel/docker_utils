package errors

import "strconv"

type BuildError struct {
	Code    int
	Message string
	Raw     []byte
}

func (buildError *BuildError) Error() string {
	return buildError.Message
}

func (buildError *BuildError) GetCode() string {
	return strconv.Itoa(buildError.Code)
}
