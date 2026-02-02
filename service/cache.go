package service

import (
	"encoding/json"
	"market-service/api"
	"market-service/custom_error"

	"github.com/New-Tatthep/microservice/util/stringutil"
)

const (
	RedisAuthConfigKey      = "auth_config"
	RedisADConfigKey        = "ad_config"
	PasswordPolicyConfigKey = "password_policy_config"
)

func (sv *service) GetAuthConfig() (*AuthConfig, bool, error) {
	redisKey := RedisAuthConfigKey
	found := false
	cacheJson, err := sv.cache.Get(redisKey)
	if err != nil {
		return nil, found, custom_error.Wrap(err)
	}

	var cache AuthConfig

	if stringutil.IsEmptyString(cacheJson) {
		// reload
		rest, err := api.New()
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		err = rest.CacheServiceReloadAuthConfig()
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		// after reload get from redis again
		cacheJson, err = sv.cache.Get(redisKey)
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		if stringutil.IsEmptyString(cacheJson) {
			return nil, found, custom_error.New("app_config not found in redis")
		}
	}

	if stringutil.IsNotEmptyString(cacheJson) {
		found = true
	}

	err = json.Unmarshal([]byte(cacheJson), &cache)
	if err != nil {
		return nil, found, custom_error.Wrap(err)
	}

	return &cache, found, nil
}

func (sv *service) GetADConfig() (map[string]ActiveDirectorySetting, bool, error) {
	redisKey := RedisADConfigKey

	found := false
	cacheJson, err := sv.cache.Get(redisKey)
	if err != nil {
		return nil, found, custom_error.Wrap(err)
	}

	cache := make(map[string]ActiveDirectorySetting)

	if stringutil.IsEmptyString(cacheJson) {
		// reload
		rest, err := api.New()
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		err = rest.CacheServiceReloadADConfig()
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		// after reload get from redis again
		cacheJson, err = sv.cache.Get(redisKey)
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		if stringutil.IsEmptyString(cacheJson) {
			return nil, found, custom_error.New("db_config not found in redis")
		}
	}

	if stringutil.IsNotEmptyString(cacheJson) {
		found = true
	}

	err = json.Unmarshal([]byte(cacheJson), &cache)
	if err != nil {
		return nil, found, custom_error.Wrap(err)
	}

	return cache, found, nil
}

func (sv *service) GetPasswordPolicy() (*PasswordPolicySetting, bool, error) {
	redisKey := PasswordPolicyConfigKey
	found := false

	cacheJson, err := sv.cache.Get(redisKey)
	if err != nil {
		return nil, found, custom_error.Wrap(err)
	}

	var cache PasswordPolicySetting

	if stringutil.IsEmptyString(cacheJson) {
		// reload
		rest, err := api.New()
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		err = rest.CacheServiceReloadPasswordPolicyConfig()
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		// after reload get from redis again
		cacheJson, err = sv.cache.Get(redisKey)
		if err != nil {
			return nil, found, custom_error.Wrap(err)
		}

		if stringutil.IsEmptyString(cacheJson) {
			return nil, found, custom_error.New("password_policy_config not found in redis")
		}
	}

	if stringutil.IsNotEmptyString(cacheJson) {
		found = true
	}

	err = json.Unmarshal([]byte(cacheJson), &cache)
	if err != nil {
		return nil, found, custom_error.Wrap(err)
	}

	return &cache, found, nil
}
