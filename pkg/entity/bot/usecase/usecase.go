package usecase

import (
	"context"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

var _ bot.Usecase = NewBotUsecase(nil, nil, nil)

type BotUsecase struct {
	usecase user.Usecase
	service bot.Service
	country country.Usecase
}

func NewBotUsecase(usecase user.Usecase, service bot.Service, country country.Usecase) *BotUsecase {
	return &BotUsecase{
		usecase: usecase,
		service: service,
		country: country,
	}
}

func (uc *BotUsecase) Welcome(ctx context.Context, chatID int64, userID, lang string) error {
	if err := uc.usecase.Create(ctx, userID, lang); err != nil {
		return err
	}
	userInfo, err := uc.usecase.Info(ctx, userID)
	if err != nil {
		return err
	}
	storedLang := uc.usecase.GetLang(ctx, userID)
	locale := l10n.GetLocale(storedLang)
	slog.Info("user's lang", "lang", storedLang, "userID", userID, "locale", locale.Name)

	msg := uc.service.FormatMessage(ctx, locale, "Welcome") + " " +
		uc.service.FormatMessage(ctx, locale, "Welcome1") + "\n\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomeCountry") + "\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomeLanguage") + "\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomeTrip") + "\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomeContribute") + "\n\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomePrompt") + "\n\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomePromptPredictEnd") + "\n" +
		uc.service.FormatMessage(ctx, locale, "WelcomePromptPredictRemain")

	uc.Send(ctx, chatID, msg, nil)
	uc.Send(ctx, chatID, userInfo, nil)
	return nil
}

const (
	BotCommandCountry    = "country"
	BotCommandLanguage   = "language"
	BotCommandContribute = "contribute"
	BotCommandTrip       = "trip"
)

func (uc *BotUsecase) Prompt(ctx context.Context, chatID int64, userID, prompt string) error {
	locale := l10n.GetLocale(uc.usecase.GetLang(ctx, userID))
	switch prompt {
	case BotCommandCountry:
		return uc.Send(ctx, chatID, locale.Message("UserCountryPrompt"), uc.CountryMarkup(ctx, userID))
	case BotCommandLanguage:
		return uc.Send(ctx, chatID, locale.Message("UserLanguagePrompt"), uc.LangMarkup(ctx, userID))
	case BotCommandContribute:
		return uc.Send(ctx, chatID, locale.Message("Contribute"), nil)
	case BotCommandTrip:
		return uc.Send(ctx, chatID,
			uc.service.FormatMessage(ctx, locale, "TripExplanation")+"\n\n"+
				uc.service.FormatMessage(ctx, locale, "TripExplanationContinual")+"\n"+
				uc.service.FormatMessage(ctx, locale, "TripExplanationLimit")+"\n"+
				uc.service.FormatMessage(ctx, locale, "TripExplanationResetInterval"),
			nil)
	default:
		return bot.ErrBotCommandNotFound
	}
}

func (uc *BotUsecase) UpdateLang(ctx context.Context, chatID int64, userID, lang string) error {
	if err := uc.usecase.UpdateLang(ctx, userID, lang); err != nil {
		return err
	}
	userInfo, err := uc.usecase.Info(ctx, userID)
	if err != nil {
		return err
	}
	uc.Send(ctx, chatID, userInfo, nil)
	return nil
}

func (uc *BotUsecase) UpdateCountry(ctx context.Context, chatID int64, userID, countryID string) error {
	if err := uc.usecase.UpdateCountry(ctx, userID, countryID); err != nil {
		return err
	}
	userInfo, err := uc.usecase.Info(ctx, userID)
	if err != nil {
		return err
	}
	uc.Send(ctx, chatID, userInfo, nil)
	return nil
}

func (uc *BotUsecase) CalculateTrip(ctx context.Context, chatID int64, userID, input string) error {
	trip, err := uc.usecase.CalculateTrip(ctx, userID, input)
	if err != nil {
		return err
	}
	return uc.Send(ctx, chatID, trip, nil)
}

func (uc *BotUsecase) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	return uc.service.Send(ctx, chatID, text, replyMarkup)
}

func (uc *BotUsecase) CountryMarkup(ctx context.Context, userID string) []*model.TelegramBotCommandRow {
	countries, err := uc.country.List(ctx)
	if err != nil {
		return nil
	}
	commands := make([]*model.TelegramBotCommand, 0, len(countries))
	for _, country := range countries {
		commands = append(commands, model.NewTelegramBotCommand(country.GetFlag()+" "+country.GetName(), "country "+country.GetCode()))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands)}
}

func (uc *BotUsecase) LangMarkup(ctx context.Context, userID string) []*model.TelegramBotCommandRow {
	locales := l10n.Locales()
	commands := make([]*model.TelegramBotCommand, 0, len(locales))
	for _, cmd := range locales {
		commands = append(commands, model.NewTelegramBotCommand(cmd.Name, "language "+cmd.Tag.String()))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands)}
}
