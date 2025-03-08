package utils

// AuthError 表示认证相关的错误
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// NewAuthError 创建一个新的认证错误
func NewAuthError(msg string) error {
	return &AuthError{Message: msg}
}

// ValidationError 表示数据验证相关的错误
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return "验证错误: " + e.Field + " - " + e.Message
	}
	return "验证错误: " + e.Message
}

// NewValidationError 创建一个新的验证错误
func NewValidationError(msg string, field string) error {
	return &ValidationError{Message: msg, Field: field}
}

// NotFoundError 表示资源未找到的错误
type NotFoundError struct {
	Message  string
	Resource string
}

func (e *NotFoundError) Error() string {
	if e.Resource != "" {
		return "未找到资源: " + e.Resource + " - " + e.Message
	}
	return "未找到资源: " + e.Message
}

// NewNotFoundError 创建一个新的资源未找到错误
func NewNotFoundError(msg string, resource string) error {
	return &NotFoundError{Message: msg, Resource: resource}
}

// DatabaseError 表示数据库操作相关的错误
type DatabaseError struct {
	Message string
	Op      string
}

func (e *DatabaseError) Error() string {
	if e.Op != "" {
		return "数据库错误: " + e.Op + " - " + e.Message
	}
	return "数据库错误: " + e.Message
}

// NewDatabaseError 创建一个新的数据库错误
func NewDatabaseError(msg string, op string) error {
	return &DatabaseError{Message: msg, Op: op}
}
