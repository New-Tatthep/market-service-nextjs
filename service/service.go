package service

import (
	"fmt"
	"market-service/custom_config"
	"market-service/custom_error"
	"market-service/datastore"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/log"
)

const (
	RedisContextName = "rediscache"
	KafkaContextName = "mq"
)

type Action interface {
	Validate() error

	CloudSession() CloudSessionAction
	ProductAction() ProductServiceAction
	EmployeeAction() EmployeeServiceAction
}

type service struct {
	ctx            microservice.IContext
	store          datastore.StoreAction
	customConfig   *custom_config.CustomConfiguration
	cache          microservice.IRedisCache
	authConfig     *AuthConfig
	producer       microservice.IKafkaProducer
	logger         microservice.IContextLogger
	adConfig       map[string]ActiveDirectorySetting
	passwordPolicy *PasswordPolicySetting
}

type Option func(sv *service) error

func WithContext(ctx microservice.IContext) Option {
	return func(sv *service) error {
		sv.ctx = ctx

		return nil
	}
}

func WithDatastore(ctx microservice.IContext, databaseContext string) Option {
	return func(sv *service) error {
		store, err := datastore.New(
			datastore.WithDBContext(sv.ctx, databaseContext),
			datastore.WithMicroserviceLogger(sv.ctx.Logger()),
		)
		if err != nil {
			return err
		}

		sv.store = store

		return nil
	}
}

func WithCustomConfig() Option {
	return func(sv *service) error {
		customConfig, err := custom_config.Config()
		if err != nil {
			return err
		}

		sv.customConfig = customConfig

		return nil
	}
}

func WithRedisCache(ctx microservice.IContext, redisContext string) Option {
	return func(sv *service) error {
		cache, found := ctx.Cache(redisContext)
		if !found {
			return fmt.Errorf("cache %s not found", redisContext)
		}

		sv.cache = cache

		return nil
	}
}

func WithKafkaProducer(ctx microservice.IContext, kafkaContext string) Option {
	return func(sv *service) error {
		producer, found := ctx.Producer(kafkaContext)
		if !found {
			return fmt.Errorf("producer %s not found", kafkaContext)
		}

		sv.producer = producer

		return nil
	}
}

func WithMicroserviceDatastore(store microservice.IDBStore) Option {
	return func(sv *service) error {
		store, err := datastore.New(
			datastore.WithMicroserviceDB(store),
		)
		if err != nil {
			return err
		}

		sv.store = store

		return nil
	}
}

func WithMicroserviceRedisCache(cache microservice.IRedisCache) Option {
	return func(sv *service) error {
		sv.cache = cache

		return nil
	}
}

func WithMicroserviceLogger(logger microservice.IContextLogger) Option {
	return func(sv *service) error {
		sv.logger = logger

		return nil
	}
}

func (sv *service) Validate() error {
	// if sv.ctx == nil {
	// 	return errors.New("not found context")
	// }

	// if sv.store == nil {
	// 	return errors.New("not found datastore")
	// }

	// if sv.cache == nil {
	// 	return errors.New("not found redis cache")
	// }

	// if sv.producer == nil {
	// 	return errors.New("not found kafka producer")
	// }

	return nil
}

func New(options ...Option) (*service, error) {
	sv := new(service)

	for _, opt := range options {
		if err := opt(sv); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	if err := sv.Validate(); err != nil {
		return nil, custom_error.Wrap(err)
	}

	// auto load custom config
	customConfig, err := custom_config.Config()
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	sv.customConfig = customConfig

	// auto load config from redis
	// if !customConfig.IsSkipReloadCache {
	authConfig, foundAuthConfig, err := sv.GetAuthConfig()
	if err != nil {
		log.Errorf("GetAuthConfig Erorr : %s", err.Error())
		// return nil, err
	}

	if foundAuthConfig {
		sv.authConfig = authConfig
	}

	adConfig, foundAdConfig, err := sv.GetADConfig()
	if err != nil {
		log.Errorf("GetADConfig Erorr : %s", err.Error())
		// return nil, err
	}

	if foundAdConfig {
		sv.adConfig = adConfig
	}

	passwordPolicyConfig, foundPasswordPolicyConfig, err := sv.GetPasswordPolicy()
	if err != nil {
		log.Errorf("GetPasswordPolicy Erorr : %s", err.Error())
		// return nil, err
	}

	if foundPasswordPolicyConfig {
		sv.passwordPolicy = passwordPolicyConfig
	}

	// }

	if sv.ctx != nil {
		if sv.cache == nil {
			cache, found := sv.ctx.Cache(RedisContextName)
			if !found {
				return nil, fmt.Errorf("cache %s not found", RedisContextName)
			}

			sv.cache = cache
		}

		sv.logger = sv.ctx.Logger()
	}

	return sv, nil
}

func WithLogger(ctx microservice.IContext) Option {
	return func(sv *service) error {
		sv.logger = ctx.Logger()

		return nil
	}
}

func NewServiceFullOptionWithContext(ctx microservice.IContext) []Option {
	var opt []Option

	opt = append(opt,
		WithContext(ctx),
		WithDatastore(ctx, datastore.PostgresContextName),
		WithRedisCache(ctx, RedisContextName),
		WithKafkaProducer(ctx, KafkaContextName),
		WithLogger(ctx),
	)

	return opt
}
