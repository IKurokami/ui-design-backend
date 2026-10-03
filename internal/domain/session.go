package domain

import (
	"context"
	"time"
)

type Session struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"userId"`
	TokenHash            string     `json:"-"`
	UserAgent            string     `json:"userAgent"`
	IP                   string     `json:"ip"`
	ExpiresAt            time.Time  `json:"expiresAt"`
	CreatedAt            time.Time  `json:"createdAt"`
	RevokedAt            *time.Time `json:"revokedAt,omitempty"`
	RotatedToID          *string    `json:"rotatedToId,omitempty"`
	GracePeriodExpiresAt *time.Time `json:"gracePeriodExpiresAt,omitempty"`
}

func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) InGracePeriod() bool {
	if s.GracePeriodExpiresAt == nil {
		return false
	}
	return time.Now().Before(*s.GracePeriodExpiresAt)
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*Session, error)
	FindByID(ctx context.Context, id string) (*Session, error)
	Rotate(ctx context.Context, oldSessionID string, newSession *Session, gracePeriod time.Duration) error
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
