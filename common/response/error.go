package response

import (
	"fmt"
	"net/http"

	"github.com/besanh/mini-crm/common/constant"
	"github.com/danielgtaylor/huma/v2"
)

type CustomError struct {
	status   int
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	ErrorMsg string   `json:"error,omitempty"`
	Details  []string `json:"details,omitempty"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) GetStatus() int {
	return e.status
}

func HandleError(err error) huma.StatusError {
	errorCode := "ERR_SYSTEM_ERROR"
	errorCodeMessage := http.StatusText(http.StatusInternalServerError)
	errorMessage := err.Error()
	status := http.StatusInternalServerError

	if code, ok := constant.MAP_ERROR_CODE[constant.ERROR_CODE(err.Error())]; ok {
		status = http.StatusOK
		errorCode = err.Error()
		errorCodeMessage = code
		errorMessage = fmt.Sprintf("%s: %s", err.Error(), code)
	}

	return &CustomError{
		status:   status,
		Message:  errorCodeMessage,
		Code:     errorCode,
		ErrorMsg: errorMessage,
	}
}

func NewHumaError() {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		details := make([]string, len(errs))
		for i, err := range errs {
			details[i] = err.Error()
		}
		code := string(constant.ERR_REQUEST_INVALID)
		if message == string(constant.ERR_UNAUTHORIZED) {
			code = string(constant.ERR_UNAUTHORIZED)
			message = "User not authorized"
		}
		return &CustomError{
			status:  http.StatusOK,
			Code:    code,
			Message: message,
			Details: details,
		}
	}
}

func ErrUnauthorized(message ...string) huma.StatusError {
	msg := "User not authorized"
	if len(message) > 0 {
		msg = message[0]
	}
	return &CustomError{
		status:   http.StatusUnauthorized,
		Message:  msg,
		Code:     string(constant.ERR_UNAUTHORIZED),
		ErrorMsg: fmt.Sprintf("%s: %s", constant.ERR_UNAUTHORIZED, http.StatusText(http.StatusUnauthorized)),
	}
}

func ErrBadRequest(message ...string) huma.StatusError {
	msg := http.StatusText(http.StatusBadRequest)
	if len(message) > 0 {
		msg = message[0]
	}
	return &CustomError{
		status:   http.StatusBadRequest,
		Message:  msg,
		Code:     string(constant.ERR_REQUEST_INVALID),
		ErrorMsg: fmt.Sprintf("%s: %s", constant.ERR_REQUEST_INVALID, message),
	}
}

func ErrServiceUnavailable(message ...string) huma.StatusError {
	msg := http.StatusText(http.StatusServiceUnavailable)
	if len(message) > 0 {
		msg = message[0]
	}
	return &CustomError{
		status:   http.StatusServiceUnavailable,
		Message:  msg,
		Code:     string(constant.ERR_SERVICE_UNAVAILABLE),
		ErrorMsg: fmt.Sprintf("%s: %s", constant.ERR_SERVICE_UNAVAILABLE, message),
	}
}

func ErrInternalServerError(message ...string) huma.StatusError {
	msg := http.StatusText(http.StatusInternalServerError)
	if len(message) > 0 {
		msg = message[0]
	}
	return &CustomError{
		status:   http.StatusInternalServerError,
		Message:  msg,
		Code:     string(constant.ERR_INTERNAL_SERVER_ERROR),
		ErrorMsg: fmt.Sprintf("%s: %s", constant.ERR_INTERNAL_SERVER_ERROR, message),
	}
}
