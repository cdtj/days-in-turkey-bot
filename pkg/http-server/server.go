package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(server *http.Server) *HttpServer {
	return &HttpServer{
		server: server,
	}
}

func (s *HttpServer) Serve(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	shutdownCtx, shutdownStopCtx := context.WithTimeout(ctx, 15*time.Second)
	defer shutdownStopCtx()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			slog.Error("unable to gracefully stop telegram bot", "error", shutdownCtx.Err())
			return
		}
	}()

	slog.Info("shutting down the server", "Addr", s.server.Addr)
	err := s.server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("unable to gracefully stop http server", "error", err)
		return
	}
	slog.Info("server is down", "Addr", s.server.Addr)
}
