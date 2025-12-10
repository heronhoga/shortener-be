package repository

import (
	"context"
	"errors"

	"github.com/heronhoga/shortener-be/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *model.User) error
	CheckExistingEmailUsername(ctx context.Context, email string, username string) (bool, error)
    GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CheckExistingEmailUsername(ctx context.Context, email string, username string) (bool, error) {

    query := `SELECT id FROM users WHERE email = $1 OR username = $2 LIMIT 1`

    row := r.db.QueryRow(ctx, query, email, username)

    var id int64
    err := row.Scan(&id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            // No conflict
            return true, nil
        }
        return false, err
    }

    return false, nil
}

func (r *userRepository) InsertUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, email, username, password, phone, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Username, user.Password, user.Phone, user.CreatedAt)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, password FROM users WHERE email = $1 LIMIT 1`

	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}




