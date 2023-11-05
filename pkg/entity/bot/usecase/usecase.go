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
	uc.Send(ctx, chatID, userInfo, nil)
	uc.Send(ctx, chatID, "Select Country:", uc.service.CountryMarkup(ctx, uc.country.Cache(ctx)))
	uc.Send(ctx, chatID, "Select Language:", uc.service.LangMarkup(ctx))
	return nil
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
