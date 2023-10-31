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

const l10nDir = "l10n"

var localizers sync.Map

var (
	ErrNoFiles    = errors.New("no localization files in l10n folder")
	ErrUknownLang = errors.New("unknown language")
	ErrUknownMsg  = errors.New("unknown message")
)

func Localization() error {
	bundle := i18n.NewBundle(language.English)
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
	return nil
}

func localizer(tag language.Tag) *i18n.Localizer {
	l, ok := localizers.Load(tag)
	if !ok {
		return nil
	}
	return (l).(*i18n.Localizer)
}

func Localaze(tag language.Tag, messageID string) string {
	return LocalazeWithCount(tag, messageID, nil)
}

func LocalazeWithCount(tag language.Tag, messageID string, count interface{}) string {
	l := localizer(tag)
	if l != nil {
		msg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: messageID,
			TemplateData: map[string]interface{}{
				"Count": count,
			},
			PluralCount: count,
		})
		if err != nil {
			fmt.Println("localize failed:", err)
			return ErrUknownMsg.Error()
		}
		return msg
	}
	fmt.Println("unknown lang:", tag)
	return ErrUknownLang.Error()
}

func ValidateLang(tag language.Tag) bool {
	_, ok := localizers.Load(tag)
	return ok
}

const dateLayout = "02/01/2006"

// TODO: add regional formats
func FormatDate(dt time.Time) string {
	return dt.Format(dateLayout)
}
