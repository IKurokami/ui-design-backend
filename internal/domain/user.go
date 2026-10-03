package domain

import (
	"context"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	DisplayName  string    `json:"displayName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserDTO struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type AuthSessionResponse struct {
	User UserDTO `json:"user"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
}
