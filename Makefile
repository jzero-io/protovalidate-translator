# 生成 Go pb 文件到 translator/testdata/pb/（需安装 protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest）
proto-go:
	@which protoc-gen-go >/dev/null 2>&1 || { echo "need protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
	@mkdir -p translator/testdata/pb
	protoc -I third_party -I. --go_out=. --go_opt=module=github.com/jzero-io/protovalidate-translator \
		translator/testdata/proto/user.proto \
		translator/testdata/proto/order.proto

count:
	grep -v '^\s*//' third_party/buf/validate/validate.proto | grep -o 'id:\s*"[^"]*"' | wc -l

extract:
	python3 scripts/extract_validate_messages.py

test: proto-go
	go test ./translator/... -v
