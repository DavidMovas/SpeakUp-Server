package echox

import (
	"errors"
	"net/http"

	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"google.golang.org/grpc/codes"

	"github.com/DavidMovas/SpeakUp-Server/contracts"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HTTPError struct {
	Message    string `json:"message"`
	IncidentID string `json:"incident_id,omitempty"`
}

func NewErrorHandler(logger *zap.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var appError *apperrors.Error

		if !errors.As(err, &appError) {
			appError = apperrors.InternalWithoutStackTrace(err)
		}

		httpError := contracts.HTTPError{
			Message:    appError.SafeError(),
			IncidentID: appError.IncidentID,
		}

		if appError.Code == codes.Internal {
			logger.Error("server error",
				zap.String("message", err.Error()),
				zap.String("incident_id", appError.IncidentID),
				zap.String("method", c.Request().Method),
				zap.String("url", c.Request().URL.String()),
				zap.String("stack_trace", appError.StackTrace),
			)
		} else {
			logger.Warn("client error", zap.String("message", err.Error()))
		}

		if err = c.JSON(toHTTPStatus(appError.Code), httpError); err != nil {
			c.Logger().Error(err)
		}
	}
}

func toHTTPStatus(code codes.Code) int {
	switch code {
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
