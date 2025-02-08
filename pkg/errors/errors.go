package errors

import (
	"errors"
	"strconv"
)

var (
	ErrNilContextReader = errors.New("nil context reader")
	ErrNilClient        = errors.New("nil client")
)

type DockerError struct {
	Message string
	Code    int
	Step    string
	Cause   error
	Raw     []byte
}

func (dockerError *DockerError) Error() string {
	return dockerError.Message
}

func (dockerError *DockerError) GetCode() string {
	if dockerError.Code == 0 {
		return ""
	}
	return strconv.Itoa(dockerError.Code)
}

func (dockerError *DockerError) Unwrap() error {
	return dockerError.Cause
}

func (dockerError *DockerError) GetCause() error {
	return dockerError.Cause
}
