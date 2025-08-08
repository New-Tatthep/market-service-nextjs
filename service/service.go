package service

import (
	"errors"
	"fmt"
	"market-service/custom_config"
	"market-service/datastore"

	"github.com/New-Tatthep/microservice"
)

const (
	RedisContextName = "rediscache"
	KafkaContextName = "mq"
)

type service struct {
	ctx          microservice.IContext
	store        datastore.IAction
	dbStore      microservice.IDBStore
	customConfig *custom_config.CustomConfiguration
	cache        microservice.IRedisCache
	producer     microservice.IKafkaProducer
	logger       microservice.IContextLogger
}

type Option func(sv *service) error

// ===== Option =====

func WithDBStore(dbStore microservice.IDBStore) Option {
	return func(sv *service) error {
		sv.dbStore = dbStore
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

func WithContext(ctx microservice.IContext) Option {
	return func(sv *service) error {
		sv.ctx = ctx

		return nil
	}
}

func WithLogger(ctx microservice.IContext) Option {
	return func(sv *service) error {
		sv.logger = ctx.Logger()

		return nil
	}
}

func (sv *service) Validate() error {
	return nil
}

func New(options ...Option) (*service, error) {
	sv := new(service)

	for _, opt := range options {
		if err := opt(sv); err != nil {
			return nil, err
		}
	}

	if err := sv.Validate(); err != nil {
		return nil, err
	}

	if sv.ctx != nil {
		// if sv.cache == nil {
		// 	cache, found := sv.ctx.Cache(RedisContextName)
		// 	if !found {
		// 		return nil, fmt.Errorf("cache %s not found", RedisContextName)
		// 	}

		// 	sv.cache = cache
		// }
	}

	// auto load custom config
	customConfig, err := custom_config.Config()
	if err != nil {
		return nil, err
	}

	sv.customConfig = customConfig

	if sv.dbStore != nil {
		sv.store = datastore.Action(nil, datastore.WithDBStore(sv.dbStore))
	} else {
		if sv.ctx == nil {
			return nil, errors.New("context is nil")
		}
		sv.store = datastore.Action(sv.ctx)
	}

	if sv.ctx != nil {
		if sv.logger == nil {
			sv.logger = sv.ctx.Logger()
		}
	}

	return sv, nil
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

func NewServiceFullOptionWithContext(ctx microservice.IContext) []Option {
	var opt []Option

	opt = append(opt,
		WithContext(ctx),
		WithDatastore(ctx),
		// WithRedisCache(ctx, RedisContextName),
		// WithKafkaProducer(ctx, KafkaContextName),
		WithLogger(ctx),
		WithCustomConfig(),
	)

	return opt
}

func WithDatastore(ctx microservice.IContext) Option {
	return func(sv *service) error {
		sv.store = datastore.Action(ctx)
		return nil
	}
}
