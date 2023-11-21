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
			slog.Error("server", "Addr", s.server.Addr, "msg", "unable to gracefully stop http server", "error", shutdownCtx.Err())
			return
		}
	}()

	slog.Info("server", "Addr", s.server.Addr, "status", "stopping")
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server", "Addr", s.server.Addr, "msg", "unable to gracefully stop http server", "error", err)
		return
	}
	slog.Info("server", "Addr", s.server.Addr, "status", "stopped")
}
