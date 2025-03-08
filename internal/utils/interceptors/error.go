package interceptors

import (
	"context"
	"errors"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"google.golang.org/grpc"
)

type Error struct {
	Message    string `json:"message"`
	IncidentID string `json:"incident_id,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewChainUnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			var appError *apperrors.Error

			if !errors.As(err, &appError) {
				appError = apperrors.InternalWithoutStackTrace(err)
			}

			return nil, appError
		}

		return resp, nil
	}
}
