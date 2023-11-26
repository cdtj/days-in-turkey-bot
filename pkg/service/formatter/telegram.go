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
	msg := ""
	firstLine := true
	firstEligible := true
	overstayed := false
	locale := f.i18n.GetLocale(language)
	for i := tree; i != nil; i = i.Prev {
		slog.Debug("TripTree", "StartDate", i.StartDate, "EndDate", i.EndDate, "TripDays", i.TripDays, "PeriodDays", i.PeriodDays, "OverstayDays", i.OverstayDays)
		if i.StartPredicted || i.EndPredicted {
			if i.StartPredicted && firstEligible {
				msg += locale.Message("TripEligibleHdr") + "\n"
				firstEligible = false
			}
			msg += locale.MessageWithTemplate("TripPredicted", map[string]interface{}{
				"StartDate":  wrapCode(locale.FormatDate(i.StartDate)),
				"EndDate":    wrapCode(locale.FormatDate(i.EndDate)),
				"TripDays":   locale.MessageWithCount("DayCounter", i.TripDays),
				"PeriodDays": locale.MessageWithCount("DayCounter", i.PeriodDays),
			}, nil)
			if i.OverstayDays > 0 {
				overstayed = true
				msg += ", " + locale.MessageWithTemplate("Overstay", map[string]interface{}{
					"OverstayDays": wrapBold(locale.MessageWithCount("DayCounter", i.OverstayDays)),
				}, nil) + " ⚠️"
			}
			msg += "\n"
		} else {
			if firstLine {
				msg += "\n" + locale.Message("TripPast") + "\n"
				firstLine = false
			}
			msg += locale.MessageWithTemplate("Trip", map[string]interface{}{
				"StartDate":  wrapCode(locale.FormatDate(i.StartDate)),
				"EndDate":    wrapCode(locale.FormatDate(i.EndDate)),
				"TripDays":   locale.MessageWithCount("DayCounter", i.TripDays),
				"PeriodDays": locale.MessageWithCount("DayCounter", i.PeriodDays),
			}, nil)
			if i.OverstayDays > 0 {
				overstayed = true
				msg += ", " + locale.MessageWithTemplate("Overstay", map[string]interface{}{
					"OverstayDays": wrapBold(locale.MessageWithCount("DayCounter", i.OverstayDays)),
				}, nil) + " ⚠️"
			}
			msg += "\n"
		}
	}
	msg += "\n" + locale.Message("OverstayCaution")
	if overstayed {
		msg += "\n" + locale.Message("OverstayExplanation")
	}
	return f.markdownWrapper(msg)
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
	msg := locale.MessageWithTemplate("CountryInfo", map[string]interface{}{
		"Flag": country.GetFlag(),
		"Name": country.GetName(),
	}, nil) +
		"\n" +
		locale.MessageWithTemplate("CountryDays", map[string]interface{}{
			"Continual":     country.GetDaysCont(),
			"Limit":         country.GetDaysLimit(),
			"ResetInterval": country.GetResetInterval(),
		}, nil)
	if !country.GetVisaFree() {
		msg += "\n\n⚠️ " + locale.Message("CountryVisaWarning") + " ⚠️"
	}
	return f.markdownWrapper(msg)
}

func (f *TelegramFormatter) FormatMessage(language language.Tag, messageID string) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.Message(messageID))
}

func (f *TelegramFormatter) FormatError(language language.Tag, err error) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.ErrorWithDefault(err, "ErrorInternal"))
}

func (f *TelegramFormatter) Welcome(language language.Tag) string {
	locale := f.i18n.GetLocale(language)
	return f.markdownWrapper(locale.Message("Welcome") + " " +
		locale.Message("Welcome1") + "\n\n" +
		locale.Message("WelcomePrompt") + "\n" +
		locale.MessageWithTemplate("WelcomePromptPredictEnd", map[string]interface{}{
			"SingleDate": wrapCode("31/12/2022"),
		}, nil) + "\n" +
		locale.MessageWithTemplate("WelcomePromptPredictRemain", map[string]interface{}{
			"MultiDate": wrapCode("31/12/2022 15/01/2023 01/02/2023 15/02/2023"),
		}, nil) + "\n\n" +
		// "/me - " + locale.Message("CommandMe") + "\n" +
		"/country - " + locale.Message("CommandCountry") + "\n" +
		"/language - " + locale.Message("CommandLanguage") + "\n" +
		"/trip - " + locale.Message("CommandTrip") + "\n" +
		// "/contribute - " + locale.Message("CommandContribute") + "\n" +
		// "/feedback - " + locale.Message("CommandFeedback") +
		"")
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
