package custom_error

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/log"
)

func ErrorResponse(ctx microservice.IContext, tag string, httpCode int, message string, err *CustomError, fields ...microservice.Field) error {
	// if an error has been handled (already mapping) return with map http
	if value, ok := mapHttpStatus[err.Code]; ok {
		httpCode = value
	}

	if value, ok := mapError[err.Code]; ok {
		message = value
	}

	// get function name and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	ctx.Logger().Log(microservice.DebugLevel, fmt.Sprintf("[%s] %s: %s	%v", frame.Function, strconv.Itoa(httpCode), message, err.ToJson()))

	return ctx.Response(
		microservice.ErrorLevel,
		tag,
		httpCode,
		err.Code,
		message,
		err.Err,
		fields...,
	)
}

func LogError(err *CustomError) {
	log.Errorf("%s", err.String())
}

func WarnError(err *CustomError) {
	log.Warnf("%s", err.String())
}
