//go:build integration

package translator

import (
	"errors"
	"testing"

	protovalidate "buf.build/go/protovalidate"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/jzero-io/protovalidate-translator/translator/testdata/pb"
)

// TestValidateProto_GeneratedGo 使用 protoc 生成的 Go pb 类型进行校验与翻译测试。
// 需先执行 make proto-go 生成 translator/testdata/pb/*.pb.go。
func TestValidateProto_GeneratedGo(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithMessages(&pb.User{}, &pb.Order{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	user := &pb.User{
		Email: "not-an-email",
		Age:   10,
		Name:  "a",
	}
	if err := validator.Validate(user); err == nil {
		t.Fatal("expected validation error for User")
	} else {
		assertZhTranslations(t, err)
	}

	order := &pb.Order{
		Id:       "a",
		Quantity: 0,
	}
	if err := validator.Validate(order); err == nil {
		t.Fatal("expected validation error for Order")
	} else {
		assertZhTranslations(t, err)
	}
}

func TestValidateProto_TranslateZh(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithMessages(&pb.User{}, &pb.Order{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	user := &pb.User{}
	if err := validator.Validate(user); err == nil {
		t.Fatal("expected validation error for User")
	} else {
		assertZhTranslations(t, err)
	}
}

func assertZhTranslations(t *testing.T, err error) {
	t.Helper()
	var valErr *protovalidate.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if len(valErr.Violations) == 0 {
		t.Fatal("expected violations")
	}
	for _, violation := range valErr.Violations {
		ruleID := violation.Proto.GetRuleId()
		data := map[string]any{"Value": normalizeRuleValue(violation.RuleValue)}
		zh, err := TranslateDefault("zh", ruleID, data)
		if err != nil {
			t.Fatalf("translate %s: %v", ruleID, err)
		}
		if zh == ruleID {
			t.Fatalf("missing zh translation for %s", ruleID)
		}
		// 获取违反规则的字段路径（如 "email"、"name"；嵌套时为 "user.address"）
		fieldPath := ""
		if fp := violation.Proto.GetField(); fp != nil {
			fieldPath = protovalidate.FieldPathString(fp)
		}
		t.Logf("translated %s: %s for field %s", ruleID, zh, fieldPath)
	}
}

func normalizeRuleValue(value any) any {
	switch v := value.(type) {
	case protoreflect.Value:
		return v.Interface()
	default:
		return value
	}
}
