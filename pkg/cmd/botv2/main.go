package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"

	tghandler "cdtj.io/days-in-turkey-bot/entity/bot/endpoint/tg-handler"
	bs "cdtj.io/days-in-turkey-bot/entity/bot/service"
	buc "cdtj.io/days-in-turkey-bot/entity/bot/usecase/v2"
)

var (
	defaultLang    = "en"
	defaultCountry = model.NewCountry("RU", "RU", 60, 90, 180)
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	i18n, err := i18n.NewI18n("i18n", defaultLang)
	if err != nil {
		panic(err)
	}

	telegramFrmtr := formatter.NewTelegramFormatter(i18n, true)

	// country service
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)
	countrySvc := cs.NewCountryService(telegramFrmtr, defaultCountry)
	countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)

	// user service
	userDB := db.NewBoltDB("users", "users")
	userRepo := ur.NewUserRepo(ur.NewUserBoltDBAdaptor(userDB))
	userSvc := us.NewUserService(telegramFrmtr, i18n, countrySvc)
	userUC := uuc.NewUserUsecase(userRepo, userSvc, countryUC)

	// telegram bot
	botSvc := bs.NewBotServicev2(telegramFrmtr, i18n)
	botUC := buc.NewBotUsecase(botSvc, userUC, countryUC)
	bot := telegrambot.NewTelegramBotv2(os.Getenv("BOT_TOKEN"), tghandler.BindBotHandlers(botUC))

	Serve(bot, userDB)
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
