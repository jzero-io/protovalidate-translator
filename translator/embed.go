package translator

import (
	"embed"
	"io/fs"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	// DefaultLocaleDir is the embedded locales directory path.
	DefaultLocaleDir = "locales"
	// DefaultLang is the fallback language used by helpers in this package.
	DefaultLang = "en"
)

// LocalesFS embeds the default locale files.
//
//go:embed locales/*.json
var LocalesFS embed.FS

// BundleCustomizer is a function that extends a bundle (e.g. load more files or add messages).
// It is applied when building the default bundle, after the embedded locales are loaded.
type BundleCustomizer func(b *i18n.Bundle) error

var (
	defaultBundleOnce    sync.Once
	defaultBundle        *i18n.Bundle
	defaultBundleErr     error
	defaultCustomizersMu sync.Mutex
	defaultCustomizers   []BundleCustomizer
)

// AddDefaultBundleCustomizer registers a customizer that will run when the default bundle
// is first built (on first call to DefaultBundle). Register before any call to DefaultBundle.
func AddDefaultBundleCustomizer(fn BundleCustomizer) {
	defaultCustomizersMu.Lock()
	defer defaultCustomizersMu.Unlock()
	defaultCustomizers = append(defaultCustomizers, fn)
}

// AddDefaultLocaleFile registers a locale file to load into the default bundle.
// The path is passed to bundle.LoadMessageFile (e.g. "locales/custom.json").
// Register before first use of DefaultBundle.
func AddDefaultLocaleFile(path string) {
	AddDefaultBundleCustomizer(func(b *i18n.Bundle) error {
		_, err := b.LoadMessageFile(path)
		return err
	})
}

// AddDefaultLocaleFromFS registers a locale file from an fs.FS (e.g. embed.FS) to load
// into the default bundle. Register before first use of DefaultBundle.
func AddDefaultLocaleFromFS(fsys fs.FS, path string) {
	AddDefaultBundleCustomizer(func(b *i18n.Bundle) error {
		_, err := b.LoadMessageFileFS(fsys, path)
		return err
	})
}

// AddDefaultMessage registers a single message for the given language.
// lang is a BCP 47 tag (e.g. "en", "zh", "zh-TW"). It overrides or adds to the default bundle.
// Invalid lang is ignored. Register before first use of DefaultBundle.
func AddDefaultMessage(lang string, id string, template string) {
	tag, err := language.Parse(lang)
	if err != nil {
		return
	}
	msg := &i18n.Message{ID: id, Other: template}
	AddDefaultBundleCustomizer(func(b *i18n.Bundle) error {
		b.AddMessages(tag, msg)
		return nil
	})
}

// DefaultBundle returns a cached bundle: first loads embedded locales (locales/*.json),
// then runs all customizers registered via AddDefaultLocaleFile, AddDefaultLocaleFromFS,
// and AddDefaultMessage. Safe for concurrent use.
func DefaultBundle() (*i18n.Bundle, error) {
	defaultBundleOnce.Do(func() {
		defaultBundle, defaultBundleErr = LoadBundleFromFS(LocalesFS, DefaultLocaleDir)
		if defaultBundleErr != nil {
			return
		}
		defaultCustomizersMu.Lock()
		customizers := append([]BundleCustomizer(nil), defaultCustomizers...)
		defaultCustomizersMu.Unlock()
		for _, fn := range customizers {
			if err := fn(defaultBundle); err != nil {
				defaultBundleErr = err
				defaultBundle = nil
				return
			}
		}
	})
	return defaultBundle, defaultBundleErr
}
