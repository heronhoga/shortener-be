package repository

import (
	"context"
	"errors"

	"github.com/heronhoga/shortener-be/model"
	"github.com/jackc/pgx/v5"
)

type LinkRepository interface {
	CheckExistingLink(ctx context.Context, name string) (bool, error)
	CreateNewLink(ctx context.Context, newLink *model.Link) error
}

type linkRepository struct {
	db *pgx.Conn
}

func NewLinkRepository(db *pgx.Conn) LinkRepository {
	return &linkRepository{db: db}
}

func (r *linkRepository) CheckExistingLink(ctx context.Context, name string) (bool, error) {
    query := `SELECT id FROM links WHERE name = $1 LIMIT 1`

    row := r.db.QueryRow(ctx, query, name)

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

func (r *linkRepository) CreateNewLink(ctx context.Context, newLink *model.Link) error {
	query := `INSERT INTO links (id, user_id, name, url, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, newLink.ID, newLink.UserID, newLink.Name, newLink.Url, newLink.CreatedAt)
	return err
}