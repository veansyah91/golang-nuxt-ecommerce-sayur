package handler

import (
	"net/http"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UserHandlerInterface interface {
	SignIn(ctx echo.Context) error
}

type UserHandler struct{
	UserService service.UserServiceInterface
}

var err error

func (u *UserHandler) SignIn(c echo.Context) error {
	var (
		req = request.SignInRequest{}
		resp = response.DefaultReponse{}
		respSignIn = response.SignInResponse{}
		ctx = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil{
		log.Errorf("[UserHandler-1] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	

}

func NewUserHandler(UserService service.UserServiceInterface) UserHandlerInterface{
	return &UserHandler(UserService: UserService)
}
