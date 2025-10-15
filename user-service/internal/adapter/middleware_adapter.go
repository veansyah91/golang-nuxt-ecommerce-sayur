package adapter

import (
	"net/http"
	"strings"
	"user-service/config"

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
			redisConn := config.NewConfig().NewRedisClient()

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Errorf("[MiddlewareAdapter-1] CheckToken: %s", "missing or invalid token")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			getSession, err := redisConn.HGetAll(c.Request().Context(), tokenString).Result()
			if err != nil || len(getSession) == 0 {
				msg := "session not found or invalid token"

				if err != nil {
					msg = err.Error()
					log.Errorf("[MiddlewareAdapter-3] CheckToken error: %s", msg)
				} else {
					log.Warnf("[MiddlewareAdapter-3] CheckToken warning: %s", msg)
				}

				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
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
