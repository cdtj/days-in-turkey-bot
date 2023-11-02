package usecase

import (
	"context"

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
	uc.Send(ctx, chatID, userInfo, nil)
	uc.Send(ctx, chatID, "Select Country:", uc.CountryMarkup(ctx, userID))
	uc.Send(ctx, chatID, "Select Language:", uc.LangMarkup(ctx, userID))
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

func (uc *BotUsecase) CountryMarkup(ctx context.Context, userID string) []*model.TelegramBotCommandRow {
	countries, err := uc.country.Keys(ctx)
	if err != nil {
		return nil
	}
	commands := make([]*model.TelegramBotCommand, 0, len(countries))
	for _, country := range countries {
		commands = append(commands, model.NewTelegramBotCommand(country, "country "+country))
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
