package validator

import (
	"errors"
    "fmt"
	"regexp"
)

// Rule 定义验证规则的类型
type Rule func(fieldName string, value interface{}) error

// StructValidator 负责管理对 struct 中字段的验证
type StructValidator struct {
	Fields map[string]*FieldValidator
}

// FieldValidator 用于验证单个字段
type FieldValidator struct {
	Value interface{}
	Rules []Rule
}

// AddField 添加字段及其验证规则
func (sv *StructValidator) AddField(fieldName string, value interface{}, rules ...Rule) {
	if sv.Fields == nil {
		sv.Fields = make(map[string]*FieldValidator)
	}
	sv.Fields[fieldName] = &FieldValidator{
		Value: value,
		Rules: rules,
	}
}

// Validate 执行所有字段的验证
func (sv *StructValidator) Validate() map[string]error {
	errorsMap := make(map[string]error)

	for fieldName, fieldValidator := range sv.Fields {
		for _, rule := range fieldValidator.Rules {
			if err := rule(fieldName, fieldValidator.Value); err != nil {
				errorsMap[fieldName] = err
				break // 遇到第一个错误即停止验证该字段
			}
		}
	}

	return errorsMap
}

// Required 验证字段不能为空
func Required() Rule {
	return func(fieldName string, value interface{}) error {
		if value == nil || value == "" {
			return errors.New(fieldName + " is required")
		}
		return nil
	}
}

// MinLength 验证字符串的最小长度
func MinLength(min int) Rule {
	return func(fieldName string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New(fieldName + " must be a string")
		}
		if len(str) < min {
			return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
		}
		return nil
	}
}

// MatchRegex 验证字符串是否符合正则表达式
func MatchRegex(pattern string) Rule {
	reg := regexp.MustCompile(pattern)
	return func(fieldName string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New(fieldName + " must be a string")
		}
		if !reg.MatchString(str) {
			return errors.New(fieldName + " format is invalid")
		}
		return nil
	}
}
