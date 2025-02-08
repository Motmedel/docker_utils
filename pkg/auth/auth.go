package auth

import (
	"encoding/base64"
	"encoding/json"
	dockerUtilsErrors "github.com/Motmedel/docker_utils/pkg/errors"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	dockerTypesRegistry "github.com/docker/docker/api/types/registry"
	"golang.org/x/oauth2"
)

func GetAuthString(registryHost string, token *oauth2.Token) (string, error) {
	if registryHost == "" {
		return "", dockerUtilsErrors.ErrEmptyRegistryHost
	}

	if token == nil {
		return "", dockerUtilsErrors.ErrNilToken
	}

	authConfigJsonString, err := json.Marshal(
		dockerTypesRegistry.AuthConfig{
			Username:      "oauth2accesstoken",
			Password:      token.AccessToken,
			ServerAddress: registryHost,
		},
	)
	if err != nil {
		return "", &motmedelErrors.CauseError{
			Message: "An error occurred when marshalling the auth config.",
			Cause:   err,
		}
	}

	return base64.StdEncoding.EncodeToString(authConfigJsonString), nil
}
