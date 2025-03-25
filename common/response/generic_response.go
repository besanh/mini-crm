package response

import "net/http"

type PaginationBodyResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Limit   int64  `json:"limit"`
	Offset  int64  `json:"offset"`
	Total   int64  `json:"total"`
}

type PaginationResponse[T any] struct {
	Body PaginationBodyResponse[T]
}

type GenericResponse[T any] struct {
	// Status int
	Body BodyResponse[T]
}

type BodyResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type MediaResponse struct {
	ContentType             string `header:"Content-Type"`
	ContentLength           int    `header:"Content-Length"`
	AcceptRanges            string `header:"Accept-Ranges"`
	AllowControlAllowOrigin string `header:"Access-Control-Allow-Origin"`
	CacheControl            string `header:"Cache-Control"`
	Body                    []byte
}

type IdResponse struct {
	Id string `json:"id"`
}

func OK[T any](data T, msgs ...string) (res *GenericResponse[T]) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	res = &GenericResponse[T]{
		// Status: http.StatusOK,
		Body: BodyResponse[T]{
			Code:    http.StatusOK,
			Message: msg,
			Data:    data,
		},
	}
	return
}

func OKOnly(msgs ...string) (res *GenericResponse[any]) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	res = &GenericResponse[any]{
		// Status: http.StatusOK,
		Body: BodyResponse[any]{
			Code:    http.StatusOK,
			Message: msg,
		},
	}
	return
}

func OKAny(data ...any) (res *GenericResponse[any]) {
	msg := "success"
	res = &GenericResponse[any]{
		// Status: http.StatusOK,
		Body: BodyResponse[any]{
			Code:    http.StatusOK,
			Message: msg,
		},
	}
	if len(data) > 0 {
		res.Body.Data = data[0]
	}
	return
}

func Pagination[T any](data T, total int64, limit int64, offset int64, msgs ...string) (res *PaginationResponse[T]) {
	msg := "success"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	res = &PaginationResponse[T]{
		// Status: http.StatusOK,
		Body: PaginationBodyResponse[T]{
			Code:    http.StatusOK,
			Message: msg,
			Data:    data,
			Total:   total,
			Limit:   limit,
			Offset:  offset,
		},
	}
	return
}

func BadRequest[T any]() (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
		},
	}
	return
}

func BadRequestWithMsg[T any](msg string) (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusBadRequest,
			Message: msg,
		},
	}
	return
}

func NotFound[T any]() (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		},
	}
	return
}

func NotFoundWithMsg[T any](msg string) (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusNotFound,
			Message: msg,
		},
	}
	return
}

func Forbidden[T any]() (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		},
	}
	return
}

func Unauthorized[T any]() (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		},
	}
	return
}

func ServiceUnavailable[T any]() (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusServiceUnavailable,
			Message: "Service Unavailable",
		},
	}
	return
}

func ServiceUnavailableWithMsg[T any](msg string) (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusServiceUnavailable,
			Message: msg,
		},
	}
	return
}

func ResponseXml[T any](field, value string) (res any) {
	res = map[string]any{field: value}
	return
}

func Created[T any]() (res *GenericResponse[T]) {
	res = &GenericResponse[T]{
		Body: BodyResponse[T]{
			Code:    http.StatusCreated,
			Message: "Created",
		},
	}
	return
}
