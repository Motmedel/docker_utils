package docker_utils

import (
	"bufio"
	"encoding/json"
	dockerUtilsErrors "github.com/Motmedel/docker_utils/pkg/errors"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	"github.com/docker/docker/pkg/jsonmessage"
	"io"
	"strings"
)

func ScanOutput(outputReader io.Reader, callback func([]byte, *jsonmessage.JSONMessage)) error {
	if outputReader == nil {
		return nil
	}

	var rawLine []byte
	var dockerError *jsonmessage.JSONError

	var penultimateMessage *jsonmessage.JSONMessage
	var lastMessage *jsonmessage.JSONMessage

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

		penultimateMessage = lastMessage
		lastMessage = &message

		if callback != nil {
			callback(rawLine, &message)
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return &motmedelErrors.CauseError{Message: "An error occurred when scanning the Docker output.", Cause: err}
	}

	if dockerError != nil {
		var cause *dockerUtilsErrors.DockerError

		if lastMessage != nil {
			lastLine := lastMessage.Stream
			if strings.HasPrefix(lastLine, "\u001b[91m") {
				cause = &dockerUtilsErrors.DockerError{
					Message: strings.TrimPrefix(
						strings.TrimSuffix(lastLine, "\u001b[0m"),
						"\u001b[91m",
					),
				}

				if penultimateMessage != nil {
					penultimateLine := penultimateMessage.Stream
					if strings.HasPrefix(penultimateLine, "Step ") {
						cause.Step = penultimateLine
					}
				}
			}
		}

		return &dockerUtilsErrors.DockerError{
			Message: dockerError.Message,
			Code:    dockerError.Code,
			Cause:   cause,
			Raw:     rawLine,
		}
	}

	return nil
}
