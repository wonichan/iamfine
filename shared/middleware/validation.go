package middleware

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator 验证器接口
type Validator interface {
	Validate(value interface{}) error
}

// StringValidator 字符串验证器
type StringValidator struct {
	MinLength int
	MaxLength int
	Pattern   *regexp.Regexp
	Required  bool
}

func (v *StringValidator) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return ValidationError{Message: "value must be a string"}
	}

	// 检查必填
	if v.Required && strings.TrimSpace(str) == "" {
		return ValidationError{Message: "field is required"}
	}

	// 如果不是必填且为空，跳过其他验证
	if !v.Required && strings.TrimSpace(str) == "" {
		return nil
	}

	// 检查长度
	if v.MinLength > 0 && len(str) < v.MinLength {
		return ValidationError{Message: fmt.Sprintf("minimum length is %d", v.MinLength)}
	}
	if v.MaxLength > 0 && len(str) > v.MaxLength {
		return ValidationError{Message: fmt.Sprintf("maximum length is %d", v.MaxLength)}
	}

	// 检查正则表达式
	if v.Pattern != nil && !v.Pattern.MatchString(str) {
		return ValidationError{Message: "format is invalid"}
	}

	return nil
}

// EmailValidator 邮箱验证器
type EmailValidator struct {
	Required bool
}

func (v *EmailValidator) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return ValidationError{Message: "value must be a string"}
	}

	if v.Required && strings.TrimSpace(str) == "" {
		return ValidationError{Message: "email is required"}
	}

	if !v.Required && strings.TrimSpace(str) == "" {
		return nil
	}

	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(str) {
		return ValidationError{Message: "invalid email format"}
	}

	return nil
}

// PhoneValidator 手机号验证器
type PhoneValidator struct {
	Required bool
}

func (v *PhoneValidator) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return ValidationError{Message: "value must be a string"}
	}

	if v.Required && strings.TrimSpace(str) == "" {
		return ValidationError{Message: "phone is required"}
	}

	if !v.Required && strings.TrimSpace(str) == "" {
		return nil
	}

	// 中国手机号验证
	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phonePattern.MatchString(str) {
		return ValidationError{Message: "invalid phone format"}
	}

	return nil
}

// UserIDValidator 用户ID验证器
type UserIDValidator struct {
	Required bool
}

func (v *UserIDValidator) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return ValidationError{Message: "value must be a string"}
	}

	if v.Required && strings.TrimSpace(str) == "" {
		return ValidationError{Message: "user ID is required"}
	}

	if !v.Required && strings.TrimSpace(str) == "" {
		return nil
	}

	// 用户ID格式验证（假设使用xid格式）
	if len(str) != 20 {
		return ValidationError{Message: "invalid user ID format"}
	}

	return nil
}

// ValidateUserRegistration 验证用户注册参数
func ValidateUserRegistration(username, password, email, phone string) []ValidationError {
	var errors []ValidationError

	// 验证用户名
	usernameValidator := &StringValidator{
		MinLength: 3,
		MaxLength: 20,
		Pattern:   regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
		Required:  true,
	}
	if err := usernameValidator.Validate(username); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "username"
			errors = append(errors, ve)
		}
	}

	// 验证密码
	passwordValidator := &StringValidator{
		MinLength: 6,
		MaxLength: 50,
		Required:  true,
	}
	if err := passwordValidator.Validate(password); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "password"
			errors = append(errors, ve)
		}
	}

	// 验证邮箱
	emailValidator := &EmailValidator{Required: false}
	if err := emailValidator.Validate(email); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "email"
			errors = append(errors, ve)
		}
	}

	// 验证手机号
	phoneValidator := &PhoneValidator{Required: true}
	if err := phoneValidator.Validate(phone); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "phone"
			errors = append(errors, ve)
		}
	}

	return errors
}

// ValidatePostCreation 验证帖子创建参数
func ValidatePostCreation(userID, title, content string) []ValidationError {
	var errors []ValidationError

	// 验证用户ID
	userIDValidator := &UserIDValidator{Required: true}
	if err := userIDValidator.Validate(userID); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "user_id"
			errors = append(errors, ve)
		}
	}

	// 验证标题
	titleValidator := &StringValidator{
		MinLength: 1,
		MaxLength: 100,
		Required:  true,
	}
	if err := titleValidator.Validate(title); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "title"
			errors = append(errors, ve)
		}
	}

	// 验证内容
	contentValidator := &StringValidator{
		MinLength: 1,
		MaxLength: 10000,
		Required:  true,
	}
	if err := contentValidator.Validate(content); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "content"
			errors = append(errors, ve)
		}
	}

	return errors
}

// ValidateFollowOperation 验证关注操作参数
func ValidateFollowOperation(userID, targetUserID string) []ValidationError {
	var errors []ValidationError

	// 验证用户ID
	userIDValidator := &UserIDValidator{Required: true}
	if err := userIDValidator.Validate(userID); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "user_id"
			errors = append(errors, ve)
		}
	}

	// 验证目标用户ID
	targetUserIDValidator := &UserIDValidator{Required: true}
	if err := targetUserIDValidator.Validate(targetUserID); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "target_user_id"
			errors = append(errors, ve)
		}
	}

	// 检查是否尝试关注自己
	if userID == targetUserID {
		errors = append(errors, ValidationError{
			Field:   "target_user_id",
			Message: "cannot follow yourself",
		})
	}

	return errors
}

// ValidateCommentCreation 验证评论创建参数
func ValidateCommentCreation(postID, userID, content string) error {
	// 验证帖子ID
	if strings.TrimSpace(postID) == "" {
		return ValidationError{
			Field:   "post_id",
			Message: "post ID is required",
		}
	}

	// 验证用户ID
	userIDValidator := &UserIDValidator{Required: true}
	if err := userIDValidator.Validate(userID); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "user_id"
			return ve
		}
	}

	// 验证评论内容
	contentValidator := &StringValidator{
		MinLength: 1,
		MaxLength: 2000,
		Required:  true,
	}
	if err := contentValidator.Validate(content); err != nil {
		if ve, ok := err.(ValidationError); ok {
			ve.Field = "content"
			return ve
		}
	}

	return nil
}
