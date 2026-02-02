package handler

import (
	"market-service/service"
	"net/http"
	"strconv"

	"github.com/New-Tatthep/microservice"
)

const (
	tagProductFilterTransactionHandler = "ProductFilterTransactionHandler"
	tagProductGetTransactionHandler    = "ProductGetTransactionHandler"
)

func ProductFilterHandler(ctx microservice.IContext) error {
	tag := tagProductFilterTransactionHandler

	request := &microservice.FilterRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.Response(
			microservice.ErrorLevel,
			tag,
			http.StatusBadRequest,
			strconv.Itoa(http.StatusBadRequest),
			"invalid request",
			err,
		)
	}

	service, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
	)
	if err != nil {
		return ctx.Response(
			microservice.ErrorLevel,
			tag,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			"new service failed",
			err,
		)
	}

	resp, total, err := service.ProductAction().FilterProduct(*request)
	if err != nil {
		return ctx.Response(
			microservice.ErrorLevel,
			tag,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			"get reserve money failed",
			err,
		)
	} else {
		return ctx.Response(
			microservice.DebugLevel,
			tag,
			http.StatusOK,
			strconv.Itoa(http.StatusOK),
			"get reserve money success",
			nil,
			microservice.Field{
				Key:   "datas",
				Value: resp,
			},
			microservice.Field{
				Key:   "total",
				Value: total,
			},
		)
	}
}

func ProductGetHandler(ctx microservice.IContext) error {
	tag := tagProductGetTransactionHandler

	request := &microservice.FilterRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.Response(
			microservice.ErrorLevel,
			tag,
			http.StatusBadRequest,
			strconv.Itoa(http.StatusBadRequest),
			"invalid request",
			err,
		)
	}

	service, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
	)
	if err != nil {
		return ctx.Response(
			microservice.ErrorLevel,
			tag,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			"new service failed",
			err,
		)
	}

	resp, total, err := service.FilterProduct(*request)
	if err != nil {
		return ctx.Response(
			microservice.ErrorLevel,
			tag,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			"get reserve money failed",
			err,
		)
	} else {
		return ctx.Response(
			microservice.DebugLevel,
			tag,
			http.StatusOK,
			strconv.Itoa(http.StatusOK),
			"get reserve money success",
			nil,
			microservice.Field{
				Key:   "datas",
				Value: resp,
			},
			microservice.Field{
				Key:   "total",
				Value: total,
			},
		)
	}
}
