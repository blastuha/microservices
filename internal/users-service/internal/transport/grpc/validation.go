package grpc

import (
	"fmt"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// validateEmail проверяет формат email
func validateEmail(email string) error {
	if email == "" {
		return ValidationError{Field: "email", Message: "email is required"}
	}

	// Простая проверка формата email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "invalid email format"}
	}

	return nil
}

// validatePassword проверяет пароль
func validatePassword(password string) error {
	if password == "" {
		return ValidationError{Field: "password", Message: "password is required"}
	}

	if len(password) < 6 {
		return ValidationError{Field: "password", Message: "password must be at least 6 characters"}
	}

	return nil
}

// validateUserID проверяет ID пользователя
func validateUserID(id uint32) error {
	if id == 0 {
		return ValidationError{Field: "id", Message: "user id is required"}
	}

	return nil
}

// validateCreateUserRequest валидирует запрос создания пользователя
func validateCreateUserRequest(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	return nil
}

// validateUpdateUserRequest валидирует запрос обновления пользователя
func validateUpdateUserRequest(id uint32, email, password *string) error {
	if err := validateUserID(id); err != nil {
		return err
	}

	// 1. В PATCH методе email/pass может быть опциональным, то есть не передан, то есть nil
	// 2. Если передан, то он не может быть пустым
	// 3. Если передан, то он должен быть валидным (проверка формата email)
	// 4. Если передан, то он должен быть не меньше 6 символов (проверка длины пароля)

	// Email и password опциональны при обновлении
	// Если поле передано, оно должно быть не пустым и валидным
	if email != nil {
		if *email == "" {
			return ValidationError{Field: "email", Message: "email cannot be empty if provided"}
		}
		if err := validateEmail(*email); err != nil {
			return err
		}
	}

	if password != nil {
		if *password == "" {
			return ValidationError{Field: "password", Message: "password cannot be empty if provided"}
		}
		if err := validatePassword(*password); err != nil {
			return err
		}
	}

	return nil
}

// handleValidationError конвертирует ошибку валидации в gRPC статус
func handleValidationError(err error) error {
	if validationErr, ok := err.(ValidationError); ok {
		return status.Errorf(codes.InvalidArgument, "%s: %s", validationErr.Field, validationErr.Message)
	}

	// Если это не ValidationError, возвращаем как есть
	return err
}
