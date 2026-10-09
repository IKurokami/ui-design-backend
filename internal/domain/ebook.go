package domain

import (
	"context"
	"encoding/json"
	"time"
)

// Ebook designs are independent of stories and chapters. Design is a versioned
// document, preserving canvas positions, styles and pages without altering prose.
type Ebook struct {
	ID        string          `json:"id"`
	UserID    string          `json:"userId"`
	Title     string          `json:"title"`
	CoverURL  string          `json:"coverUrl"`
	PageCount int             `json:"pageCount"`
	Revision  int             `json:"revision"`
	Design    json.RawMessage `json:"design,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type SaveEbookDTO struct {
	Title            string          `json:"title"`
	CoverURL         string          `json:"coverUrl"`
	Design           json.RawMessage `json:"design"`
	ExpectedRevision int             `json:"expectedRevision"`
}

type EbookRepository interface {
	Create(ctx context.Context, ebook *Ebook) error
	Get(ctx context.Context, userID, id string) (*Ebook, error)
	List(ctx context.Context, userID string) ([]*Ebook, error)
	Update(ctx context.Context, ebook *Ebook, expectedRevision int) error
	Delete(ctx context.Context, userID, id string) error
}
