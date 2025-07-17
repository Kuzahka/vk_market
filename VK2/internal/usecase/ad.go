package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vk/internal/adapter/repository"
	"vk/internal/domain"
)

var (
	ErrAdNotFound = errors.New("ad not found")
)

type AdUseCase struct {
	adRepo repository.AdRepository
}

func NewAdUseCase(adRepo repository.AdRepository) *AdUseCase {
	return &AdUseCase{adRepo: adRepo}
}

// CreateAd создает новое объявление.
func (uc *AdUseCase) CreateAd(userID, title, description, imageURL string, price float64) (*domain.Ad, error) {
	// TODO: Добавить валидацию входных данных (длина, формат, цена > 0)
	if title == "" || price <= 0 {
		return nil, errors.New("title cannot be empty and price must be greater than 0")
	}
	if len(title) > 255 {
		return nil, errors.New("title is too long")
	}
	if len(description) > 1000 {
		return nil, errors.New("description is too long")
	}

	newAd := &domain.Ad{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		CreatedAt:   time.Now().UTC(),
	}

	if err := uc.adRepo.CreateAd(newAd); err != nil {
		return nil, fmt.Errorf("failed to create ad: %w", err)
	}

	return newAd, nil
}

type ListAdsParameters struct {
	Page      int
	Limit     int
	SortBy    string // created_at, price
	SortOrder string // asc, desc
	MinPrice  float64
	MaxPrice  float64
}

// ListAds возвращает список объявлений с учетом пагинации, сортировки и фильтрации.
func (uc *AdUseCase) ListAds(params ListAdsParameters) ([]domain.Ad, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 { // Ограничение на размер страницы
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	ads, err := uc.adRepo.ListAds(offset, params.Limit, params.SortBy, params.SortOrder, params.MinPrice, params.MaxPrice)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list ads: %w", err)
	}

	totalCount, err := uc.adRepo.CountAds(params.MinPrice, params.MaxPrice)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ads: %w", err)
	}

	return ads, totalCount, nil
}
