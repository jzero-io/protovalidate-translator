count:
	grep -v '^\s*//' third_party/buf/validate/validate.proto | grep -o 'id:\s*"[^"]*"' | wc -l

extract:
	python3 scripts/extract_validate_messages.py

# 主库无测试，仅确保构建通过
test:
	go build ./...

# 在 examples 目录运行全部测试（单元 + 依赖 pb 的集成测试）
test-examples:
	cd examples && make proto-go && go test ./... -v
