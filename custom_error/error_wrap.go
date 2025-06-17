package custom_error

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/New-Tatthep/microservice/util/stringutil"
)

type CustomError struct {
	Err      error  `json:"error"`
	Code     string `json:"code"`
	Function string `json:"function"`
	Line     int    `json:"line"`
}

func (err *CustomError) ToJson() string {
	return stringutil.Json(*err)
}

func (err *CustomError) ToMap() map[string]any {
	return map[string]any{
		"error":    err.Err.Error(),
		"code":     err.Code,
		"function": err.Function,
		"line":     err.Line,
	}
}

func (err *CustomError) String() string {
	return fmt.Sprintf("Function %s - Line %d - Code %s\n%s", err.Function, err.Line, err.Code, err.Err.Error())
}

func (err *CustomError) Error() string {
	return err.Err.Error()
}

func ToMap(err error) map[string]any {
	if customErr, ok := err.(*CustomError); ok {
		return map[string]any{
			"error":    customErr.Err.Error(),
			"code":     customErr.Code,
			"function": customErr.Function,
			"line":     customErr.Line,
		}
	} else {
		return map[string]any{
			"error": err.Error(),
		}
	}
}

// if var err match to mapError
// this func will set code and error from mapError
// else code will be "NOT_SPECIFIED" and use var err as error
func Wrap(err error) *CustomError {
	if customErr, ok := err.(*CustomError); ok {
		// return if already wrap
		return customErr
	}

	// get function name and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	resp := &CustomError{
		Function: frame.Function,
		Line:     frame.Line,
	}

	errMsg, ok := mapError[err.Error()]
	if ok {
		resp.Code = err.Error()
		resp.Err = errors.New(errMsg)
	} else {
		resp.Code = "NOT_SPECIFIED"
		resp.Err = err
	}

	return resp
}

func WrapMicroservice(err error, code string) *CustomError {
	if customErr, ok := err.(*CustomError); ok {
		// return if already wrap
		return customErr
	}

	// get function name and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	resp := &CustomError{
		Function: frame.Function,
		Line:     frame.Line,
	}

	_, ok := mapError[code]
	if ok {
		resp.Code = code
		resp.Err = err
	} else {
		resp.Code = "NOT_SPECIFIED"
		resp.Err = err
	}

	return resp
}

func New(errStr string) *CustomError {
	err := errors.New(errStr)

	// get function name and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	resp := &CustomError{
		Function: frame.Function,
		Line:     frame.Line,
	}

	errMsg, ok := mapError[err.Error()]
	if ok {
		resp.Code = err.Error()
		resp.Err = errors.New(errMsg)
	} else {
		resp.Code = "NOT_SPECIFIED"
		resp.Err = err
	}

	return resp
}

func Newf(format string, args ...any) *CustomError {
	err := fmt.Errorf(format, args...)

	// get function name and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	resp := &CustomError{
		Function: frame.Function,
		Line:     frame.Line,
	}

	errMsg, ok := mapError[err.Error()]
	if ok {
		resp.Code = err.Error()
		resp.Err = errors.New(errMsg)
	} else {
		resp.Code = "NOT_SPECIFIED"
		resp.Err = err
	}

	return resp
}
