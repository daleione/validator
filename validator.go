package validator

import (
	"errors"
	"fmt"
	"regexp"
)

// Rule 定义验证规则的类型
type Rule func(fieldName string, value any, ec *ErrorsCollector) error

// StructValidator 负责管理对 struct 中字段的验证
type StructValidator struct {
	Fields         map[string]*FieldValidator
	StopOnFirstErr bool // 控制是否只返回第一个错误
}

// FieldValidator 用于验证单个字段
type FieldValidator struct {
	Value any
	Rules []Rule
}

// AddField 添加字段及其验证规则
func (sv *StructValidator) AddField(fieldName string, value any, rules ...Rule) {
	if sv.Fields == nil {
		sv.Fields = make(map[string]*FieldValidator)
	}
	sv.Fields[fieldName] = &FieldValidator{
		Value: value,
		Rules: rules,
	}
}

// AddFieldGroup 添加一组字段及其统一的验证规则
func (sv *StructValidator) AddFieldGroup(fieldNames []string, rules ...Rule) {
	for _, fieldName := range fieldNames {
		if field, exists := sv.Fields[fieldName]; exists {
			field.Rules = append(field.Rules, rules...)
		}
	}
}

// Validate 执行所有字段的验证，并返回错误收集器
func (sv *StructValidator) Validate() *ErrorsCollector {
	ec := &ErrorsCollector{}

	for fieldName, fieldValidator := range sv.Fields {
		for _, rule := range fieldValidator.Rules {
			if err := rule(fieldName, fieldValidator.Value, ec); err != nil {
				if sv.StopOnFirstErr {
					return ec
				}
			}
		}
	}

	return ec
}

// Required 验证字段不能为空
func Required() Rule {
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		if value == nil || value == "" {
			ec.Add(fieldName, "is required", 1001)
			return errors.New(fieldName + " is required")
		}
		return nil
	}
}

// RequiredWithMsg 验证字段不能为空，支持自定义错误消息
func RequiredWithMsg(customMsg string) Rule {
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		if value == nil || value == "" {
			if customMsg != "" {
				ec.Add(fieldName, customMsg, 1002)
				return errors.New(customMsg)
			}
			ec.Add(fieldName, "is required", 1001)
			return errors.New(fieldName + " is required")
		}
		return nil
	}
}

// MinLength 验证字符串的最小长度
func MinLength(min int) Rule {
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		str, ok := value.(string)
		if !ok {
			ec.Add(fieldName, "must be a string", 1003)
			return errors.New(fieldName + " must be a string")
		}
		if len(str) < min {
			ec.Add(fieldName, fmt.Sprintf("must be at least %d characters long", min), 1004)
			return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
		}
		return nil
	}
}

// MinValue 验证整数的最小值
func MinValue(min int) Rule {
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		v, ok := value.(int)
		if !ok {
			ec.Add(fieldName, "must be an integer", 1005)
			return errors.New(fieldName + " must be an integer")
		}
		if v < min {
			ec.Add(fieldName, fmt.Sprintf("must be at least %d", min), 1006)
			return fmt.Errorf("%s must be at least %d", fieldName, min)
		}
		return nil
	}
}

// MatchRegex 验证字符串是否符合正则表达式
func MatchRegex(pattern string) Rule {
	reg := regexp.MustCompile(pattern)
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		str, ok := value.(string)
		if !ok {
			ec.Add(fieldName, "must be a string", 1007)
			return errors.New(fieldName + " must be a string")
		}
		if !reg.MatchString(str) {
			ec.Add(fieldName, "format is invalid", 1008)
			return errors.New(fieldName + " format is invalid")
		}
		return nil
	}
}

func Enums(pattern ...string) Rule {
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		str, ok := value.(string)
		if !ok {
			ec.Add(fieldName, "must be a string", 1007)
			return errors.New(fieldName + " must be a string")
		}
		for _, validValue := range pattern {
			if str == validValue {
				return nil
			}
		}
		ec.Add(fieldName, "value is not in enums", 1009)
		return errors.New(fieldName + "value is not in enums")
	}
}

// Conditional 验证规则在满足特定条件时才生效
func Conditional(condition func() bool, rule Rule) Rule {
	return func(fieldName string, value any, ec *ErrorsCollector) error {
		if condition() {
			return rule(fieldName, value, ec)
		}
		return nil
	}
}
