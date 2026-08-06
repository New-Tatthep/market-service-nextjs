package handler

import (
	"market-service/custom_error"
	"market-service/datastore"
	"market-service/request"
	"market-service/service"
	"net/http"
	"strconv"

	"github.com/New-Tatthep/microservice"
)

const (
	tagLoginHandler  = "LoginHandler"
	tagLogoutHandler = "LogoutHandler"
)

func LoginHandler(ctx microservice.IContext) error {
	tag := tagLoginHandler

	req := &request.LoginRequest{}
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	service, err := service.New(
		service.WithContext(ctx),
		service.WithDatastore(ctx, datastore.PostgresContextName),
		service.WithRedisCache(ctx, service.RedisContextName),
		service.WithKafkaProducer(ctx, service.KafkaContextName),
	)
	if err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusInternalServerError,
			"new service failed",
			custom_error.Wrap(err),
		)
	}

	var resp *microservice.Field
	resp, err = service.CloudSession().CloudLogin(*req)
	if err != nil {
		if resp == nil {
			return custom_error.ErrorResponse(
				ctx,
				tag,
				http.StatusInternalServerError,
				"failed",
				custom_error.Wrap(err),
			)
		} else {
			return custom_error.ErrorResponse(
				ctx,
				tag,
				http.StatusInternalServerError,
				"failed",
				custom_error.Wrap(err),
				*resp,
			)
		}
	} else {
		return ctx.Response(
			microservice.DebugLevel,
			tag,
			http.StatusOK,
			strconv.Itoa(http.StatusOK),
			"success",
			nil,
			*resp,
		)
	}
}

func LogoutHandler(ctx microservice.IContext) error {
	tag := tagLogoutHandler

	req := &request.LogoutRequest{}
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	service, err := service.New(
		service.WithContext(ctx),
		service.WithDatastore(ctx, datastore.PostgresContextName),
		service.WithRedisCache(ctx, service.RedisContextName),
		service.WithKafkaProducer(ctx, service.KafkaContextName),
	)
	if err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusInternalServerError,
			"new service failed",
			custom_error.Wrap(err),
		)
	}

	err = service.CloudSession().CloudLogout(*req)
	if err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusInternalServerError,
			"failed",
			custom_error.Wrap(err),
		)
	} else {
		return ctx.Response(
			microservice.DebugLevel,
			tag,
			http.StatusOK,
			strconv.Itoa(http.StatusOK),
			"success",
			nil,
		)
	}
}
