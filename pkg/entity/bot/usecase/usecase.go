package usecase

import (
	"context"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
)

var _ bot.Usecase = NewBotUsecase(nil, nil, nil)

type BotUsecase struct {
	service   bot.Service
	userUC    user.Usecase
	countryUC country.Usecase
}

func NewBotUsecase(service bot.Service, userUC user.Usecase, countryUC country.Usecase) *BotUsecase {
	return &BotUsecase{
		service:   service,
		userUC:    userUC,
		countryUC: countryUC,
	}
}

func (uc *BotUsecase) Welcome(ctx context.Context, chatID int64, userID, lang string) error {
	if err := uc.userUC.Create(ctx, userID, lang); err != nil {
		return err
	}
	userInfo, err := uc.userUC.Info(ctx, userID)
	if err != nil {
		return err
	}
	storedLang := uc.userUC.GetLang(ctx, userID)
	locale := i18n.GetLocale(storedLang)
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
	locale := i18n.GetLocale(uc.userUC.GetLang(ctx, userID))
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
	if err := uc.userUC.UpdateLang(ctx, userID, lang); err != nil {
		return err
	}
	userInfo, err := uc.userUC.Info(ctx, userID)
	if err != nil {
		return err
	}
	uc.Send(ctx, chatID, userInfo, nil)
	return nil
}

func (uc *BotUsecase) UpdateCountry(ctx context.Context, chatID int64, userID, countryID string) error {
	if err := uc.userUC.UpdateCountry(ctx, userID, countryID); err != nil {
		return err
	}
	userInfo, err := uc.userUC.Info(ctx, userID)
	if err != nil {
		return err
	}
	uc.Send(ctx, chatID, userInfo, nil)
	return nil
}

func (uc *BotUsecase) CalculateTrip(ctx context.Context, chatID int64, userID, input string) error {
	trip, err := uc.userUC.CalculateTrip(ctx, userID, input)
	if err != nil {
		return err
	}
	return uc.Send(ctx, chatID, trip, nil)
}

func (uc *BotUsecase) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	return uc.service.Send(ctx, chatID, text, replyMarkup)
}
