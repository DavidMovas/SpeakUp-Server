package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewLoggingMiddleware(logger *zap.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			now := time.Now()

			err := next(c)

			logger.Debug("Request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("query", c.Request().URL.RawQuery),
				zap.String("remote_addr", c.Request().RemoteAddr),
				zap.String("user_agent", c.Request().UserAgent()),
				zap.String("referer", c.Request().Referer()),
				zap.Duration("duration", time.Since(now)),
			)

			return err
		}
	}
}
