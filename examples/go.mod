module github.com/jzero-io/protovalidate-translator/examples

go 1.24.3

replace github.com/jzero-io/protovalidate-translator => ../

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.10-20251209175733-2a1774d88802.1
	buf.build/go/protovalidate v1.1.0
	github.com/jzero-io/protovalidate-translator v0.0.0
	github.com/nicksnyder/go-i18n/v2 v2.6.1
	golang.org/x/text v0.32.0
	google.golang.org/protobuf v1.36.10
)

require (
	cel.dev/expr v0.24.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/google/cel-go v0.26.1 // indirect
	github.com/stoewer/go-strcase v1.3.1 // indirect
	golang.org/x/exp v0.0.0-20250813145105-42675adae3e6 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250811230008-5f3141c8851a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250811230008-5f3141c8851a // indirect
)
