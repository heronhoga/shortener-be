package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/heronhoga/shortener-be/model"
	"github.com/jackc/pgx/v5"
)

type LinkRepository interface {
	CheckExistingLink(ctx context.Context, name string) (bool, error)
	CreateNewLink(ctx context.Context, newLink *model.Link) error
    GetSpecificLinkById(ctx context.Context, uuid string) (*model.Link, error)
    UpdateSpecificLink(ctx context.Context, existingLink *model.Link) error
	GetShortLinks(ctx context.Context, userId string, limit int, offset int) ([]*model.Link, error)
	DeleteLink(ctx context.Context, linkID string, userID string) error 
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

func (r *linkRepository) GetSpecificLinkById(ctx context.Context, id string) (*model.Link, error) {
	query := `SELECT id, user_id, name, url, created_at, updated_at FROM links WHERE id = $1 LIMIT 1`

	var link model.Link

	err := r.db.QueryRow(ctx, query, id).Scan(
		&link.ID,
		&link.UserID,
		&link.Name,
		&link.Url,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
		return nil, err
	}

	return &link, nil
}

func (r *linkRepository) GetShortLinks(ctx context.Context, userID string, limit int, offset int) ([]*model.Link, error) {
	query := `
		SELECT *
		FROM links
		WHERE user_id = $1
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*model.Link

	for rows.Next() {
		var link model.Link
		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Name,
			&link.Url,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}


func (r *linkRepository) UpdateSpecificLink(ctx context.Context, existingLink *model.Link) error {
    query := `UPDATE links SET name = $1, url = $2, updated_at = $3 WHERE id = $4`
    _, err := r.db.Exec(ctx, query, existingLink.Name, existingLink.Url, existingLink.UpdatedAt, existingLink.ID)

    if err != nil {
        return err
    }

    return nil
}

func (r *linkRepository) DeleteLink(ctx context.Context, linkID string, userID string) error {
	query := `DELETE FROM links WHERE id = $1 AND user_id = $2`
	cmd, err := r.db.Exec(ctx, query, linkID, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

