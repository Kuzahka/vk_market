package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"vk/internal/usecase"
)

// AuthHandler обрабатывает HTTP-запросы, связанные с аутентификацией.
type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// writeJSONResponse - вспомогательная функция для отправки JSON-ответов.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterUser обрабатывает запрос на регистрацию нового пользователя.
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload", Details: err.Error()})
		return
	}

	user, err := h.authUseCase.RegisterUser(req.Login, req.Password)
	if err != nil {
		var validationErr *usecase.ValidationErr
		if errors.As(err, &validationErr) {
			writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Validation error", Details: err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			writeJSONResponse(w, http.StatusConflict, ErrorResponse{Message: err.Error()})
			return
		}
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: "Failed to register user", Details: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, RegisterResponse{
		ID:        user.ID,
		Login:     user.Login,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse представляет структуру ответа на вход.
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginUser обрабатывает запрос на вход пользователя.
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload", Details: err.Error()})
		return
	}

	token, err := h.authUseCase.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
			return
		}
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: "Failed to authenticate user", Details: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, LoginResponse{Token: token})
}
