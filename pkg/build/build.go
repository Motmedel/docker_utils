package build

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	dockerUtilsErrors "github.com/Motmedel/docker_utils/pkg/errors"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelTarTypes "github.com/Motmedel/utils_go/pkg/tar/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"io"
)

const DockerignoreFilename = ".dockerignore"
const DockerfileFilename = "Dockerfile"

func GetDockerIgnorePatternsWithPath(archive motmedelTarTypes.Archive, path string) []string {
	if len(archive) == 0 {
		return nil
	}

	if tarEntry, ok := archive[path]; ok {
		dockerignoreContent := tarEntry.Content
		var lines []string
		scanner := bufio.NewScanner(bytes.NewReader(dockerignoreContent))

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		return lines
	}

	return nil
}

func GetDockerIgnorePatterns(archive motmedelTarTypes.Archive) []string {
	if len(archive) == 0 {
		return nil
	}
	return GetDockerIgnorePatternsWithPath(archive, DockerignoreFilename)
}

func Build(
	contextReader io.Reader,
	options *dockerTypes.ImageBuildOptions,
	callback func([]byte, *jsonmessage.JSONMessage),
) error {
	if options == nil {
		options = &dockerTypes.ImageBuildOptions{}
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &motmedelErrors.CauseError{
			Message: "An error occurred when creating a Docker client.",
			Cause:   err,
		}
	}

	imageBuildResponse, err := dockerClient.ImageBuild(context.Background(), contextReader, *options)
	if err != nil {
		return &motmedelErrors.CauseError{
			Message: "An error occurred when sending a request to the Docker daemon to build.",
			Cause:   err,
		}
	}
	defer imageBuildResponse.Body.Close()

	var rawLine []byte
	var buildError *jsonmessage.JSONError

	scanner := bufio.NewScanner(imageBuildResponse.Body)
	for scanner.Scan() {
		rawLine = scanner.Bytes()

		var message jsonmessage.JSONMessage
		if err := json.Unmarshal(rawLine, &message); err != nil {
			return &motmedelErrors.InputError{
				Message: "An error occurred when unmarshalling build output.",
				Cause:   err,
				Input:   rawLine,
			}
		}

		if messageError := message.Error; messageError != nil {
			buildError = messageError
			break
		}

		if callback != nil {
			callback(rawLine, &message)
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return &motmedelErrors.CauseError{Message: "An error occurred when scanning the build output.", Cause: err}
	}

	if buildError != nil {
		return &dockerUtilsErrors.BuildError{Message: buildError.Message, Code: buildError.Code, Raw: rawLine}
	}

	return nil
}
