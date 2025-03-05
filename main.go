package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DavidMovas/SpeakUp-Server/internal"
	"github.com/DavidMovas/SpeakUp-Server/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Error loading config: ", "error", err)
		return
	}

	startCtx, startCancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer startCancel()

	srv, err := internal.NewServer(startCtx, cfg)
	if err != nil {
		slog.Error("Error starting server: ", "error", err)
		return
	}

	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		<-signalCh
		slog.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer shutdownCancel()

		if err = srv.Shutdown(shutdownCtx); err != nil {
			slog.Warn("Server forced to shutdown", "error", err)
		}
	}()

	if err = srv.Start(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed to start", "error", err)
	}
}
