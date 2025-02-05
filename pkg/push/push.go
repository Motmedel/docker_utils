package push

import (
	"context"
	"github.com/Motmedel/docker_utils/pkg/docker_utils"
	"github.com/Motmedel/docker_utils/pkg/errors"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	dockerTypesImage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

func PushWithClient(
	ctx context.Context,
	image string,
	dockerClient *client.Client,
	options dockerTypesImage.PushOptions,
	callback func([]byte, *jsonmessage.JSONMessage),
) error {
	if image == "" {
		return nil
	}

	if dockerClient == nil {
		return errors.ErrNilClient
	}

	imagePushResponse, err := dockerClient.ImagePush(ctx, image, options)
	if err != nil {
		return &motmedelErrors.CauseError{
			Message: "An error occurred when sending a request to the Docker daemon to push.",
			Cause:   err,
		}
	}
	defer imagePushResponse.Close()

	return docker_utils.ScanOutput(imagePushResponse, callback)
}

func Push(
	ctx context.Context,
	image string,
	options dockerTypesImage.PushOptions,
	callback func([]byte, *jsonmessage.JSONMessage),
) error {
	if image == "" {
		return nil
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &motmedelErrors.CauseError{
			Message: "An error occurred when creating the Docker client.",
			Cause:   err,
		}
	}

	return PushWithClient(ctx, image, dockerClient, options, callback)
}
