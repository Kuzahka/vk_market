package repository

import "vk/internal/domain"

// AdRepository определяет интерфейс для взаимодействия с хранилищем объявлений.
type AdRepository interface {
	// CreateAd сохраняет новое объявление в хранилище.
	CreateAd(ad *domain.Ad) error
	// GetAdByID находит объявление по ID.
	GetAdByID(id string) (*domain.Ad, error)
	// ListAds возвращает список объявлений с учетом пагинации, сортировки и фильтрации.
	ListAds(offset, limit int, sortBy, sortOrder string, minPrice, maxPrice float64) ([]domain.Ad, error)
	// CountAds возвращает общее количество объявлений с учетом фильтрации.
	CountAds(minPrice, maxPrice float64) (int, error)
}
