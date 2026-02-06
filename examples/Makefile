# 生成 Go pb 到 translate/testdata/pb/（需在 examples 目录执行，依赖上级 third_party）
# 需安装: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
proto-go:
	@which protoc-gen-go >/dev/null 2>&1 || { echo "need protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
	@mkdir -p translate/testdata/pb
	protoc -I . -I third_party -I translate/testdata --go_out=.. --go_opt=module=github.com/jzero-io/protovalidate-translator \
		translate/testdata/proto/user.proto \
		translate/testdata/proto/order.proto

test:
	go test ./ -v
