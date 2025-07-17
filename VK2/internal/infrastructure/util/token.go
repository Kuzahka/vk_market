package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("неверный или истекший токен")
)

// GenerateToken генерирует простой токен, состоящий из userID и временной метки,
// подписанный с использованием HMAC-SHA256.
// Формат токена: Base64(userID.timestamp.HMAC_Signature)
func GenerateToken(userID string, secretKey string, expiration time.Duration) (string, error) {
	// Создаем временную метку истечения срока действия
	expiresAt := time.Now().Add(expiration).Unix()

	// Составляем данные для подписи: userID.expiresAt
	dataToSign := fmt.Sprintf("%s.%d", userID, expiresAt)

	// Создаем HMAC-подпись
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(dataToSign))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	token := fmt.Sprintf("%s.%d.%s", userID, expiresAt, signature)

	encodedToken := base64.URLEncoding.EncodeToString([]byte(token))

	return encodedToken, nil
}

// ParseToken парсит и валидирует кастомный токен, возвращая ID пользователя.
func ParseToken(tokenString string, secretKey string) (string, error) {
	// Декодируем токен из Base64
	decodedBytes, err := base64.URLEncoding.DecodeString(tokenString)
	if err != nil {
		return "", ErrInvalidToken
	}
	decodedToken := string(decodedBytes)

	parts := strings.Split(decodedToken, ".")
	if len(parts) != 3 {
		return "", ErrInvalidToken
	}

	userID := parts[0]
	expiresAtStr := parts[1]
	receivedSignature := parts[2]

	// Проверяем срок действия
	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return "", ErrInvalidToken
	}
	if time.Now().Unix() > expiresAt {
		return "", ErrInvalidToken // Токен истек
	}

	// Повторно генерируем подпись для проверки
	dataToSign := fmt.Sprintf("%s.%d", userID, expiresAt)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(dataToSign))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Сравниваем подписи
	if !hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		return "", ErrInvalidToken
	}

	return userID, nil
}
