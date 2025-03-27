package v1

import (
	"context"
	"net/http"

	"github.com/besanh/mini-crm/common/response"
	"github.com/besanh/mini-crm/services"
	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	api         hureg.APIGen
	userService services.IUsers
}

func NewUsers(api *hureg.APIGen, userService services.IUsers) {
	handler := &UserHandler{
		api:         *api,
		userService: userService,
	}

	group := api.AddBasePath("aaa/v1/users")
	{
		hureg.Register(group, huma.Operation{
			Method:      http.MethodGet,
			Path:        "{id}",
			Tags:        []string{"users"},
			OperationID: "GetUserByID",
		}, handler.GetUserByID)
	}
}

func (handler *UserHandler) GetUserByID(c context.Context, req *struct {
	Id uuid.UUID `path:"id" required:"true"`
}) (res *response.GenericResponse[any], err error) {
	result, err := handler.userService.GetUserByID(c, req.Id)
	if err != nil {
		err = response.HandleError(err)
		return
	}

	res = response.OKAny(result)
	return
}
