package main

import (
	"log/slog"
	"net/http"
	"os"

	"cdtj.io/days-in-turkey-bot/cmd"
	"cdtj.io/days-in-turkey-bot/db"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	bwh "cdtj.io/days-in-turkey-bot/entity/bot/endpoint/web-hook/echo"
	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"

	bs "cdtj.io/days-in-turkey-bot/entity/bot/service"
	buc "cdtj.io/days-in-turkey-bot/entity/bot/usecase"
)

var (
	defaultLang    = "en"
	defaultCountry = model.NewCountry("CUSTOM", "📝", "", 60, 90, 180, true)
)

// bot (v1) is deprecated even before become public but
// it might be interested to return back to it when the origin library
// will be updated to the actual API (if it)
func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	i18n, err := i18n.NewI18n("i18n", defaultLang)
	if err != nil {
		panic(err)
	}

	telegramFrmtr := formatter.NewTelegramFormatter(i18n, false)

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
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), os.Getenv("BOT_WEBHOOK"))
	botSvc := bs.NewBotService(bot, telegramFrmtr, i18n)
	botUC := buc.NewBotUsecase(botSvc, userUC, countryUC)

	// we are using webserver to deploy telegram webhook,
	// please note that bot (v1) based on [github.com/go-telegram-bot-api/telegram-bot-api/v5]
	// library that doesn't implement webhook secret header, so anyone can post data
	// if they know your host and hanlder
	router := httpserver.NewEchoRouter()
	bwh.RegisterWebhookEndpointsEcho(router, botUC)

	// testing on insecure
	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
	cmd.Serve(srv, bot, userDB)
}
