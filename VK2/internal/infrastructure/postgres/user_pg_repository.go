package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"vk/internal/adapter/repository"
	"vk/internal/domain"
)

type PGUserRepository struct {
	db *sql.DB
}

func NewPGUserRepository(db *sql.DB) repository.UserRepository {
	return &PGUserRepository{db: db}
}

// CreateUser реализует метод создания пользователя для PostgreSQL.
func (r *PGUserRepository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (id, login, password_hash, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, user.ID, user.Login, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user in postgres: %w", err)
	}
	return nil
}

// GetUserByLogin реализует метод получения пользователя по логину для PostgreSQL.
func (r *PGUserRepository) GetUserByLogin(login string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	err := r.db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by login from postgres: %w", err)
	}
	return user, nil
}

// GetUserByID реализует метод получения пользователя по ID для PostgreSQL.
func (r *PGUserRepository) GetUserByID(id string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, login, password_hash, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID from postgres: %w", err)
	}
	return user, nil
}

// NewPostgresDB создает и возвращает новое соединение с базой данных PostgreSQL.
func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
