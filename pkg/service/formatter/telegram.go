package formatter

import (
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"
)

type TelegramFormatter struct {
	i18n i18n.I18ner
}

func NewTelegramFormatter(i18n i18n.I18ner) *TelegramFormatter {
	return &TelegramFormatter{
		i18n: i18n,
	}
}

var _ Formatter = NewTelegramFormatter(nil)

func (f *TelegramFormatter) TripTree(language language.Tag, tree *model.TripTree) string {
	result := ""
	firstLine := true
	firstEligible := true
	locale := f.i18n.GetLocale(language)
	for i := tree; i != nil; i = i.Prev {
		if i.StartPredicted {
			if firstEligible {
				result += locale.Message("TripEligibleHdr") + "\n"
				firstEligible = false
			}
			result += locale.MessageWithTemplate("TripPredicted", map[string]interface{}{
				"StartDate":  wrapCode(locale.FormatDate(i.StartDate)),
				"EndDate":    wrapCode(locale.FormatDate(i.EndDate)),
				"TripDays":   locale.MessageWithCount("DayCounter", i.TripDays),
				"PeriodDays": locale.MessageWithCount("DayCounter", i.PeriodDays),
			}, nil) + "\n"
		} else if i.EndPredicted {
			result += locale.MessageWithTemplate("TripPredicted", map[string]interface{}{
				"StartDate":  wrapCode(locale.FormatDate(i.StartDate)),
				"EndDate":    wrapCode(locale.FormatDate(i.EndDate)),
				"TripDays":   locale.MessageWithCount("DayCounter", i.TripDays),
				"PeriodDays": locale.MessageWithCount("DayCounter", i.PeriodDays),
			}, nil) + "\n"
		} else {
			if firstLine {
				result += "\n" + locale.Message("TripPast") + "\n"
				firstLine = false
			}
			result += locale.MessageWithTemplate("Trip", map[string]interface{}{
				"StartDate":  wrapCode(locale.FormatDate(i.StartDate)),
				"EndDate":    wrapCode(locale.FormatDate(i.EndDate)),
				"TripDays":   locale.MessageWithCount("DayCounter", i.TripDays),
				"PeriodDays": locale.MessageWithCount("DayCounter", i.PeriodDays),
			}, nil) + "\n"
		}
	}
	return result
}

func (f *TelegramFormatter) User(language language.Tag, user *model.User) string {
	locale := f.i18n.GetLocale(language)
	return locale.MessageWithTemplate("UserInfo", map[string]interface{}{
		"Language": locale.Message("LanguageName"),
	}, nil) + "\n" + f.Country(language, &user.Country)
}

func (f *TelegramFormatter) Country(language language.Tag, country *model.Country) string {
	locale := f.i18n.GetLocale(language)
	return locale.MessageWithTemplate("CountryInfo", map[string]interface{}{
		"Flag": country.GetFlag(),
		"Name": country.GetName(),
	}, nil) + "\n" + locale.MessageWithTemplate("CountryDays", map[string]interface{}{
		"Continual":     country.GetDaysCont(),
		"Limit":         country.GetDaysLimit(),
		"ResetInterval": country.GetResetInterval(),
	}, nil)
}

func (f *TelegramFormatter) FormatMessage(language language.Tag, messageID string) string {
	locale := f.i18n.GetLocale(language)
	return locale.Message(messageID)
}

func (f *TelegramFormatter) Welcome(language language.Tag) string {
	locale := f.i18n.GetLocale(language)
	return locale.Message("Welcome") + " " +
		locale.Message("Welcome1") + "\n\n" +
		locale.Message("WelcomeCountry") + "\n" +
		locale.Message("WelcomeLanguage") + "\n" +
		locale.Message("WelcomeTrip") + "\n" +
		locale.Message("WelcomeContribute") + "\n\n" +
		locale.Message("WelcomePrompt") + "\n\n" +
		locale.Message("WelcomePromptPredictEnd") + "\n" +
		locale.Message("WelcomePromptPredictRemain")
}

func (f *TelegramFormatter) TripExplanation(language language.Tag) string {
	locale := f.i18n.GetLocale(language)
	return locale.Message("TripExplanation") + "\n\n" +
		locale.Message("TripExplanationContinual") + "\n" +
		locale.Message("TripExplanationLimit") + "\n" +
		locale.Message("TripExplanationResetInterval")
}

func wrapCode(str string) string {
	return "`" + str + "`"
}

func wrapItalic(str string) string {
	return "_" + str + "_"
}

func wrapBold(str string) string {
	return "*" + str + "*"
}
