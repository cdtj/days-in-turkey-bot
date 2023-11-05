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

func init() {
	if err := i18n.I18n(); err != nil {
		panic(err)
	}
}

func main() {
	// country service
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)
	countrySvc := cs.NewCountryService(formatter.NewTelegramFormatter())
	countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)

	// user service
	userDB := db.NewBoltDB("users", "users")
	userRepo := ur.NewUserRepo(userDB)
	userSvc := us.NewUserService(formatter.NewTelegramFormatter())
	userUC := uuc.NewUserUsecase(userRepo, countryRepo, userSvc)

	// telegram bot
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), os.Getenv("BOT_WEBHOOK"))
	botSvc := bs.NewBotService(bot, formatter.NewTelegramFormatter())
	botUC := buc.NewBotUsecase(userUC, botSvc, countryUC)

	router := httpserver.NewChiRouter()
	bwh.RegisterWebhookEndpoints(router, botUC)

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
