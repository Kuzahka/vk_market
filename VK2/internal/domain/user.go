package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(id, login, passwordHash string, createdAt time.Time) *User {
	return &User{
		ID:           id,
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}
}
