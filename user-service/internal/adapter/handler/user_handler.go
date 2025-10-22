package handler

import (
	"net/http"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var err error

type UserHandlerInterface interface {
	SignIn(c echo.Context) error
	CreateUserAccount(c echo.Context) error
	ForgotPassword(c echo.Context) error
	VerifyAccount(c echo.Context) error
	UpdatePassword(c echo.Context) error
}

type UserHandler struct {
	UserService service.UserServiceInterface
}

// UpdatePassword implements UserHandlerInterface.
func (u *UserHandler) UpdatePassword(c echo.Context) error {
	var (
		resp = response.DefaultReponse{}
		req  = request.UpdatePasswordRequest{}
		ctx  = c.Request().Context()
	)

	tokenString := c.QueryParam("token")

	if tokenString == "" {
		resp.Message = "missing or invalid token"
		log.Infof("[UserHandler-1] UpdatePassword: %s", resp.Message)
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}

	if err := c.Bind(&req); err != nil {
		log.Infof("[UserHandler-2] UpdatePassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UserHandler-3] UpdatePassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if req.NewPassword != req.PasswordConfirmation {
		resp.Message = "new password and password confirmation does not match"
		resp.Data = nil
		log.Errorf("[UserHandler-4] UpdatePassword: %v", resp.Message)
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		Password: req.NewPassword,
		Token:    tokenString,
	}

	err = u.UserService.UpdatePassword(ctx, reqEntity)

	if err != nil {
		log.Errorf("[UserHandler-5] UpdatePassword: %v", err)
		if err.Error() == "404" {
			resp.Message = "User not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}

		if err.Error() == "401" {
			resp.Message = "Token expired or invalid"
			resp.Data = nil
			return c.JSON(http.StatusUnauthorized, resp)
		}
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	resp.Data = nil
	resp.Message = "Password updated successfully"

	return c.JSON(http.StatusOK, resp)
}

// VerifyAccount implements UserHandlerInterface.
func (u *UserHandler) VerifyAccount(c echo.Context) error {
	var (
		resp       = response.DefaultReponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	tokenString := c.QueryParam("token")

	if tokenString == "" {
		resp.Message = "missing or invalid token"
		log.Infof("[UserHandler-1] VerifyAccount: %s", resp.Message)
		resp.Data = nil
		return c.JSON(http.StatusUnauthorized, resp)
	}

	user, err := u.UserService.VerifyToken(ctx, tokenString)
	if err != nil {
		log.Errorf("[UserHandler-2] VerifyAccount: %v", err)

		if err.Error() == "404" {
			resp.Message = "User not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}

		if err.Error() == "404" {
			resp.Message = "Token expired or invalid"
			resp.Data = nil
			return c.JSON(http.StatusUnauthorized, resp)
		}

		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Role = user.RoleName
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.Phone = user.Phone
	respSignIn.AccessToken = user.Token

	resp.Message = "Success"
	resp.Data = respSignIn

	return c.JSON(http.StatusOK, resp)
}

// ForgotPassword implements UserHandlerInterface.
func (u *UserHandler) ForgotPassword(c echo.Context) error {
	var (
		req  = request.ForgotPasswordRequest{}
		resp = response.DefaultReponse{}
		ctx  = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] ForgotPassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UserHandler-2] ForgotPassword: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		Email: req.Email,
	}

	err = u.UserService.ForgotPassword(ctx, reqEntity)

	if err != nil {
		log.Errorf("[UserHandler-3] ForgotPassword: %v", err)
		if err.Error() == "404" {
			resp.Message = "User not found"
			resp.Data = nil
			return c.JSON(http.StatusInternalServerError, resp)
		}

		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	resp.Message = "Success"
	resp.Data = nil

	return c.JSON(http.StatusOK, resp)
}

// CreateUserAccount implements UserHandlerInterface.
func (u *UserHandler) CreateUserAccount(c echo.Context) error {
	var (
		req  = request.SignUpRequest{}
		resp = response.DefaultReponse{}
		ctx  = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] CreateUserAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UserHandler-2] CreateUserAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if req.Password != req.PasswordConfirmation {
		log.Errorf("[UserHandler-3] CreateUserAccount: %v", err)
		resp.Message = "Password not Match"
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err = u.UserService.CreateUserAccount(ctx, reqEntity)

	if err != nil {
		log.Errorf("[UserHandler-4] CreateUserAccount: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	resp.Message = "Success"
	resp.Data = nil

	return c.JSON(http.StatusCreated, resp)
}

func (u *UserHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		resp       = response.DefaultReponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UserHandler-2] SignIn: %v", err)

		resp.Message = err.Error()
		resp.Data = nil

		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := u.UserService.SignIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			log.Errorf("[UserHandler-3] SignIn: %s", "User No Found")

			resp.Message = "User Not Found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		log.Errorf("[UserHandler-4] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Role = user.RoleName
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.Phone = user.Phone
	respSignIn.AccessToken = token

	resp.Message = "Success"
	resp.Data = respSignIn

	return c.JSON(http.StatusOK, resp)
}

func NewUserHandler(e *echo.Echo, UserService service.UserServiceInterface, cfg *config.Config) UserHandlerInterface {
	userHandler := &UserHandler{UserService: UserService}

	e.Use(middleware.Recover())

	e.POST("/signin", userHandler.SignIn)
	e.POST("/signup", userHandler.CreateUserAccount)
	e.POST("/forgot-password", userHandler.ForgotPassword)
	e.GET("/verify-account", userHandler.VerifyAccount)
	e.PUT("/update-password", userHandler.UpdatePassword)

	mid := adapter.NewMiddlewareAdapter(cfg)
	adminGroup := e.Group("/admin", mid.CheckToken())

	adminGroup.GET("/check", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	return userHandler
}
