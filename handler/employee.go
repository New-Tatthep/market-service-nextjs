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

func FilterEmployeeHandler(ctx microservice.IContext) error {
	tag := "tagFilterEmployeeHandler"

	req := new(datastore.FilterData)
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	sv, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
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

	resp, err := sv.FilterEmployee(req)
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
			resp...,
		)
	}
}

func GetEmployeeHandler(ctx microservice.IContext) error {
	tag := "tagGetEmployeeHandler"

	req := new(datastore.FilterData)
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	sv, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
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

	resp, err := sv.FilterEmployee(req)
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
			resp...,
		)
	}
}

func CreateEmployeeHandler(ctx microservice.IContext) error {
	tag := "CreateEmployeeHandler"

	req := new(request.EmployeeRequest)
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	sv, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
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

	resp, err := sv.CreateEmployee(req)
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
			resp...,
		)
	}
}

func UpdateEmployeeHandler(ctx microservice.IContext) error {
	tag := "tagUpdateEmployeeHandler"

	req := new(request.EmployeeRequest)
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	sv, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
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

	resp, err := sv.UpdateEmployee(req)
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
			resp...,
		)
	}
}

func DeleteEmployeeHandler(ctx microservice.IContext) error {
	tag := "tagDeleteEmployeeHandler"

	req := new(request.EmployeeRequest)
	if err := ctx.Bind(req); err != nil {
		return custom_error.ErrorResponse(
			ctx,
			tag,
			http.StatusBadRequest,
			"invalid request",
			custom_error.Wrap(err),
		)
	}

	sv, err := service.New(
		service.NewServiceFullOptionWithContext(ctx)...,
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

	resp, err := sv.DeleteEmployee(req)
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
			resp...,
		)
	}
}
