package i18n

import (
	"errors"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer interface {
	GetName() string
	GetLanguage() string

	Message(messageID string) string
	MessageWithCount(messageID string, count interface{}) string
	MessageWithTemplate(messageID string, tpl map[string]interface{}, plural interface{}) string

	Error(err *model.LError) string
	ErrorWithDefault(err error, defaultMessageID string) string

	FormatDate(dt time.Time) string
}

type Locale struct {
	Name      string
	Tag       language.Tag
	localizer *i18n.Localizer
}

func (l *Locale) GetName() string {
	return l.Name
}

func (l *Locale) GetLanguage() string {
	return l.Tag.String()
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

func (l *Locale) ErrorWithDefault(err error, defaultMessageID string) string {
	switch lerr := err.(type) {
	case *model.LError:
		return l.Error(lerr)
	}
	return l.Message(defaultMessageID)
}

func (l *Locale) Error(err *model.LError) string {
	if err.Template != nil {
		for key := range err.Template {
			switch val := err.Template[key].(type) {
			case model.LErrorExpandable:
				err.Template[key] = l.Message(string(val))
			}
		}
	}
	return l.MessageWithTemplate(err.Code, err.Template, nil)
}

const dateLayout = "02/01/2006"

// TODO: add regional formats
func (l *Locale) FormatDate(dt time.Time) string {
	return dt.Format(dateLayout)
}
