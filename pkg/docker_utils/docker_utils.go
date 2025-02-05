package docker_utils

import (
	"bufio"
	"encoding/json"
	dockerUtilsErrors "github.com/Motmedel/docker_utils/pkg/errors"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	"github.com/docker/docker/pkg/jsonmessage"
	"io"
)

func ScanOutput(outputReader io.Reader, callback func([]byte, *jsonmessage.JSONMessage)) error {
	if outputReader == nil {
		return nil
	}

	var rawLine []byte
	var dockerError *jsonmessage.JSONError

	scanner := bufio.NewScanner(outputReader)
	for scanner.Scan() {
		rawLine = scanner.Bytes()

		var message jsonmessage.JSONMessage
		if err := json.Unmarshal(rawLine, &message); err != nil {
			return &motmedelErrors.InputError{
				Message: "An error occurred when unmarshalling Docker output.",
				Cause:   err,
				Input:   rawLine,
			}
		}

		if messageError := message.Error; messageError != nil {
			dockerError = messageError
			break
		}

		if callback != nil {
			callback(rawLine, &message)
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return &motmedelErrors.CauseError{Message: "An error occurred when scanning the Docker output.", Cause: err}
	}

	if dockerError != nil {
		return &dockerUtilsErrors.DockerError{
			Message: dockerError.Message,
			Code:    dockerError.Code,
			Raw:     rawLine,
		}
	}

	return nil
}
