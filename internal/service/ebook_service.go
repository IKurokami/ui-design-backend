package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"ui-design-backend/internal/domain"
)

var ErrInvalidEbook = errors.New("invalid ebook design")
var ebookColor = regexp.MustCompile(`^#[0-9a-fA-F]{3,8}$`)

type EbookService struct{ repo domain.EbookRepository }

func NewEbookService(repo domain.EbookRepository) *EbookService { return &EbookService{repo: repo} }

func validateEbook(req *domain.SaveEbookDTO) (int, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" || utf8.RuneCountInString(req.Title) > 200 || len(req.Design) > 4*1024*1024 {
		return 0, ErrInvalidEbook
	}
	var document struct {
		Version int `json:"version"`
		Pages   []struct {
			ID         string `json:"id"`
			Kind       string `json:"kind"`
			Background string `json:"background"`
			Accent     string `json:"accent"`
		} `json:"pages"`
		PageBackground string `json:"pageBackground"`
		Accent         string `json:"accent"`
	}
	if json.Unmarshal(req.Design, &document) != nil || (document.Version != 1 && document.Version != 2) || len(document.Pages) == 0 || len(document.Pages) > 500 {
		return 0, ErrInvalidEbook
	}
	ids := make(map[string]bool)
	for _, p := range document.Pages {
		if p.ID == "" || ids[p.ID] || (p.Kind != "cover" && p.Kind != "chapter" && p.Kind != "blank") {
			return 0, ErrInvalidEbook
		}
		if (p.Background != "" && !ebookColor.MatchString(p.Background)) || (p.Accent != "" && !ebookColor.MatchString(p.Accent)) {
			return 0, ErrInvalidEbook
		}
		ids[p.ID] = true
	}
	if (document.PageBackground != "" && !ebookColor.MatchString(document.PageBackground)) || (document.Accent != "" && !ebookColor.MatchString(document.Accent)) {
		return 0, ErrInvalidEbook
	}
	if len(req.CoverURL) > 2048 || (req.CoverURL != "" && !strings.HasPrefix(req.CoverURL, "/") && !strings.HasPrefix(req.CoverURL, "https://") && !strings.HasPrefix(req.CoverURL, "http://")) {
		return 0, ErrInvalidEbook
	}
	return len(document.Pages), nil
}

func (s *EbookService) Create(ctx context.Context, userID string, req *domain.SaveEbookDTO) (*domain.Ebook, error) {
	count, err := validateEbook(req)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	ebook := &domain.Ebook{ID: uuid.New().String(), UserID: userID, Title: req.Title, CoverURL: req.CoverURL, PageCount: count, Revision: 1, Design: req.Design, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.Create(ctx, ebook); err != nil {
		return nil, err
	}
	return ebook, nil
}
func (s *EbookService) Get(ctx context.Context, userID, id string) (*domain.Ebook, error) {
	return s.repo.Get(ctx, userID, id)
}
func (s *EbookService) List(ctx context.Context, userID string) ([]*domain.Ebook, error) {
	return s.repo.List(ctx, userID)
}
func (s *EbookService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}
func (s *EbookService) Update(ctx context.Context, userID, id string, req *domain.SaveEbookDTO) (*domain.Ebook, error) {
	count, err := validateEbook(req)
	if err != nil || req.ExpectedRevision < 1 {
		return nil, ErrInvalidEbook
	}
	ebook, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	ebook.Title, ebook.CoverURL, ebook.PageCount, ebook.Design, ebook.UpdatedAt = req.Title, req.CoverURL, count, req.Design, time.Now().UTC()
	if err := s.repo.Update(ctx, ebook, req.ExpectedRevision); err != nil {
		return nil, err
	}
	return ebook, nil
}
