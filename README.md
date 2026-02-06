# protovalidate-translator

Localized messages for [protovalidate](https://github.com/bufbuild/protovalidate) rule IDs. Use with [go-i18n](https://github.com/nicksnyder/go-i18n) and template placeholders such as `{{.Value}}`.

## Install

```bash
go get github.com/jzero-io/protovalidate-translator/translator
```

## Quick start

Use the built-in embedded locales (en + zh):

```go
import "github.com/jzero-io/protovalidate-translator/translator"

msg, _ := translator.TranslateDefault("zh", "float.lt", map[string]any{"Value": 100})
// msg == "值必须小于 100"

msgTW,_ := translator.TranslateDefault("zh-TW", "float.lt", map[string]any{"Value": 100})
// msg == "值必須小於 100"

msgEn, _ := translator.TranslateDefault("en", "float.lt", map[string]any{"Value": 100})
// msgEn == "value must be less than 100"
```

## Example

Run the example:

```bash
go run ./examples/translate/main.go
```

It prints sample translations for a few rule IDs in English and Chinese.

## Usage with validation errors

After validating with protovalidate, use the violation’s rule ID and value to get a localized message:

```go
err := validator.Validate(msg)
var valErr *protovalidate.ValidationError
if errors.As(err, &valErr) {
    for _, v := range valErr.Violations {
        ruleID := v.Proto.GetRuleId()
        data := map[string]any{"Value": v.RuleValue.Interface()}
        message, _ := translator.TranslateDefault("zh", ruleID, data)
        // message is the localized error string
    }
}
```

## Custom locales

Load your own locale directory (go-i18n JSON):

```go
bundle, err := translator.LoadBundleFromDir("./locales")
msg, _ := translator.Translate(bundle, "zh", "en", "float.lt", map[string]any{"Value": 100})
```

Or from an embedded FS:

```go
//go:embed locales/*.json
var locales embed.FS

bundle, _ := translator.LoadBundleFromFS(locales, "locales")
msg, _ := translator.Translate(bundle, "zh", "en", "float.lt", data)
```

## Extending the default bundle

You can add locales or single messages to the default bundle (used by `TranslateDefault`). **Register before the first call to `DefaultBundle` or `TranslateDefault`.**

Add a locale file from disk:

```go
translator.AddDefaultLocaleFile("./locales/custom.json")
msg, _ := translator.TranslateDefault("zh", "float.lt", data)
```

Add a locale file from an `fs.FS` (e.g. embed):

```go
//go:embed custom.json
var customFS embed.FS

translator.AddDefaultLocaleFromFS(customFS, "custom.json")
```

Add a single message (overrides or adds for that language):

```go
translator.AddDefaultMessage("zh", "my.rule", "自訂：{{.Value}}")
msg, _ := translator.TranslateDefault("zh", "my.rule", map[string]any{"Value": 1})
```

For full control, use a customizer:

```go
translator.AddDefaultBundleCustomizer(func(b *i18n.Bundle) error {
    // e.g. b.LoadMessageFile(path), b.AddMessages(tag, msgs...)
    return nil
})
```

## Supported languages

- **en** (default) – English  
- **zh** – 简体中文  
- **zh-TW** – 繁體中文  

Message IDs follow the rule IDs from `buf/validate` (e.g. `float.lt`, `string.min_len`, `int32.gt`). Add more languages by placing go-i18n JSON files in `translator/locales/` and rebuilding, or by loading your own bundle.

## Development

```bash
make test-unit     # unit tests only (no generated pb; safe after clone)
make test          # generate testdata/pb + run all tests including protovalidate integration
make extract       # regenerate en.json from third_party/buf/validate/validate.proto
```

The protovalidate integration tests live in a file built only with `-tags=integration`, so `go build` and `go test ./translator/...` without the tag do not require the (uncommitted) generated `translator/testdata/pb` package. This keeps `go mod tidy` and CI green for consumers who do not run `make proto-go`.

## License

[LICENSE](LICENSE)
