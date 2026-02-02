package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"market-service/request"
	"market-service/response"
	"net/http"

	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/go-resty/resty/v2"
)

type LoginToCloudResponse struct {
	RequestId  string            `json:"request_id,omitempty"`
	StatusCode int               `json:"status_code"`
	Code       string            `json:"code,omitempty"`
	Message    string            `json:"message,omitempty"`
	Data       LoginToCloundData `json:"data,omitempty"`
	Error      Errors            `json:"errors,omitempty"`
}

type Errors []Error

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

type LoginToCloundData struct {
	Data response.LoginWithLocalResponse `json:"data"`
}

func (api *API) CloudLogin(body request.LoginWithLocalRequest) (*LoginToCloudResponse, error) {
	endpoint := api.customConfig.CloudLoginEndpoint

	if stringutil.IsEmptyString(endpoint) {
		return nil, errors.New("endpoint is required")
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(body).
		Post(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%d : %s", resp.StatusCode(), string(resp.Body()))
	}

	var data LoginToCloudResponse
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (api *API) CloudExtend(token string) error {
	endpoint := api.customConfig.CloudExtendEndpoint

	if stringutil.IsEmptyString(endpoint) {
		return errors.New("endpoint is required")
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader(http.CanonicalHeaderKey("x-auth-token"), token).
		Post(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%d : %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func (api *API) CloudReSetPassword(body request.ResetPasswordWithLocalRequest) (*LoginToCloudResponse, error) {
	endpoint := api.customConfig.CloudResetPasswordEndpoint

	if stringutil.IsEmptyString(endpoint) {
		return nil, errors.New("endpoint is required")
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(body).
		Post(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%d : %s", resp.StatusCode(), string(resp.Body()))
	}

	var data LoginToCloudResponse
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
