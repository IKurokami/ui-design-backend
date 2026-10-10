package domain

import (
	"context"
	"time"
)

type Story struct {
	ID              string           `json:"id"`
	UserID          string           `json:"userId"`
	AuthorName      string           `json:"authorName,omitempty"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	CoverURL        string           `json:"coverUrl"`
	Category        string           `json:"category"`
	Tags            string           `json:"tags"`
	Status          string           `json:"status"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	Chapters        []Chapter        `json:"chapters,omitempty"`
	ChapterCount    int              `json:"chapterCount"`
	PreviewImages   []string         `json:"previewImages"`
	ReadingProgress *ReadingProgress `json:"readingProgress,omitempty"`
}

type Chapter struct {
	ID            string    `json:"id"`
	StoryID       string    `json:"storyId"`
	ChapterNumber int       `json:"chapterNumber"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Illustrations string    `json:"illustrations"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateStoryDTO struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	CoverURL       string            `json:"coverUrl"`
	Category       string            `json:"category"`
	Tags           string            `json:"tags"`
	Status         string            `json:"status"`
	InitialChapter *CreateChapterDTO `json:"initialChapter,omitempty"`
}

type UpdateStoryDTO struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	CoverURL    *string `json:"coverUrl,omitempty"`
	Category    *string `json:"category,omitempty"`
	Tags        *string `json:"tags,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type CreateChapterDTO struct {
	ChapterNumber int    `json:"chapterNumber"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Illustrations string `json:"illustrations"`
}

type UpdateChapterDTO struct {
	Title         *string `json:"title,omitempty"`
	Content       *string `json:"content,omitempty"`
	Illustrations *string `json:"illustrations,omitempty"`
}

type StoryRepository interface {
	SocialRepository
	GetReadingProgress(ctx context.Context, userID, storyID string) (*ReadingProgress, error)
	SaveReadingProgress(ctx context.Context, userID, storyID, chapterID string) error
	CreateStory(ctx context.Context, story *Story) error
	GetStory(ctx context.Context, id string) (*Story, error)
	ListStories(ctx context.Context, userID string) ([]*Story, error)
	ListAll(ctx context.Context) ([]*Story, error)
	ListAllPage(ctx context.Context, page, pageSize int, query string) ([]*Story, int, error)
	ListReadingStories(ctx context.Context, userID string) ([]*Story, error)
	ListPage(ctx context.Context, userID, readerID string, page, pageSize int) ([]*Story, error)
	SearchPage(ctx context.Context, query string, page, pageSize int) ([]*Story, error)
	ListFavorites(ctx context.Context, userID string, page, pageSize int) ([]*Story, error)
	FavoriteIDs(ctx context.Context, userID string) ([]string, error)
	SaveFavorite(ctx context.Context, userID, storyID string) (bool, error)
	RemoveFavorite(ctx context.Context, userID, storyID string) error
	UpdateStory(ctx context.Context, story *Story) error
	DeleteStory(ctx context.Context, id, userID string) error

	CreateChapter(ctx context.Context, chapter *Chapter) error
	GetChapter(ctx context.Context, id string) (*Chapter, error)
	ListChaptersByStory(ctx context.Context, storyID string) ([]*Chapter, error)
	UpdateChapter(ctx context.Context, chapter *Chapter) error
	DeleteChapter(ctx context.Context, id string) error
}
