package usecase

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
)

var _ bot.Usecasev2 = NewBotUsecase(nil, nil, nil)

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

func (uc *BotUsecase) Welcome(ctx context.Context, userID int64, lang string) *model.TelegramMessage {
	if err := uc.userUC.Create(ctx, userID, lang); err != nil {
		return model.NewTelegramMessage(err.Error())
	}

	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), bot.FmtdMsgWelcome))
}

func (uc *BotUsecase) Country(ctx context.Context, userID int64) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "UserCountryPrompt"))
	// uc.service.CountryMarkup(ctx, uc.countryUC.ListFromCache(ctx))
}

func (uc *BotUsecase) Language(ctx context.Context, userID int64) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "UserLanguagePrompt"))
	// uc.service.LanguageMarkup(ctx))
}

func (uc *BotUsecase) Contribute(ctx context.Context, userID int64) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "Contribute"))
}

func (uc *BotUsecase) Trip(ctx context.Context, userID int64) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), bot.FmtdMsgTripExplanation))
}

func (uc *BotUsecase) UpdateLanguage(ctx context.Context, userID int64, languageCode string) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	if err := uc.userUC.UpdateLanguage(ctx, user, languageCode); err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	userInfo, err := uc.userUC.GetInfo(ctx, user)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(userInfo)
}

func (uc *BotUsecase) UpdateCountry(ctx context.Context, userID int64, countryID string, daysCont, daysLimit, resetInterval int) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	country, err := uc.countryUC.Lookup(ctx, countryID, daysCont, daysLimit, resetInterval)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	if err := uc.userUC.UpdateCountry(ctx, user, country); err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	userInfo, err := uc.userUC.GetInfo(ctx, user)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(userInfo)
}

func (uc *BotUsecase) CalculateTrip(ctx context.Context, chatID int64, userID int64, datesInput string) *model.TelegramMessage {
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	trip, err := uc.userUC.GetTrip(ctx, user, datesInput)
	if err != nil {
		return model.NewTelegramMessage(err.Error())
	}
	return model.NewTelegramMessage(trip)
}
