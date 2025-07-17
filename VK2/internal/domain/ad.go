package domain

import "time"

type Ad struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewAd(id, userID, title, description, imageURL string, price float64, createdAt time.Time) *Ad {
	return &Ad{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		CreatedAt:   createdAt,
	}
}
