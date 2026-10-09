package domain

import (
	"context"
	"time"
)

// A regular reader/author has no management role.
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
)

func (r Role) IsStaff() bool { return r == RoleAdmin || r == RoleModerator }
func (r Role) IsValid() bool { return r == "" || r.IsStaff() }

type User struct {
	ID           string    `json:"id"`
	DisplayName  string    `json:"displayName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserDTO struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Role        Role   `json:"role,omitempty"`
}

type AuthSessionResponse struct {
	User UserDTO `json:"user"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	ListPage(ctx context.Context, page, pageSize int, query, role string) ([]*User, int, error)
	UpdateRole(ctx context.Context, userID string, role Role) error
	Delete(ctx context.Context, userID string) error
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
	UpdateDisplayName(ctx context.Context, userID, displayName string) error
}
