package usecase

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"cdtj.io/days-in-turkey-bot/entity/bot/v2"
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

/*
I feel like logs here are overloaded
*/

// Welcome creates a new user with id and language code,
func (uc *BotUsecase) Welcome(ctx context.Context, userID int64, lang string) *model.TelegramMessage {
	mth := "Welcome"
	// since we do not store anything important, recreating a user on /start is correct behavior
	if err := uc.userUC.Create(ctx, userID, lang); err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "lang", lang, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "lang", lang, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), bot.FmtdMsgWelcome), nil)
}

func (uc *BotUsecase) Country(ctx context.Context, userID int64) *model.TelegramMessage {
	mth := "Country"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "UserCountryPrompt"),
		uc.service.CommandsToInlineKeboard(ctx, uc.service.CountryMarkup(ctx, uc.countryUC.ListFromCache(ctx))))
}

func (uc *BotUsecase) Language(ctx context.Context, userID int64) *model.TelegramMessage {
	mth := "Language"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "UserLanguagePrompt"),
		uc.service.CommandsToInlineKeboard(ctx, uc.service.LanguageMarkup(ctx)))
}

func (uc *BotUsecase) Contribute(ctx context.Context, userID int64) *model.TelegramMessage {
	mth := "Contribute"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "Contribute"), nil)
}

func (uc *BotUsecase) Trip(ctx context.Context, userID int64) *model.TelegramMessage {
	mth := "Trip"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), bot.FmtdMsgTripExplanation), nil)
}

func (uc *BotUsecase) Me(ctx context.Context, userID int64) *model.TelegramMessage {
	mth := "Me"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	userInfo := uc.userUC.GetInfo(ctx, user)
	return model.NewTelegramMessage(userInfo, nil)
}

func (uc *BotUsecase) UpdateLanguage(ctx context.Context, userID int64, languageCodeInput string) *model.TelegramMessage {
	mth := "UpdateLanguage"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", languageCodeInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	languageCodeArr := strings.Split(languageCodeInput, " ")
	languageCode := ""
	if len(languageCodeArr) == 2 && languageCodeArr[0] == "language" {
		languageCode = languageCodeArr[1]
	}
	if err := uc.userUC.UpdateLanguage(ctx, user, languageCode); err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", languageCodeInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInternal"), nil)
	}
	userInfo := uc.userUC.GetInfo(ctx, user)
	msg := uc.service.FormatMessage(ctx, user.GetLanguage(), "UserLanguageChanged")
	return model.NewTelegramMessage(msg+"\n"+userInfo, nil)
}

func (uc *BotUsecase) UpdateCountry(ctx context.Context, userID int64, countryInput string) *model.TelegramMessage {
	mth := "UpdateCountry"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", countryInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	countryArr := strings.Split(countryInput, " ")
	var countryID string
	var daysCont, daysLimit, resetInterval int

	if len(countryArr) == 2 && countryArr[0] == "country" {
		countryID = countryArr[1]
	} else if len(countryArr) == 4 && countryArr[0] == "/custom" {
		daysCont, err = strconv.Atoi(countryArr[1])
		if err != nil {
			slog.Error("usecase failed", "method", mth, "userID", userID, "input", countryInput, "err", err)
			return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInvalidCustomCountry"), nil)
		}
		daysLimit, err = strconv.Atoi(countryArr[2])
		if err != nil {
			slog.Error("usecase failed", "method", mth, "userID", userID, "input", countryInput, "err", err)
			return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInvalidCustomCountry"), nil)
		}
		resetInterval, err = strconv.Atoi(countryArr[3])
		if err != nil {
			slog.Error("usecase failed", "method", mth, "userID", userID, "input", countryInput, "err", err)
			return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInvalidCustomCountry"), nil)
		}
	} else {
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInvalidCustomCountry"), nil)
	}
	country, err := uc.countryUC.Lookup(ctx, countryID, daysCont, daysLimit, resetInterval)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", countryInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInternal"), nil)
	}
	if err := uc.userUC.UpdateCountry(ctx, user, country); err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", countryInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInternal"), nil)
	}
	userInfo := uc.userUC.GetInfo(ctx, user)
	msg := uc.service.FormatMessage(ctx, user.GetLanguage(), "UserCountryChanged")
	return model.NewTelegramMessage(msg+"\n"+userInfo, nil)
}

func (uc *BotUsecase) CalculateTrip(ctx context.Context, userID int64, datesInput string) *model.TelegramMessage {
	mth := "CalculateTrip"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", datesInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	trip, err := uc.userUC.CalculateTrip(ctx, user, datesInput)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", datesInput, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "ErrorInvalidDate"), nil)
	}
	return model.NewTelegramMessage(trip, nil)
}

func (uc *BotUsecase) Hint(ctx context.Context, userID int64, messageCode bot.FmtdMsg) *model.TelegramMessage {
	mth := "Hint"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "input", messageCode, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), messageCode), nil)
}

func (uc *BotUsecase) Feedback(ctx context.Context, userID int64) *model.TelegramMessage {
	mth := "Feedback"
	user, err := uc.userUC.Get(ctx, userID)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return model.NewTelegramMessage(uc.service.FormatMessage(ctx, i18n.DefaultLang(), "ErrorInternal"), nil)
	}
	return model.NewTelegramMessage(uc.service.FormatMessage(ctx, user.GetLanguage(), "Feedback"), nil)
}
