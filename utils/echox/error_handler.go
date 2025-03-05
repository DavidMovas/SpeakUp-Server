package echox

import (
	"errors"
	"github.com/DavidMovas/SpeakUp-Server/contracts"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/error"
	"github.com/DavidMovas/SpeakUp-Server/internal/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type HTTPError struct {
	Message    string `json:"message"`
	IncidentID string `json:"incident_id,omitempty"`
}

func ErrorHandler(err error, c echo.Context) {
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

	logger, ok := log.FromContext(c.Request().Context())

	if appError.Code == apperrors.InternalCode {
		logger.Error("server error",
			zap.String("message", err.Error()),
			zap.String("incident_id", appError.IncidentID),
			zap.String("method", c.Request().Method),
			zap.String("url", c.Request().URL.String()),
			zap.String("stack_trace", appError.StackTrace),
		)
	} else {
		if ok {
			logger.Warn("client error", zap.String("message", err.Error()))
		}
	}

	if err = c.JSON(toHTTPStatus(appError.Code), httpError); err != nil {
		c.Logger().Error(err)
	}
}

func toHTTPStatus(code apperrors.Code) int {
	switch code {
	case apperrors.InternalCode:
		return http.StatusInternalServerError
	case apperrors.BadRequestCode:
		return http.StatusBadRequest
	case apperrors.NotFoundCode:
		return http.StatusNotFound
	case apperrors.UnauthorizedCode:
		return http.StatusUnauthorized
	case apperrors.ForbiddenCode:
		return http.StatusForbidden
	case apperrors.AlreadyExistsCode, apperrors.VersionMismatchCode:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
