package v1

import (
	"context"
	"net/http"

	"github.com/besanh/mini-crm/common/response"
	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
)

type HealthCheckHandler struct {
	api hureg.APIGen
}

func NewHealthCheck(api *hureg.APIGen) {
	handler := &HealthCheckHandler{api: *api}

	group := handler.api.AddBasePath("aaa/v1/health-check")
	{
		hureg.Register(group, huma.Operation{
			Method:   http.MethodGet,
			Path:     "",
			Tags:     []string{"health_check"},
			Security: nil,
		}, handler.PingCheckHealth)
	}
}

func (handler *HealthCheckHandler) PingCheckHealth(c context.Context, req *struct{}) (res *response.GenericResponse[any], err error) {
	res = response.OKAny("pong")
	return
}
