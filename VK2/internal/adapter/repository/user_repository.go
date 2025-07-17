package repository

import "vk/internal/domain"

// UserRepository определяет интерфейс для взаимодействия с хранилищем пользователей.
type UserRepository interface {
	// CreateUser сохраняет нового пользователя в хранилище.
	CreateUser(user *domain.User) error
	// GetUserByLogin находит пользователя по логину.
	GetUserByLogin(login string) (*domain.User, error)
	// GetUserByID находит пользователя по ID.
	GetUserByID(id string) (*domain.User, error)
}
