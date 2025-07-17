package handler

import (
	"context"
	"net/http"
	"strings"

	"vk/internal/infrastructure/util"
	"vk/internal/usecase"
)

type ContextKey string

const ContextKeyUserID ContextKey = "userID"

// AuthMiddleware проверяет наличие и валидность авторизационного токена.
func AuthMiddleware(tokenSecretKey string, authUseCase *usecase.AuthUseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Message: "Не авторизован: отсутствует заголовок Authorization"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Message: "Не авторизован: неверный формат заголовка Authorization"})
			return
		}

		tokenString := parts[1]

		// Парсим и валидируем кастомный токен
		userID, err := util.ParseToken(tokenString, tokenSecretKey)
		if err != nil {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Message: "Не авторизован: неверный или истекший токен", Details: err.Error()})
			return
		}

		// Добавляем userID в контекст запроса
		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
