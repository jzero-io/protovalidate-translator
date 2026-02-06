package translator

import (
	"path/filepath"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func newTestBundle() *i18n.Bundle {
	bundle := NewBundle()
	bundle.AddMessages(language.English,
		&i18n.Message{ID: "float.lt", Other: "value must be less than {{.Value}}"},
		&i18n.Message{ID: "float.finite", Other: "value must be finite"},
	)
	bundle.AddMessages(language.Chinese,
		&i18n.Message{ID: "float.lt", Other: "值必须小于 {{.Value}}"},
		&i18n.Message{ID: "string.min_len", Other: "长度至少为 {{.Value}}"},
	)
	return bundle
}

func TestTranslate_byLang(t *testing.T) {
	bundle := newTestBundle()
	en, err := Translate(bundle, "en", "", "float.lt", map[string]any{"Value": 100})
	if err != nil {
		t.Fatal(err)
	}
	if en != "value must be less than 100" {
		t.Errorf("en: got %q", en)
	}
	zh, err := Translate(bundle, "zh", "", "float.lt", map[string]any{"Value": 100})
	if err != nil {
		t.Fatal(err)
	}
	if zh != "值必须小于 100" {
		t.Errorf("zh: got %q", zh)
	}
}

func TestTranslate_fallbackToDefaultLang(t *testing.T) {
	bundle := newTestBundle()
	out, err := Translate(bundle, "fr", "en", "float.lt", map[string]any{"Value": 1})
	if err != nil {
		t.Fatal(err)
	}
	if out != "value must be less than 1" {
		t.Errorf("fallback to en: got %q", out)
	}
}

func TestTranslate_fallbackIdInDefaultOnly(t *testing.T) {
	bundle := newTestBundle()
	out, err := Translate(bundle, "zh", "en", "float.lt", map[string]any{"Value": 10})
	if err != nil {
		t.Fatal(err)
	}
	if out != "值必须小于 10" {
		t.Errorf("zh float.lt: got %q", out)
	}
	out2, err := Translate(bundle, "zh", "en", "float.finite", nil)
	if err != nil {
		t.Fatal(err)
	}
	if out2 != "value must be finite" {
		t.Errorf("fallback to en: got %q", out2)
	}
}

func TestTranslate_missingIdReturnsId(t *testing.T) {
	bundle := newTestBundle()
	out, err := Translate(bundle, "en", "", "nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != "nonexistent" {
		t.Errorf("expected id as fallback, got %q", out)
	}
}

func TestLoadBundleFromDir(t *testing.T) {
	bundle, err := LoadBundleFromDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	out, err := Translate(bundle, "sample", "", "float.lt", map[string]any{"Value": 42})
	if err != nil {
		t.Fatal(err)
	}
	if out != "value must be less than 42" {
		t.Errorf("got %q", out)
	}
	outZh, err := Translate(bundle, "zh", "", "float.lt", map[string]any{"Value": 42})
	if err != nil {
		t.Fatal(err)
	}
	if outZh != "值必须小于 42" {
		t.Errorf("zh: got %q", outZh)
	}
}

func TestLoadBundleFromFS(t *testing.T) {
	bundle, err := LoadBundleFromFS(LocalesFS, DefaultLocaleDir)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Translate(bundle, "en", "", "float.finite", nil)
	if err != nil {
		t.Fatal(err)
	}
	if out == "float.finite" {
		t.Errorf("expected translation from embedded locales")
	}
}

func TestDefaultBundle(t *testing.T) {
	bundle, err := DefaultBundle()
	if err != nil {
		t.Fatal(err)
	}
	if bundle == nil {
		t.Fatal("DefaultBundle returned nil")
	}
}

func TestTranslateDefault(t *testing.T) {
	out, err := TranslateDefault("en", "float.lt", map[string]any{"Value": 5})
	if err != nil {
		t.Fatal(err)
	}
	if out == "float.lt" {
		t.Errorf("expected translation from default bundle")
	}
}

func TestLoadBundleFromDir_missing(t *testing.T) {
	_, err := LoadBundleFromDir(filepath.Join("testdata", "missing"))
	if err == nil {
		t.Fatal("expected error for missing dir")
	}
}
