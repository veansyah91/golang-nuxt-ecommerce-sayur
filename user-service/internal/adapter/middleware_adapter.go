package adapter

import (
	"net/http"
	"strings"
	"user-service/config"
	"user-service/internal/adapter/handler/response"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}

type MiddlewareAdapter struct {
	cfg *config.Config
}

// CheckToken implements MiddlewareAdapterInterface.
func (m *MiddlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			respErr := response.DefaultReponse{}
			redisConn := config.NewRedisClient()

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Errorf("[MiddlewareAdapter-1] CheckToken: %s", "missing or invalid token")
				respErr.Message = "missing or invalid token"
				respErr.Data = nil
				return c.JSON(http.StatusUnauthorized, respErr)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			getSession, err := redisConn.HGetAll(c.Request().Context(), tokenString).Result()

			if err != nil {
				log.Errorf("[MiddlewareAdapter-2] Redis error: %v", err)
				respErr.Message = "failed to verify session"
				respErr.Data = nil
				return c.JSON(http.StatusUnauthorized, respErr)
			}

			if len(getSession) == 0 {
				log.Info(getSession)
				log.Warnf("[MiddlewareAdapter-3] session not found or invalid token")
				respErr.Message = "session not found or invalid token"
				respErr.Data = nil
				return c.JSON(http.StatusUnauthorized, respErr)
			}

			c.Set("user", getSession)
			return next(c)
		}
	}

}

func NewMiddlewareAdapter(cfg *config.Config) MiddlewareAdapterInterface {
	return &MiddlewareAdapter{
		cfg: cfg,
	}
}
