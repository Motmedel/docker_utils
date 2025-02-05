package build

import (
	"bufio"
	"bytes"
	"context"
	"github.com/Motmedel/docker_utils/pkg/docker_utils"
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

func BuildWithClient(
	dockerClient *client.Client,
	contextReader io.Reader,
	options *dockerTypes.ImageBuildOptions,
	callback func([]byte, *jsonmessage.JSONMessage),
) error {
	if dockerClient == nil {
		return dockerUtilsErrors.ErrNilClient
	}

	if contextReader == nil {
		return dockerUtilsErrors.ErrNilContextReader
	}

	if options == nil {
		options = &dockerTypes.ImageBuildOptions{}
	}

	imageBuildResponse, err := dockerClient.ImageBuild(context.Background(), contextReader, *options)
	if err != nil {
		return &motmedelErrors.CauseError{
			Message: "An error occurred when sending a request to the Docker daemon to build.",
			Cause:   err,
		}
	}
	defer imageBuildResponse.Body.Close()

	return docker_utils.ScanOutput(imageBuildResponse.Body, callback)
}

func Build(
	contextReader io.Reader,
	options *dockerTypes.ImageBuildOptions,
	callback func([]byte, *jsonmessage.JSONMessage),
) error {
	if contextReader == nil {
		return dockerUtilsErrors.ErrNilContextReader
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &motmedelErrors.CauseError{
			Message: "An error occurred when creating the Docker client.",
			Cause:   err,
		}
	}

	return BuildWithClient(dockerClient, contextReader, options, callback)
}
