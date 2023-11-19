package formatter

import (
	"log/slog"
	"strings"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	tgapi "github.com/go-telegram/bot"
	"golang.org/x/text/language"
)

type TelegramFormatter struct {
	i18n i18n.I18ner
	v2   bool
}

func NewTelegramFormatter(i18n i18n.I18ner, v2 bool) *TelegramFormatter {
	return &TelegramFormatter{
		i18n: i18n,
		v2:   v2,
	}
}

var _ Formatter = NewTelegramFormatter(nil, false)

func (f *TelegramFormatter) TripTree(language language.Tag, tree *model.TripTree) string {
	result := ""
	firstLine := true
	firstEligible := true
	overstayed := false
	locale := f.i18n.GetLocale(language)
	for i := tree; i != nil; i = i.Prev {
		slog.Info("TripTree", "StartDate", i.StartDate, "EndDate", i.EndDate, "TripDays", i.TripDays, "PeriodDays", i.PeriodDays, "OverstayDays", i.OverstayDays)
		if i.StartPredicted || i.EndPredicted {
			if i.StartPredicted && firstEligible {
				result += locale.Message("TripEligibleHdr") + "\n"
				firstEligible = false
			}
			result += locale.MessageWithTemplate("TripPredicted", map[string]interface{}{
				"StartDate":  wrapCode(locale.FormatDate(i.StartDate)),
				"EndDate":    wrapCode(locale.FormatDate(i.EndDate)),
				"TripDays":   locale.MessageWithCount("DayCounter", i.TripDays),
				"PeriodDays": locale.MessageWithCount("DayCounter", i.PeriodDays),
			}, nil)
			if i.OverstayDays > 0 {
				overstayed = true
				result += ", " + locale.MessageWithTemplate("Overstay", map[string]interface{}{
					"OverstayDays": wrapBold(locale.MessageWithCount("DayCounter", i.OverstayDays)),
				}, nil) + " ⚠️"
			}
			result += "\n"
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
			}, nil)
			if i.OverstayDays > 0 {
				overstayed = true
				result += ", " + locale.MessageWithTemplate("Overstay", map[string]interface{}{
					"OverstayDays": wrapBold(locale.MessageWithCount("DayCounter", i.OverstayDays)),
				}, nil) + " ⚠️"
			}
			result += "\n"
		}
	}
	result += "\n\n" + locale.Message("OverstayCaution")
	if overstayed {
		result += "\n" + locale.Message("OverstayExplanation")
	}
	return f.markdownWrapper(result)
}

func (f *TelegramFormatter) User(language language.Tag, user *model.User) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.MessageWithTemplate("UserInfo", map[string]interface{}{
		"Language": locale.Message("LanguageName"),
	}, nil)) +
		"\n" +
		f.Country(language, &user.Country)
}

func (f *TelegramFormatter) Country(language language.Tag, country *model.Country) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.MessageWithTemplate("CountryInfo", map[string]interface{}{
		"Flag": country.GetFlag(),
		"Name": country.GetName(),
	}, nil) +
		"\n" +
		locale.MessageWithTemplate("CountryDays", map[string]interface{}{
			"Continual":     country.GetDaysCont(),
			"Limit":         country.GetDaysLimit(),
			"ResetInterval": country.GetResetInterval(),
		}, nil))
}

func (f *TelegramFormatter) FormatMessage(language language.Tag, messageID string) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.Message(messageID))
}

func (f *TelegramFormatter) Welcome(language language.Tag) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.Message("Welcome") + " " +
		locale.Message("Welcome1") + "\n\n" +
		locale.Message("WelcomePrompt") + "\n" +
		locale.Message("WelcomePromptPredictEnd") + "\n" +
		locale.Message("WelcomePromptPredictRemain") + "\n\n" +
		locale.Message("WelcomeMe") + "\n" +
		locale.Message("WelcomeCountry") + "\n" +
		locale.Message("WelcomeLanguage") + "\n" +
		locale.Message("WelcomeTrip") + "\n" +
		locale.Message("WelcomeContribute"))
}

func (f *TelegramFormatter) TripExplanation(language language.Tag) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.Message("TripExplanation") + "\n\n" +
		locale.Message("TripExplanationContinual") + "\n" +
		locale.Message("TripExplanationLimit") + "\n" +
		locale.Message("TripExplanationResetInterval"))
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

func (f *TelegramFormatter) markdownWrapper(str string) string {
	if !f.v2 {
		return str
	}
	return strings.ReplaceAll(tgapi.EscapeMarkdown(saveMarkdown(str)), "\\\\", "")
}

var saveableMarkdown = "_*`"

// saveMarkdown escapes allowed markdown before global markdown escape
func saveMarkdown(s string) string {
	var result []rune
	for _, r := range s {
		if strings.ContainsRune(saveableMarkdown, r) {
			result = append(result, '\\')
		}
		result = append(result, r)
	}
	return string(result)
}
