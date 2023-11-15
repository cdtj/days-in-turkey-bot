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
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	bwh "cdtj.io/days-in-turkey-bot/entity/bot/endpoint/web-hook"
	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"

	bs "cdtj.io/days-in-turkey-bot/entity/bot/service"
	buc "cdtj.io/days-in-turkey-bot/entity/bot/usecase"
)

var (
	defaultLang    = "en"
	defaultCountry = model.NewCountry("RU", "RU", 60, 90, 180)
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	i18n, err := i18n.NewI18n("i18n", defaultLang)
	if err != nil {
		panic(err)
	}

	telegramFrmtr := formatter.NewTelegramFormatter(i18n)

	// country service
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)
	countrySvc := cs.NewCountryService(telegramFrmtr, defaultCountry)
	countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)

	// user service
	userDB := db.NewBoltDB("users", "users")
	userRepo := ur.NewUserRepo(userDB)
	userSvc := us.NewUserService(telegramFrmtr, i18n, countrySvc)
	userUC := uuc.NewUserUsecase(userRepo, userSvc, countryUC)

	// telegram bot
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), os.Getenv("BOT_WEBHOOK"))
	botSvc := bs.NewBotService(bot, telegramFrmtr, i18n)
	botUC := buc.NewBotUsecase(botSvc, userUC, countryUC)

	router := httpserver.NewEchoRouter()
	bwh.RegisterWebhookEndpointsEcho(router, botUC)

	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
	Serve(srv, bot, userDB)
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
