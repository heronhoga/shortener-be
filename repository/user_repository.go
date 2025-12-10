package repository

import (
	"context"

	"github.com/heronhoga/shortener-be/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.RegisterUser) error
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.RegisterUser) error {
	query := `INSERT INTO users (id, email, username, password, phone, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Username, user.Password, user.Phone, user.CreatedAt)
	return err
}