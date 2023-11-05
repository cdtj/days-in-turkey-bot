package usecase

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
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

func (uc *BotUsecase) Welcome(ctx context.Context, chatID int64, userID int64, lang string) error {
	if err := uc.userUC.Create(ctx, userID, lang); err != nil {
		return err
	}

	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return err
	}
	msg := uc.service.FormatMessage(ctx, user.GetLanguage(), bot.FmtdMsgWelcome)
	if err := uc.Send(ctx, chatID, msg, nil); err != nil {
		return err
	}

	userInfo, err := uc.userUC.GetInfo(ctx, user)
	if err != nil {
		return err
	}
	if err := uc.Send(ctx, chatID, userInfo, nil); err != nil {
		return err
	}

	return nil
}

const (
	BotCommandCountry    = "country"
	BotCommandLanguage   = "language"
	BotCommandContribute = "contribute"
	BotCommandTrip       = "trip"
)

func (uc *BotUsecase) Prompt(ctx context.Context, chatID int64, userID int64, prompt string) error {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return err
	}

	switch prompt {
	case BotCommandCountry:
		return uc.Send(ctx, chatID, uc.service.FormatMessage(ctx, user.GetLanguage(), "UserCountryPrompt"), uc.service.CountryMarkup(ctx, uc.countryUC.ListFromCache(ctx)))
	case BotCommandLanguage:
		return uc.Send(ctx, chatID, uc.service.FormatMessage(ctx, user.GetLanguage(), "UserLanguagePrompt"), uc.service.LanguageMarkup(ctx))
	case BotCommandContribute:
		return uc.Send(ctx, chatID, uc.service.FormatMessage(ctx, user.GetLanguage(), "Contribute"), nil)
	case BotCommandTrip:
		return uc.Send(ctx, chatID,
			uc.service.FormatMessage(ctx, user.GetLanguage(), bot.FmtdMsgTripExplanation),
			nil)
	default:
		return bot.ErrBotCommandNotFound
	}
}

func (uc *BotUsecase) UpdateLanguage(ctx context.Context, chatID int64, userID int64, lang string) error {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return err
	}

	if err := uc.userUC.UpdateLanguage(ctx, user, lang); err != nil {
		return err
	}
	userInfo, err := uc.userUC.GetInfo(ctx, user)
	if err != nil {
		return err
	}
	return uc.Send(ctx, chatID, userInfo, nil)
}

func (uc *BotUsecase) UpdateCountry(ctx context.Context, chatID int64, userID int64, countryID string, daysCont, daysLimit, resetInterval int) error {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return err
	}

	country, err := uc.countryUC.Lookup(ctx, countryID, daysCont, daysLimit, resetInterval)
	if err != nil {
		return err
	}
	if err := uc.userUC.UpdateCountry(ctx, user, country); err != nil {
		return err
	}
	userInfo, err := uc.userUC.GetInfo(ctx, user)
	if err != nil {
		return err
	}
	return uc.Send(ctx, chatID, userInfo, nil)
}

func (uc *BotUsecase) CalculateTrip(ctx context.Context, chatID int64, userID int64, datesInput string) error {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return err
	}

	trip, err := uc.userUC.GetTrip(ctx, user, datesInput)
	if err != nil {
		return err
	}
	return uc.Send(ctx, chatID, trip, nil)
}

func (uc *BotUsecase) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	return uc.service.Send(ctx, chatID, text, replyMarkup)
}
