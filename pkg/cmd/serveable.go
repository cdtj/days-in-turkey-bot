package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Serveable maintaining daemons for us
type Serveable interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context)
}

func Serve(serveables ...Serveable) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errch := make(chan error)

	for k, s := range serveables {
		slog.Info("starting", "key", k, "s", s)
		go func(s Serveable) {
			if err := s.Serve(ctx); err != nil {
				errch <- fmt.Errorf("%T stopped with error: %w", s, err)
			}
		}(s)
	}

	go func() {
		rcvd := <-sig
		slog.Info("stopping serveables", "signal", rcvd)
		for _, s := range serveables {
			s.Shutdown(ctx)
		}
		cancel()
	}()
	go func() {
		panic(<-errch)
	}()
	<-ctx.Done()
}
