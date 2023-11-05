package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cdtj.io/days-in-turkey-bot/db"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
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

func init() {
	if err := i18n.Localization(); err != nil {
		panic(err)
	}
}

func main() {
	tgFmt := formatter.NewTelegramFormatter()
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)
	countrySvc := cs.NewCountryService(tgFmt)
	countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)

	userDB := db.NewMapDB()
	userRepo := ur.NewUserRepo(userDB)
	userSvc := us.NewUserService(tgFmt)
	userUC := uuc.NewUserUsecase(userRepo, countryRepo, userSvc)

	router := httpserver.NewChiRouter()
	uep.RegisterHTTPEndpoints(router, userUC)
	cep.RegisterHTTPEndpoints(router, countryUC)

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errch := make(chan error)

	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
	go func() {
		rcvd := <-sig
		slog.Info("stopping daemon", "signal", rcvd)
		srv.Shutdown(ctx)
		cancel()
	}()
	go func() {
		err := srv.Serve(ctx)
		if err != nil {
			errch <- fmt.Errorf("server stopped: %w", err)
		}
	}()
	go func() {
		panic(<-errch)
	}()

	<-ctx.Done()
}
