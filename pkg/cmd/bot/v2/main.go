package main

import (
	"log/slog"
	"os"

	"cdtj.io/days-in-turkey-bot/cmd"
	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"

	bh "cdtj.io/days-in-turkey-bot/entity/bot/v2/endpoint/tg-handler"
	bs "cdtj.io/days-in-turkey-bot/entity/bot/v2/service"
	buc "cdtj.io/days-in-turkey-bot/entity/bot/v2/usecase"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot/v2"
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
	botSvc := bs.NewBotService(telegramFrmtr, i18n)
	botUC := buc.NewBotUsecase(botSvc, userUC, countryUC)
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), bh.BindBotHandlers(botUC))

	// using botv2 (based on [github.com/go-telegram/bot]) to read all updates directly without callbacks
	// so we're not using webserver to process with webhooks
	cmd.Serve(bot, userDB)
}
