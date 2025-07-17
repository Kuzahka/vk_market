package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"vk/internal/usecase"
)

// AdHandler обрабатывает HTTP-запросы, связанные с объявлениями.
type AdHandler struct {
	adUseCase *usecase.AdUseCase
}

func NewAdHandler(adUseCase *usecase.AdUseCase) *AdHandler {
	return &AdHandler{adUseCase: adUseCase}
}

type CreateAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

type AdResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	IsOwner     bool      `json:"is_owner,omitempty"` // Дополнительное поле для авторизованных пользователей
}

// CreateAd обрабатывает запрос на создание нового объявления.
func (h *AdHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Message: "Не авторизован: ID пользователя не найден в контексте"})
		return
	}

	var req CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Неверная полезная нагрузка запроса", Details: err.Error()})
		return
	}

	ad, err := h.adUseCase.CreateAd(userID, req.Title, req.Description, req.ImageURL, req.Price)
	if err != nil {
		var validationErr *usecase.ValidationErr
		if errors.As(err, &validationErr) {
			writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Ошибка валидации", Details: err.Error()})
			return
		}
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: "Не удалось создать объявление", Details: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusCreated, AdResponse{
		ID:          ad.ID,
		UserID:      ad.UserID,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price,
		CreatedAt:   ad.CreatedAt,
		IsOwner:     true, // Создатель всегда является владельцем
	})
}

type ListAdsResponse struct {
	Ads        []AdResponse `json:"ads"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
}

// GetAdsFeed обрабатывает запрос на получение ленты объявлений.
func (h *AdHandler) GetAdsFeed(w http.ResponseWriter, r *http.Request) {
	params := usecase.ListAdsParameters{}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	var err error
	params.Page, err = strconv.Atoi(pageStr)
	if err != nil || params.Page < 1 {
		params.Page = 1
	}

	params.Limit, err = strconv.Atoi(limitStr)
	if err != nil || params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	params.SortBy = sortBy
	params.SortOrder = sortOrder

	params.MinPrice, err = strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		params.MinPrice = 0
	}

	params.MaxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		params.MaxPrice = 0
	}

	ads, totalCount, err := h.adUseCase.ListAds(params)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Message: "Не удалось получить объявления", Details: err.Error()})
		return
	}

	// Получаем ID текущего пользователя из контекста (если авторизован)
	currentUserID, _ := r.Context().Value(ContextKeyUserID).(string)

	var adResponses []AdResponse
	for _, ad := range ads {
		isOwner := ad.UserID == currentUserID
		adResponses = append(adResponses, AdResponse{
			ID:          ad.ID,
			UserID:      ad.UserID,
			Title:       ad.Title,
			Description: ad.Description,
			ImageURL:    ad.ImageURL,
			Price:       ad.Price,
			CreatedAt:   ad.CreatedAt,
			IsOwner:     isOwner,
		})
	}

	writeJSONResponse(w, http.StatusOK, ListAdsResponse{
		Ads:        adResponses,
		TotalCount: totalCount,
		Page:       params.Page,
		Limit:      params.Limit,
	})
}
