# protovalidate-translator

**阅读语言：** [English](README.md) | 中文

为 [protovalidate](https://github.com/bufbuild/protovalidate) 的 rule ID 提供本地化文案，配合 [go-i18n](https://github.com/nicksnyder/go-i18n) 使用，支持 `{{.Value}}` 等模板占位符。

## 安装

```bash
go get github.com/jzero-io/protovalidate-translator/translator
```

## 快速开始

使用内置的嵌入文案（en + zh）：

```go
import "github.com/jzero-io/protovalidate-translator/translator"

msg, _ := translator.TranslateDefault("zh", "float.lt", map[string]any{"Value": 100})
// msg == "值必须小于 100"

msgTW,_ := translator.TranslateDefault("zh-TW", "float.lt", map[string]any{"Value": 100})
// msg == "值必須小於 100"

msgEn, _ := translator.TranslateDefault("en", "float.lt", map[string]any{"Value": 100})
// msgEn == "value must be less than 100"
```

## 示例

运行示例程序：

```bash
go run ./examples/translate/main.go
```

会打印若干 rule ID 的英文与中文翻译示例。

## 与校验错误一起使用

使用 protovalidate 校验后，用违规的 rule ID 和 value 获取本地化文案：

```go
err := validator.Validate(msg)
var valErr *protovalidate.ValidationError
if errors.As(err, &valErr) {
    for _, v := range valErr.Violations {
        ruleID := v.Proto.GetRuleId()
        data := map[string]any{"Value": v.RuleValue.Interface()}
        message, _ := translator.TranslateDefault("zh", ruleID, data)
        // message 即为本地化错误文案
    }
}
```

## 自定义文案

从目录加载自己的 go-i18n JSON 文案：

```go
bundle, err := translator.LoadBundleFromDir("./locales")
msg, _ := translator.Translate(bundle, "zh", "en", "float.lt", map[string]any{"Value": 100})
```

或从 embed FS 加载：

```go
//go:embed locales/*.json
var locales embed.FS

bundle, _ := translator.LoadBundleFromFS(locales, "locales")
msg, _ := translator.Translate(bundle, "zh", "en", "float.lt", data)
```

## 扩展默认文案包

可在默认文案包（供 `TranslateDefault` 使用）上增加语言或单条文案。**请在首次调用 `DefaultBundle` 或 `TranslateDefault` 之前注册。**

从磁盘增加语言文件：

```go
translator.AddDefaultLocaleFile("./locales/custom.json")
msg, _ := translator.TranslateDefault("zh", "float.lt", data)
```

从 `fs.FS`（如 embed）增加语言文件：

```go
//go:embed custom.json
var customFS embed.FS

translator.AddDefaultLocaleFromFS(customFS, "custom.json")
```

增加单条文案（覆盖或追加该语言）：

```go
translator.AddDefaultMessage("zh", "my.rule", "自訂：{{.Value}}")
msg, _ := translator.TranslateDefault("zh", "my.rule", map[string]any{"Value": 1})
```

需要完全控制时，可使用 customizer：

```go
translator.AddDefaultBundleCustomizer(func(b *i18n.Bundle) error {
    // 例如 b.LoadMessageFile(path), b.AddMessages(tag, msgs...)
    return nil
})
```

## 支持的语言

- **en**（默认）– 英文  
- **zh** – 简体中文  
- **zh-TW** – 繁體中文  

文案 ID 与 `buf/validate` 的 rule ID 一致（如 `float.lt`、`string.min_len`、`int32.gt`）。可在 `translator/locales/` 下放置 go-i18n JSON 并重新构建以增加语言，或自行加载 bundle。

## 开发说明

- **主库**（仓库根目录）：仅包含可供外部 import 的 `translator` 包，无测试、不依赖生成的 pb，直接执行 `go build ./...` 与 `go mod tidy` 即可通过。
- **测试与示例** 均在 `examples/` 下，且 **examples 使用独立 `go.mod`**，避免主库引用不存在的 pb 包。

```bash
make test              # 主库：仅构建
cd examples && make proto-go && go test ./... -v   # 生成 pb 并运行全部测试
# 或在仓库根目录执行：
make test-examples     # 同上
make extract           # 从 validate.proto 重新生成 en.json
```

在 `examples` 目录下执行 `go mod tidy` 和 `go test ./...` 即可。集成测试依赖 `examples/translate/testdata/pb`，需先在 examples 目录执行 `make proto-go`。

## 许可证

[LICENSE](LICENSE)
