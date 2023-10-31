package usecase

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/entity/user"
)

type BotUsecase struct {
	usecase user.Usecase
	service bot.Service
}

func NewBotUsecase(usecase user.Usecase, service bot.Service) *BotUsecase {
	return &BotUsecase{
		usecase: usecase,
		service: service,
	}
}

func (uc *BotUsecase) UpdateLang(ctx context.Context, userID, lang string) (string, error) {
	return uc.usecase.UpdateLang(ctx, userID, lang)
}

func (uc *BotUsecase) UpdateCountry(ctx context.Context, userID, countryID string) (string, error) {
	return uc.usecase.UpdateCountry(ctx, userID, countryID)
}

func (uc *BotUsecase) CalculateTrip(ctx context.Context, userID, input string) (string, error) {
	return uc.usecase.CalculateTrip(ctx, userID, input)
}

func (uc *BotUsecase) Reply(ctx context.Context, chatID int64, text string) error {
	return uc.service.Reply(ctx, chatID, text)
}
