package validator

import (
	"fmt"
    "testing"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func TestStruct(t *testing.T) {
	// 创建一个用户实例
	user := User{
		Name:  "",
		Email: "invalid_email",
		Age:   25,
	}

	// 创建 StructValidator 实例
	structValidator := &StructValidator{}

	// 为每个字段添加验证规则
	structValidator.AddField("Name", user.Name, Required(), MinLength(3))
	structValidator.AddField("Email", user.Email, Required(), MatchRegex("^[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,}$"))
	structValidator.AddField("Age", user.Age, Required())

	// 执行验证
	errors := structValidator.Validate()

	// 打印验证错误
	for field, err := range errors {
		fmt.Printf("Error in %s: %s\n", field, err)
	}
}
