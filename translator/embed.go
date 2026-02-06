package translator

import (
	"embed"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
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

var (
	defaultBundleOnce sync.Once
	defaultBundle     *i18n.Bundle
	defaultBundleErr  error
)

// DefaultBundle returns a cached bundle loaded from the embedded locales.
func DefaultBundle() (*i18n.Bundle, error) {
	defaultBundleOnce.Do(func() {
		defaultBundle, defaultBundleErr = LoadBundleFromFS(LocalesFS, DefaultLocaleDir)
	})
	return defaultBundle, defaultBundleErr
}
