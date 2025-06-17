package custom_error

import (
	"net/http"
)

const (
	DataNotFound = "DATA_NOT_FOUND"
)

var mapError = map[string]string{
	DataNotFound: "data not found"}

var mapHttpStatus = map[string]int{
	DataNotFound: http.StatusBadRequest,
}
