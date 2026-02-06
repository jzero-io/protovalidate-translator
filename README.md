# protovalidate-translator

Localized messages for [protovalidate](https://github.com/bufbuild/protovalidate) constraint IDs. Use with [go-i18n](https://github.com/nicksnyder/go-i18n) and template placeholders such as `{{.Value}}`.

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

msgEn, _ := translator.TranslateDefault("en", "float.lt", map[string]any{"Value": 100})
// msgEn == "value must be less than 100"
```

## Example

Run the example:

```bash
go run ./examples/translate/main.go
```

It prints sample translations for a few constraint IDs in English and Chinese.

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

## Supported languages

- **en** (default) – English  
- **zh** – 简体中文  

Message IDs follow the constraint IDs from `buf/validate` (e.g. `float.lt`, `string.min_len`, `int32.gt`). Add more languages by placing go-i18n JSON files in `translator/locales/` and rebuilding, or by loading your own bundle.

## Development

```bash
make test          # run tests (generates proto artifacts if needed)
make extract       # regenerate en.json from third_party/buf/validate/validate.proto
```

## License

[LICENSE](LICENSE)
