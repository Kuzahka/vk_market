package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"vk/internal/adapter/repository"
	"vk/internal/domain"
)

type PGAdRepository struct {
	db *sql.DB
}

func NewPGAdRepository(db *sql.DB) repository.AdRepository {
	return &PGAdRepository{db: db}
}

// CreateAd реализует метод создания объявления для PostgreSQL.
func (r *PGAdRepository) CreateAd(ad *domain.Ad) error {
	query := `INSERT INTO ads (id, user_id, title, description, image_url, price, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, ad.ID, ad.UserID, ad.Title, ad.Description, ad.ImageURL, ad.Price, ad.CreatedAt)
	if err != nil {
		log.Printf("Error creating ad in postgres: %v", err)
		return fmt.Errorf("failed to create ad in postgres: %w", err)
	}
	return nil
}

// GetAdByID реализует метод получения объявления по ID для PostgreSQL.
func (r *PGAdRepository) GetAdByID(id string) (*domain.Ad, error) {
	ad := &domain.Ad{}
	query := `SELECT id, user_id, title, description, image_url, price, created_at FROM ads WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&ad.ID, &ad.UserID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ad by ID from postgres: %w", err)
	}
	return ad, nil
}

// ListAds реализует метод получения списка объявлений с пагинацией, сортировкой и фильтрацией для PostgreSQL.
func (r *PGAdRepository) ListAds(offset, limit int, sortBy, sortOrder string, minPrice, maxPrice float64) ([]domain.Ad, error) {
	var ads []domain.Ad
	args := []interface{}{}
	whereClauses := []string{}
	argCounter := 1

	if minPrice > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argCounter))
		args = append(args, minPrice)
		argCounter++
	}
	if maxPrice > 0 && maxPrice >= minPrice {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argCounter))
		args = append(args, maxPrice)
		argCounter++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	orderByClause := " ORDER BY created_at DESC"
	if sortBy != "" {
		switch sortBy {
		case "created_at":
			orderByClause = " ORDER BY created_at"
		case "price":
			orderByClause = " ORDER BY price"
		default:

		}
	}

	if sortOrder != "" {
		if strings.ToUpper(sortOrder) == "ASC" {
			orderByClause += " ASC"
		} else {
			orderByClause += " DESC"
		}
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, title, description, image_url, price, created_at
		FROM ads
		%s
		%s
		OFFSET $%d LIMIT $%d`,
		whereClause,
		orderByClause,
		argCounter,
		argCounter+1,
	)
	args = append(args, offset, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list ads from postgres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		ad := domain.Ad{}
		if err := rows.Scan(&ad.ID, &ad.UserID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan ad row: %w", err)
		}
		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return ads, nil
}

// CountAds реализует метод подсчета объявлений с учетом фильтрации для PostgreSQL.
func (r *PGAdRepository) CountAds(minPrice, maxPrice float64) (int, error) {
	args := []interface{}{}
	whereClauses := []string{}
	argCounter := 1

	if minPrice > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argCounter))
		args = append(args, minPrice)
		argCounter++
	}
	if maxPrice > 0 && maxPrice >= minPrice {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argCounter))
		args = append(args, maxPrice)
		argCounter++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM ads %s`, whereClause)

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count ads from postgres: %w", err)
	}
	return count, nil
}
