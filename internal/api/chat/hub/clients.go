package hub

import (
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"google.golang.org/grpc"
)

type client struct {
	userID string
	grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]
}

func newClient(userID string, stream grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]) *client {
	return &client{
		userID:              userID,
		BidiStreamingServer: stream,
	}
}
