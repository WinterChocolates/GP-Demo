package utils

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

//func NewAuthError(msg string) error {
//	return &AuthError{Message: msg}
//}
