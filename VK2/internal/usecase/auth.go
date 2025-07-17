package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vk/internal/adapter/repository"
	"vk/internal/domain"
	"vk/internal/infrastructure/util"
)

var (
	ErrUserAlreadyExists  = errors.New("пользователь с таким логином уже существует")
	ErrInvalidCredentials = errors.New("неверные учетные данные")
)

type ValidationErr struct {
	Message string
}

func (e *ValidationErr) Error() string {
	return e.Message
}

type AuthUseCase struct {
	userRepo        repository.UserRepository
	tokenSecretKey  string
	tokenExpiration time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, tokenSecretKey string, tokenExpiration time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		tokenSecretKey:  tokenSecretKey,
		tokenExpiration: tokenExpiration,
	}
}

// RegisterUser регистрирует нового пользователя.
func (uc *AuthUseCase) RegisterUser(login, password string) (*domain.User, error) {

	if len(login) < 3 || len(login) > 50 {
		return nil, &ValidationErr{Message: "логин должен быть от 3 до 50 символов"}
	}
	if !isValidLogin(login) {
		return nil, &ValidationErr{Message: "логин может содержать только буквы, цифры, подчеркивания и дефисы"}
	}
	if len(password) < 8 || len(password) > 100 {
		return nil, &ValidationErr{Message: "пароль должен быть от 8 до 100 символов"}
	}
	if !isValidPassword(password) {
		return nil, &ValidationErr{Message: "пароль должен содержать хотя бы одну заглавную букву, одну строчную букву, одну цифру и один специальный символ"}
	}

	// Проверка на существование пользователя
	existingUser, err := uc.userRepo.GetUserByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("не удалось проверить существующего пользователя: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Хеширование пароля
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("не удалось хешировать пароль: %w", err)
	}

	// Создание нового пользователя
	newUser := &domain.User{
		ID:           uuid.New().String(),
		Login:        login,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now().UTC(),
	}

	// Сохранение пользователя в репозитории
	if err := uc.userRepo.CreateUser(newUser); err != nil {
		return nil, fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	return newUser, nil
}

// AuthenticateUser аутентифицирует пользователя и возвращает токен.
func (uc *AuthUseCase) AuthenticateUser(login, password string) (string, error) {
	user, err := uc.userRepo.GetUserByLogin(login)
	if err != nil {
		return "", fmt.Errorf("не удалось получить пользователя по логину: %w", err)
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	// Проверка пароля
	if !util.CheckPasswordHash(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	// Генерация кастомного токена
	token, err := util.GenerateToken(user.ID, uc.tokenSecretKey, uc.tokenExpiration)
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать токен: %w", err)
	}

	return token, nil
}

// isValidLogin проверяет соответствие логина требованиям (буквы, цифры, _, -)
func isValidLogin(login string) bool {

	for _, r := range login {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return false
		}
	}
	return true
}

// isValidPassword проверяет соответствие пароля требованиям (минимум 1 заглавная, 1 строчная, 1 цифра, 1 спецсимвол)
func isValidPassword(password string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;':\",.<>/?`~", r):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
