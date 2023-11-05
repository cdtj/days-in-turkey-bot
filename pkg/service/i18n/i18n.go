package i18n

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

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
		GetLocale(msg.Tag)
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
