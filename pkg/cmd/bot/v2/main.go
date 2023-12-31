package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
	defaultCountry = model.NewCountry("CUSTOM", "📝", "/country", 60, 90, 180, true)
)

func main() {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		if a.Value.Kind().String() == "Time" && a.Key != slog.TimeKey {
			return slog.Attr{
				Key:   a.Key,
				Value: slog.StringValue(a.Value.Time().Format("2006-01-02")),
			}
		}
		return a
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replace,
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
	dbPath := os.Getenv("BOLDTB_PATH")
	if dbPath == "" {
		dbPath = "/db/"
	}
	userDB := db.NewBoltDB(dbPath+"users", "users")
	userRepo := ur.NewUserRepo(ur.NewUserBoltDBAdaptor(userDB))
	userSvc := us.NewUserService(telegramFrmtr, i18n, countrySvc)
	userUC := uuc.NewUserUsecase(userRepo, userSvc, countryUC)

	// telegram bot
	botSvc := bs.NewBotService(telegramFrmtr, i18n)
	botUC := buc.NewBotUsecase(botSvc, userUC, countryUC)

	botOptions := telegrambot.NewTelegramBotOptions(bh.BindBotHandlers(botUC),
		model.NewTelegramBotDescription("BotDescription", "BotAbout"),
		nil,
		func(err error) {
			slog.Error("bop-api", "err", err)
		},
		func(format string, args ...any) {
			slog.Debug("bop-api", "msg", fmt.Sprintf(format, args))
		},
	)
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), botOptions)

	locales := i18n.Locales()

	for _, locale := range locales {
		lCommands := LocalizeCommands(locale, botOptions.GetCommands())
		lDescription := LocalizeDescription(locale, botOptions.GetDescription())
		bot.SetCommands(context.Background(), lCommands.GetCommands(), locale.GetLanguage())
		bot.SetDescription(context.Background(), lDescription.GetDescription(), lDescription.GetAbout(), locale.GetLanguage())
	}

	// using botv2 (based on [github.com/go-telegram/bot]) to read all updates directly without callbacks
	// so we're not using webserver to process with webhooks.
	// mayber we will use http in future for logs
	cmd.Serve(bot, userDB)
}

// LocalizeCommands translates commands using i18n'er and adding Language Tag
func LocalizeCommands(locale *i18n.Locale, commands []*model.TelegramBotCommand) *model.TelegramBotCommandRow {
	return locale.LocalizeCommands(commands)
}

// LocalizeDescription translates description using i18n'er and adding Language Tag
func LocalizeDescription(locale *i18n.Locale, description *model.TelegramBotDescription) *model.TelegramBotDescription {
	return locale.LocalizeDescription(description)
}
