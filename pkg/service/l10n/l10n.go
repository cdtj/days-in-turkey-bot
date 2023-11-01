package l10n

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"cdtj.io/days-in-turkey-bot/assets"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry
const l10nDir = "l10n"

var localizers sync.Map

var (
	ErrNoDefaultPkg = errors.New("no package for default locale, can't run the app")
	ErrNoFiles      = errors.New("no localization files in l10n folder")
	ErrUknownLang   = errors.New("unknown language")
	ErrUknownMsg    = errors.New("unknown message")
)

var (
	defaultLang = language.English
)

func Localization() error {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	fs, err := assets.L10n.ReadDir(l10nDir)
	if err != nil {
		return fmt.Errorf("unable to read l10n dir: %w", err)
	}
	if len(fs) == 0 {
		return ErrNoFiles
	}
	for _, f := range fs {
		msg, err := bundle.LoadMessageFileFS(assets.L10n, filepath.Join(l10nDir, f.Name()))
		if err != nil {
			return fmt.Errorf("unable to read %s: %w", f.Name(), err)
		}
		localizers.Store(msg.Tag, i18n.NewLocalizer(bundle, msg.Tag.String()))
	}
	// to avoid nil check on every method call, we make sure that we have default locale
	if !ValidateLang(defaultLang) {
		return ErrNoDefaultPkg
	}
	return nil
}

func DefaultLang() language.Tag {
	return defaultLang
}

func localizer(tag language.Tag) *i18n.Localizer {
	l, ok := localizers.Load(tag)
	if !ok {
		return nil
	}
	return (l).(*i18n.Localizer)
}

func ValidateLang(tag language.Tag) bool {
	_, ok := localizers.Load(tag)
	return ok
}

type Locale struct {
	tag       language.Tag
	localizer *i18n.Localizer
}

func NewLocale(tag language.Tag) *Locale {
	// oh my gd what im doing here?!
	for _, t := range []language.Tag{tag, defaultLang} {
		l := localizer(t)
		if l != nil {
			return &Locale{
				tag:       tag,
				localizer: l,
			}
		}
	}
	return nil
}

func (l *Locale) Lang() string {
	return l.tag.String()
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
