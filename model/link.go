package model

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name string `json:"string"`
	Url string `json:"url"`
	Label string `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateLink struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Label string `json:"label"`
}

type EditLink struct {
	ID	string `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
	Label string `json:"label"`
}

type GetLink struct {
	LinkID string `json:"link_id"`
	Page int `json:"page"`
}

type DeleteLink struct {
	LinkID string `json:"link_id"`
}