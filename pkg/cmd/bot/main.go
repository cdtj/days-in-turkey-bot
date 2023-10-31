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
	"cdtj.io/days-in-turkey-bot/service/l10n"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"

	bwh "cdtj.io/days-in-turkey-bot/entity/bot/endpoint/web-hook"
	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"

	bs "cdtj.io/days-in-turkey-bot/entity/bot/service"
	buc "cdtj.io/days-in-turkey-bot/entity/bot/usecase"
)

func init() {
	if err := l10n.Localization(); err != nil {
		panic(err)
	}
}

func main() {
	// country service
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)

	// user service
	userDB := db.NewMapDB()
	userRepo := ur.NewUserRepo(userDB)
	userSvc := us.NewUserService(formatter.NewTelegramFormatter())
	userUC := uuc.NewUserUsecase(userRepo, countryRepo, userSvc)

	// telegram bot
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), os.Getenv("BOT_WEBHOOK"))
	botSvc := bs.NewBotService(bot)
	botUC := buc.NewBotUsecase(userUC, botSvc)

	router := httpserver.NewChiRouter()
	bwh.RegisterWebhookEndpoints(router, botUC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errch := make(chan error)

	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
	go func() {
		err := srv.Serve(ctx)
		if err != nil {
			errch <- fmt.Errorf("server stopped: %w", err)
		}
	}()

	go func() {
		err := bot.Serve(ctx)
		if err != nil {
			errch <- fmt.Errorf("bot stopped: %w", err)
		}
	}()

	go func() {
		rcvd := <-sig
		slog.Info("stopping daemon", "signal", rcvd)
		srv.Shutdown(ctx)
		bot.Shutdown(ctx)
		cancel()
	}()
	go func() {
		panic(<-errch)
	}()

	<-ctx.Done()
}
