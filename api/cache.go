package api

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	CategoryCacheService                   = "cache-service"
	CacheServiceReloadAuthConfig           = "reload-auth-config"
	CacheServiceReloadADConfig             = "reload-ad-config"
	CacheServiceReloadPasswordPolicyConfig = "reload-password-policy-config"
)

func (api *API) CacheServiceReloadAuthConfig() error {
	url, err := getApiEndpoint(CategoryCacheService, CacheServiceReloadAuthConfig)
	if err != nil {
		return err
	}

	client := resty.New()

	resp, err := client.R().
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("call %s error status %d, message : %s", url, resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func (api *API) CacheServiceReloadADConfig() error {
	url, err := getApiEndpoint(CategoryCacheService, CacheServiceReloadADConfig)
	if err != nil {
		return err
	}

	client := resty.New()

	resp, err := client.R().
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("call %s error status %d, message : %s", url, resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func (api *API) CacheServiceReloadPasswordPolicyConfig() error {
	url, err := getApiEndpoint(CategoryCacheService, CacheServiceReloadPasswordPolicyConfig)
	if err != nil {
		return err
	}

	client := resty.New()

	resp, err := client.R().
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("call %s error status %d, message : %s", url, resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
