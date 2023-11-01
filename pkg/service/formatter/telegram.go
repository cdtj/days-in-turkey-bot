package formatter

import (
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

type TelegramFormatter struct {
}

func NewTelegramFormatter() *TelegramFormatter {
	return &TelegramFormatter{}
}

var _ Formatter = NewTelegramFormatter()

func (f *TelegramFormatter) TripTree(locale *l10n.Locale, tree *model.TripTree) string {
	result := ""
	firstLine := true
	firstEligible := true
	for i := tree; i != nil; i = i.Prev {
		if i.StartPredicted {
			if firstEligible {
				result += locale.Message("TripEligibleHdr") + "\n"
				firstEligible = false
			}
			result += locale.MessageWithTemplate("TripPredicted", map[string]interface{}{
				"StartDate":  locale.FormatDate(i.StartDate),
				"EndDate":    locale.FormatDate(i.EndDate),
				"TripDays":   i.TripDays,
				"PeriodDays": i.PeriodDays,
			}, nil) + locale.MessageWithCount("DayCounter", i.PeriodDays) + "\n"
		} else if i.EndPredicted {
			result += locale.MessageWithTemplate("TripPredicted", map[string]interface{}{
				"StartDate":  locale.FormatDate(i.StartDate),
				"EndDate":    locale.FormatDate(i.EndDate),
				"TripDays":   i.TripDays,
				"PeriodDays": i.PeriodDays,
			}, nil) + locale.MessageWithCount("DayCounter", i.PeriodDays) + "\n"
		} else {
			if firstLine {
				result += "\n" + locale.Message("TripPast") + "\n"
				firstLine = false
			}
			result += locale.MessageWithTemplate("Trip", map[string]interface{}{
				"StartDate":  locale.FormatDate(i.StartDate),
				"EndDate":    locale.FormatDate(i.EndDate),
				"TripDays":   i.TripDays,
				"PeriodDays": i.PeriodDays,
			}, nil) + locale.MessageWithCount("DayCounter", i.PeriodDays) + "\n"
		}
	}
	return result
}

func (f *TelegramFormatter) User(locale *l10n.Locale, user *model.User) string {
	return locale.MessageWithTemplate("UserInfo", map[string]interface{}{
		"Language": user.GetLang(),
	}, nil) + "\n" + f.Country(locale, user.Country)
}

func (f *TelegramFormatter) Country(locale *l10n.Locale, country *model.Country) string {
	return locale.MessageWithTemplate("CountryInfo", map[string]interface{}{
		"Code": country.GetCode(),
		"Flag": country.GetFlag(),
	}, nil) + "\n" + locale.MessageWithTemplate("CountryDays", map[string]interface{}{
		"Continual":     country.GetDaysCont(),
		"Limit":         country.GetDaysLimit(),
		"ResetInterval": country.GetResetInterval(),
	}, nil)
}
