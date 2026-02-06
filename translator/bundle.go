package translator

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// NewBundle creates a go-i18n bundle.
func NewBundle() *i18n.Bundle {
	return i18n.NewBundle(language.Und)
}

// LoadBundleFromFS loads locale JSON files from an fs.FS (e.g., embed.FS).
// The directory should contain files like en.json, zh.json, etc.
func LoadBundleFromFS(fsys fs.FS, dir string) (*i18n.Bundle, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	bundle := NewBundle()
	loaded := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		filePath := path.Join(dir, name)
		if _, err := bundle.LoadMessageFileFS(fsys, filePath); err != nil {
			return nil, err
		}
		loaded++
	}
	if loaded == 0 {
		return nil, fmt.Errorf("no locale files found in %s", dir)
	}
	return bundle, nil
}

// LoadBundleFromDir loads locale JSON files from an OS directory.
func LoadBundleFromDir(dir string) (*i18n.Bundle, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	bundle := NewBundle()
	loaded := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		filePath := filepath.Join(dir, name)
		if _, err := bundle.LoadMessageFile(filePath); err != nil {
			return nil, err
		}
		loaded++
	}
	if loaded == 0 {
		return nil, fmt.Errorf("no locale files found in %s", dir)
	}
	return bundle, nil
}
