package translator

import (
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Translate renders a message by id using the specified language and optional fallback language.
// If the message is missing in both languages, the id itself is returned.
func Translate(bundle *i18n.Bundle, lang string, defaultLang string, id string, data map[string]any) (string, error) {
	if bundle == nil {
		return id, nil
	}
	if msg, ok, err := localize(bundle, lang, id, data); err != nil {
		return "", err
	} else if ok {
		return msg, nil
	}
	if defaultLang != "" && defaultLang != lang {
		if msg, ok, err := localize(bundle, defaultLang, id, data); err != nil {
			return "", err
		} else if ok {
			return msg, nil
		}
	}
	return id, nil
}

// MustTranslate is like Translate but panics on error.
func MustTranslate(bundle *i18n.Bundle, lang string, defaultLang string, id string, data map[string]any) string {
	msg, err := Translate(bundle, lang, defaultLang, id, data)
	if err != nil {
		panic(err)
	}
	return msg
}

// TranslateDefault uses the embedded locales and DefaultLang as fallback.
func TranslateDefault(lang string, id string, data map[string]any) (string, error) {
	bundle, err := DefaultBundle()
	if err != nil {
		return "", err
	}
	return Translate(bundle, lang, DefaultLang, id, data)
}

// MustTranslateDefault is like TranslateDefault but panics on error.
func MustTranslateDefault(lang string, id string, data map[string]any) string {
	msg, err := TranslateDefault(lang, id, data)
	if err != nil {
		panic(err)
	}
	return msg
}

func localize(bundle *i18n.Bundle, lang string, id string, data map[string]any) (string, bool, error) {
	if lang == "" {
		return "", false, nil
	}
	localizer := i18n.NewLocalizer(bundle, lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		var notFound *i18n.MessageNotFoundErr
		if errors.As(err, &notFound) {
			return "", false, nil
		}
		return "", false, err
	}
	return msg, true, nil
}
