package tests

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"testing"

	"context"
	"fmt"
	"os/signal"
	"syscall"

	"cdtj.io/days-in-turkey-bot/db"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"

	cep "cdtj.io/days-in-turkey-bot/entity/country/endpoint/http"
	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	uep "cdtj.io/days-in-turkey-bot/entity/user/endpoint/http"
	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"
)

func BenchmarkHttpServer(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	incr := atomic.Int64{}

	srv := srvr()
	go Serve(srv)

	cli := cli()
	for _, err := cli.Get("https://1030-188-119-27-163.ngrok.io/user/info/1"); err != nil; {
		b.Error(err)
		return
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cli.Get("https://1030-188-119-27-163.ngrok.io/user/info/" + strconv.FormatInt(incr.Add(1), 10))
		}
	})
}

func cli() *http.Client {
	return &http.Client{}
}

var (
	defaultLang    = "en"
	defaultCountry = model.NewCountry("RU", "RU", 60, 90, 180)
)

func srvr() *httpserver.HttpServer {
	i18n, err := i18n.NewI18n("i18n", defaultLang)
	if err != nil {
		panic(err)
	}
	tgFmt := formatter.NewTelegramFormatter(i18n)
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)
	countrySvc := cs.NewCountryService(tgFmt, defaultCountry)
	countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)

	userDB := db.NewMapDB()
	userRepo := ur.NewUserRepo(userDB)
	userSvc := us.NewUserService(tgFmt, i18n, countrySvc)
	userUC := uuc.NewUserUsecase(userRepo, userSvc, countryUC)

	router := httpserver.NewChiRouter()
	uep.RegisterHTTPEndpoints(router, userUC)
	cep.RegisterHTTPEndpoints(router, countryUC)

	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
	return srv
}

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
			for {
				if err := s.Serve(ctx); err != nil {
					errch <- fmt.Errorf("%T stopped with error: %w", s, err)
				}
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
		slog.Error("serving error", "error", <-errch)
	}()
	<-ctx.Done()
}
