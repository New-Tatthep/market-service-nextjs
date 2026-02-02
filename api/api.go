package api

import (
	"errors"
	"fmt"
	"market-service/custom_config"

	"github.com/New-Tatthep/microservice/util/stringutil"
)

type API struct {
	customConfig *custom_config.CustomConfiguration
}

type Option func(api *API) error

func WithCustomConfig() Option {
	return func(api *API) error {
		customConfig, err := custom_config.Config()
		if err != nil {
			return err
		}

		api.customConfig = customConfig

		return nil
	}
}

func New(options ...Option) (*API, error) {
	api := &API{}

	for _, opt := range options {
		if err := opt(api); err != nil {
			return nil, err
		}
	}

	// auto load custom config
	customConfig, err := custom_config.Config()
	if err != nil {
		return nil, err
	}

	api.customConfig = customConfig

	return api, nil
}

func getApiEndpoint(category string, endpointName string) (string, error) {
	if stringutil.IsEmptyString(category) {
		return "", errors.New("category is required")
	}

	if stringutil.IsEmptyString(endpointName) {
		return "", errors.New("apiName is required")
	}

	customConfig, err := custom_config.Config()
	if err != nil {
		return "", err
	}

	apiData, ok := customConfig.Api[category]
	if !ok {
		return "", fmt.Errorf("api category %s not found", category)
	}

	endpoint, ok := apiData.Endpoint[endpointName]
	if !ok {
		return "", fmt.Errorf("api category %s endpoint %s not found", category, endpointName)
	}

	url := apiData.BaseUrl + endpoint

	return url, nil
}
