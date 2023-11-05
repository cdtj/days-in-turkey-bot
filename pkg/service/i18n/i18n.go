package i18n

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"

	"cdtj.io/days-in-turkey-bot/assets"
	"github.com/BurntSushi/toml"
	i18nlib "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18ner interface {
	GetLocale(tag language.Tag) Localizer
}

type I18n struct {
	dir         string
	defaultLang language.Tag
	localizers  sync.Map
	locales     sync.Map
}

var (
	ErrNoDefaultPkg = errors.New("no package for default locale, can't run the app")
	ErrNoFiles      = errors.New("no localization files in i18n folder")
	ErrUknownLang   = errors.New("unknown language")
	ErrUknownMsg    = errors.New("unknown message")
)

func NewI18n(dir string, defaultLang string) (*I18n, error) {
	tag, err := language.Parse(defaultLang)
	if err != nil {
		return nil, err
	}
	bundle := i18nlib.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	fs, err := assets.I18n.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read [%s]: %w", dir, err)
	}
	if len(fs) == 0 {
		return nil, ErrNoFiles
	}
	i18n := &I18n{
		dir:         dir,
		defaultLang: tag,
	}
	for _, f := range fs {
		msg, err := bundle.LoadMessageFileFS(assets.I18n, filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("unable to read %s: %w", f.Name(), err)
		}
		i18n.localizers.Store(msg.Tag, i18nlib.NewLocalizer(bundle, msg.Tag.String()))
	}
	// to avoid nil, we make sure that we have default locale
	if !i18n.ValidateLang(tag) {
		return nil, ErrNoDefaultPkg
	}
	return i18n, nil
}

func (i *I18n) DefaultLang() language.Tag {
	return i.defaultLang
}

func (i *I18n) GetLocale(tag language.Tag) Localizer {
	if locale, ok := i.locales.Load(tag); ok {
		return locale.(*Locale)
	}
	for _, t := range []language.Tag{tag, i.defaultLang} {
		l := i.localizer(t)
		slog.Debug("localizer", "tag", t, "l", l)
		if l != nil {
			locale := &Locale{
				Tag:       t,
				localizer: l,
			}
			locale.Name = locale.Message("LanguageName")
			i.locales.Store(tag, locale)
			return locale
		}
	}
	slog.Error("empty locale", "tag", tag)
	return nil
}

func (i *I18n) GetLocaleByString(lang string) Localizer {
	tag, err := language.Parse(lang)
	if err != nil {
		return i.GetLocale(i.defaultLang)
	}
	return i.GetLocale(tag)
}

func (i *I18n) Locales() []*Locale {
	arr := make([]*Locale, 0)
	i.locales.Range(func(key, value interface{}) bool {
		arr = append(arr, value.(*Locale))
		return true
	})
	return arr
}

func (i *I18n) ValidateLang(tag language.Tag) bool {
	return i.localizer(tag) != nil
}

func (i *I18n) localizer(tag language.Tag) *i18nlib.Localizer {
	l, ok := i.localizers.Load(tag)
	if !ok {
		return nil
	}
	return (l).(*i18nlib.Localizer)
}
