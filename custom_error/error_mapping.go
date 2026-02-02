package custom_error

import (
	"net/http"
)

const (
	DataNotFound            = "DATA_NOT_FOUND"
	SessionExpired          = "SESSION_EXPIRED"
	SessionNotFound         = "SESSION_NOT_FOUND"
	InvalidUsernamePassword = "INVALID_USERNAME_PASSWORD"
	TokenNotFound           = "TOKEN_NOT_FOUND"
	UserLocked              = "USER_LOCKED"
	UserInactive            = "USER_INACTIVE"
	InvalidRequest          = "INVALID_REQUEST"
	CounterNotRegister      = "COUNTER_NOT_REGISTER"
	CounterInactive         = "COUNTER_INACTIVE"
	CounterInvalid          = "COUNTER_INVALID"
	NoPermission            = "NO_PERMISSION"
	BranchDataNotFound      = "BRANCH_NOTFOUND"
	CounterIsLimit          = "COUNTER_LIMIT"
	AlreadyLogin            = "ALREADY_LOGIN"
	CurrentCompanyNotFound  = "CURRENT_COMPANY_NOT_FOUND"
)

var mapError = map[string]string{
	DataNotFound:            "data not found",
	SessionExpired:          "session expired",
	SessionNotFound:         "session not found",
	InvalidUsernamePassword: "invalid username or password",
	TokenNotFound:           "token not found",
	UserLocked:              "user has been locked",
	InvalidRequest:          "invalid request",
	CounterNotRegister:      "counter not register",
	CounterInactive:         "counter is inactive",
	CounterInvalid:          "counter data is invalid",
	NoPermission:            "user dont have permission to use this api",
	BranchDataNotFound:      "data not found",
	CounterIsLimit:          "counter limit",
	UserInactive:            "user inactive",
	AlreadyLogin:            "user is in use",
	CurrentCompanyNotFound:  "current company not found",
}

var mapHttpStatus = map[string]int{
	DataNotFound:            http.StatusBadRequest,
	SessionExpired:          http.StatusUnauthorized,
	SessionNotFound:         http.StatusUnauthorized,
	InvalidUsernamePassword: http.StatusUnauthorized,
	TokenNotFound:           http.StatusUnauthorized,
	UserLocked:              http.StatusUnprocessableEntity,
	InvalidRequest:          http.StatusBadRequest,
	CounterNotRegister:      http.StatusUnauthorized,
	CounterInactive:         http.StatusUnprocessableEntity,
	CounterInvalid:          http.StatusUnauthorized,
	NoPermission:            http.StatusUnprocessableEntity,
	BranchDataNotFound:      http.StatusUnauthorized,
	CounterIsLimit:          http.StatusUnauthorized,
	UserInactive:            http.StatusBadRequest,
	AlreadyLogin:            http.StatusBadRequest,
	CurrentCompanyNotFound:  http.StatusBadRequest,
}
