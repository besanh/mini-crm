package response

import (
	"net/http"
)

const (
	ERR_TOKEN_IS_EMPTY            = "token is empty"
	ERR_TOKEN_IS_INVALID          = "token is invalid"
	ERR_TOKEN_IS_EXPIRED          = "token is expired"
	ERR_EMPTY_CONN                = "empty connection"
	ERR_DATA_NOT_FOUND            = "data not found"
	ERR_DATA_INVALID              = "data is invalid"
	ERR_INSERT_FAILED             = "insert failed"
	ERR_GET_FAILED                = "get failed"
	ERR_PUT_FAILED                = "put failed"
	ERR_PATCH_FAILED              = "patch failed"
	ERR_DELETE_FAILED             = "delete failed"
	ERR_VALIDATION_FAILED         = "validation failed"
	ERR_INVALID_USERNAME_PASSWORD = "invalid username or password"

	ERR_REQUEST_INVALID         = "request invalid"
	ERR_REQUEST_IS_EXISTED      = "request is existed"
	ERR_REQUEST_NOTFOUND        = "request not found"
	ERR_REQUEST_IS_EXPIRED      = "request is expired"
	ERR_REQUEST_CODE_IS_INVALID = "request code is invalid"

	ERR_PERMISSION_DENIED = "permission denied"

	// Success
	SUCCESS = "success"
)

var MAP_ERR_RESPONSE = map[string]struct {
	Code    string
	Message string
}{
	ERR_TOKEN_IS_EMPTY: {
		Code:    "ERR_UNAUTHORIZE",
		Message: ERR_TOKEN_IS_EMPTY,
	},
	ERR_TOKEN_IS_INVALID: {
		Code:    "ERR_UNAUTHORIZE",
		Message: ERR_TOKEN_IS_INVALID,
	},
	ERR_TOKEN_IS_EXPIRED: {
		Code:    "ERR_UNAUTHORIZE",
		Message: ERR_TOKEN_IS_EXPIRED,
	},
	ERR_EMPTY_CONN: {
		Code:    "ERR_EMPTY_CONN",
		Message: ERR_EMPTY_CONN,
	},
	ERR_DATA_NOT_FOUND: {
		Code:    "ERR_DATA_NOT_FOUND",
		Message: ERR_DATA_NOT_FOUND,
	},
	ERR_DATA_INVALID: {
		Code:    "ERR_DATA_INVALID",
		Message: ERR_DATA_INVALID,
	},
	ERR_INSERT_FAILED: {
		Code:    "ERR_INSERT_FAILED",
		Message: ERR_INSERT_FAILED,
	},
	ERR_GET_FAILED: {
		Code:    "ERR_GET_FAILED",
		Message: ERR_GET_FAILED,
	},
	ERR_PUT_FAILED: {
		Code:    "ERR_PUT_FAILED",
		Message: ERR_PUT_FAILED,
	},
	ERR_PATCH_FAILED: {
		Code:    "ERR_PATCH_FAILED",
		Message: ERR_PATCH_FAILED,
	},
	ERR_DELETE_FAILED: {
		Code:    "ERR_DELETE_FAILED",
		Message: ERR_DELETE_FAILED,
	},
}

type (
	ErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	ValidatorFieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func OKResponse() (int, any) {
	return http.StatusOK, map[string]any{
		"message": "SUCCESS",
		"code":    http.StatusText(http.StatusOK),
	}
}

func BadRequest() (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusText(http.StatusBadRequest),
		"message": http.StatusText(http.StatusBadRequest),
	}
}

func BadRequestMsg(msg any) (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusText(http.StatusBadRequest),
		"message": msg,
	}
}

func NotFound() (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusText(http.StatusNotFound),
		"message": http.StatusText(http.StatusNotFound),
	}
}

func NotFoundMsg(msg any) (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusText(http.StatusNotFound),
		"message": msg,
	}
}

func Forbidden() (int, any) {
	return http.StatusForbidden, map[string]any{
		"error":   "Do not have permission for the request.",
		"code":    http.StatusText(http.StatusForbidden),
		"message": http.StatusText(http.StatusForbidden),
	}
}

func Unauthorized() (int, any) {
	return http.StatusUnauthorized, map[string]any{
		"error":   http.StatusText(http.StatusUnauthorized),
		"code":    http.StatusText(http.StatusUnauthorized),
		"message": http.StatusText(http.StatusUnauthorized),
	}
}

func ServiceUnavailable() (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusText(http.StatusServiceUnavailable),
		"message": http.StatusText(http.StatusServiceUnavailable),
	}
}

func ServiceUnavailableMsg(msg any) (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusText(http.StatusServiceUnavailable),
		"message": msg,
	}
}

func ResponseXml(field, val string) (int, any) {
	return http.StatusOK, map[string]any{field: val}
}

func Created(data any) (int, any) {
	result := map[string]any{
		"code":    http.StatusCreated,
		"message": "SUCCESS",
		"data":    data,
	}

	return http.StatusCreated, result
}

func Pagination(data, total, limit, offset any) (int, any) {
	return http.StatusOK, map[string]any{
		"data":   data,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
}

func OK(data any) (int, any) {
	return http.StatusOK, data
}
