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
	defaultCountry = model.NewCountry("RU", "RU", 60, 90, 180)
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

	// registering our handlers and replacing defaults
	botHandlers := bh.BindBotHandlers(botUC)
	botHandlers = append(botHandlers, telegrambot.BindHandlerDefaultDebug(func(format string, args ...any) {
		slog.Debug("bop-api", "msg", fmt.Sprintf(format, args))
	}))
	botHandlers = append(botHandlers, telegrambot.BindHandlerDefaultError(func(err error) {
		slog.Error("bop-api", "err", err)
	}))
	bot := telegrambot.NewTelegramBot(os.Getenv("BOT_TOKEN"), botHandlers)
	if commandsRows := botUC.LocalizeCommands(context.Background(), BotCommands()); commandsRows != nil {
		for _, commandRow := range commandsRows {
			bot.SetCommands(context.Background(), commandRow.Commands, commandRow.LanguageCode)
		}
	}
	if descriptionRows := botUC.LocalizeDescription(context.Background(), BotDescription()); descriptionRows != nil {
		for _, descriptionRow := range descriptionRows {
			bot.SetDescription(context.Background(), descriptionRow.Description, descriptionRow.About, descriptionRow.LanguageCode)
		}
	}

	// using botv2 (based on [github.com/go-telegram/bot]) to read all updates directly without callbacks
	// so we're not using webserver to process with webhooks.
	// mayber we will use http in future for logs
	cmd.Serve(bot, userDB)
}

// BotCommands returns list of known commands,
// feels like it needs to be exported as a config
func BotCommands() []*model.TelegramBotCommand {
	return []*model.TelegramBotCommand{
		model.NewTelegramBotCommand("CommandMe", "me"),
		model.NewTelegramBotCommand("CommandCountry", "country"),
		model.NewTelegramBotCommand("CommandLanguage", "language"),
		model.NewTelegramBotCommand("CommandTrip", "trip"),
		model.NewTelegramBotCommand("CommandContribute", "contribute"),
		model.NewTelegramBotCommand("CommandFeedback", "feedback"),
	}
}

// BotDescription returns placeholders for Bot About and Bot Description used in l10n
func BotDescription() *model.TelegramBotDescription {
	return model.NewTelegramBotDescription("BotDescription", "BotAbout", "")
}
