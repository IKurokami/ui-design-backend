package domain

import (
	"context"
	"time"
)

type PasswordReset struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UsedAt    *time.Time `json:"usedAt,omitempty"`
}

func (pr *PasswordReset) IsValid() bool {
	return pr.UsedAt == nil && time.Now().Before(pr.ExpiresAt)
}

type PasswordResetRepository interface {
	Create(ctx context.Context, pr *PasswordReset) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordReset, error)
	MarkAsUsed(ctx context.Context, id string) error
	InvalidateActiveForUser(ctx context.Context, userID string) error
}
