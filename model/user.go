package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone string `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Session string `json:"session"`
	Otp string `json:"otp"`
}

type RegisterUser struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}