package i18n

import (
	"errors"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer interface {
	Message(messageID string) string
	MessageWithCount(messageID string, count interface{}) string
	MessageWithTemplate(messageID string, tpl map[string]interface{}, plural interface{}) string
	Error(messageID string, err error) string
	FormatDate(dt time.Time) string
}

type Locale struct {
	Name      string
	Tag       language.Tag
	localizer *i18n.Localizer
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

func (l *Locale) Error(messageID string, err error) string {
	return l.MessageWithTemplate(messageID, map[string]interface{}{"Error": err}, nil)
}

const dateLayout = "02/01/2006"

// TODO: add regional formats
func (l *Locale) FormatDate(dt time.Time) string {
	return dt.Format(dateLayout)
}
