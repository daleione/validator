package validator

import "fmt"

// ValidationError 统一错误结构
type ValidationError struct {
	Field   string // 出错的字段
	Message string // 错误消息
	Code    int    // 错误码（可选）
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Field '%s': %s", e.Field, e.Message)
}

// ErrorsCollector 用于收集和处理多个验证错误
type ErrorsCollector struct {
	Errors []ValidationError
}

// Add 添加一个新的验证错误
func (ec *ErrorsCollector) Add(field, message string, code int) {
	ec.Errors = append(ec.Errors, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// HasErrors 判断是否存在错误
func (ec *ErrorsCollector) HasErrors() bool {
	return len(ec.Errors) > 0
}

// Error 返回所有错误的汇总消息
func (ec *ErrorsCollector) Error() string {
	var errorMessages string
	for _, err := range ec.Errors {
		errorMessages += err.Error() + "\n"
	}
	return errorMessages
}

// String 返回格式化的错误信息字符串
func (ec *ErrorsCollector) String() string {
	if !ec.HasErrors() {
		return "No errors"
	}

	var result string
	for i, err := range ec.Errors {
		result += fmt.Sprintf("Error %d:\n", i+1)
		result += fmt.Sprintf("  Field: %s\n", err.Field)
		result += fmt.Sprintf("  Message: %s\n", err.Message)
		if err.Code != 0 {
			result += fmt.Sprintf("  Code: %d\n", err.Code)
		}
		result += "\n"
	}
	return result
}
