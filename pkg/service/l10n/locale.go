package l10n

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var locales sync.Map

type Locale struct {
	Name      string
	Tag       language.Tag
	localizer *i18n.Localizer
}

func GetLocale(tag language.Tag) *Locale {
	if locale, ok := locales.Load(tag); ok {
		return locale.(*Locale)
	}
	for _, t := range []language.Tag{tag, defaultLang} {
		l := localizer(t)
		slog.Info("localizer", "tag", t, "l", l)
		if l != nil {
			locale := &Locale{
				Tag:       t,
				localizer: l,
			}
			locale.Name = locale.Message("Name")
			locales.Store(tag, locale)
			return locale
		}
	}
	slog.Error("empty locale", "tag", tag)
	return nil
}

func Locales() []*Locale {
	arr := make([]*Locale, 0)
	slog.Info("getting locales")
	locales.Range(func(key, value interface{}) bool {
		slog.Info("appending", "locale", value)
		arr = append(arr, value.(*Locale))
		return true
	})
	return arr
}

func (l *Locale) Lang() string {
	return l.Name
}

func (l *Locale) Message(messageID string) string {
	msg, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return errors.Join(ErrUknownMsg, err).Error()
	}
	return msg
}

func (l *Locale) Error(messageID string, err error) string {
	return l.MessageWithTemplate(messageID, map[string]interface{}{"Error": err}, nil)
}

func (l *Locale) MessageWithCount(messageID string, count interface{}) string {
	return l.MessageWithTemplate(messageID, map[string]interface{}{"Count": count}, count)
}

func (l *Locale) MessageWithTemplate(messageID string, tpl map[string]interface{}, plural interface{}) string {
	msg, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: tpl,
		PluralCount:  plural,
	})
	if err != nil {
		return errors.Join(ErrUknownMsg, err).Error()
	}
	return msg
}

const dateLayout = "02/01/2006"

// TODO: add regional formats
func (l *Locale) FormatDate(dt time.Time) string {
	return dt.Format(dateLayout)
}
