package domain

import (
	"context"
	"time"
)

type StoryComment struct {
	ID         string    `json:"id"`
	StoryID    string    `json:"storyId"`
	UserID     string    `json:"userId"`
	AuthorName string    `json:"authorName"`
	Body       string    `json:"body"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}
type StoryNotification struct {
	StoryID      string    `json:"storyId"`
	ChapterID    string    `json:"chapterId"`
	StoryTitle   string    `json:"storyTitle"`
	ChapterTitle string    `json:"chapterTitle"`
	CreatedAt    time.Time `json:"createdAt"`
	Read         bool      `json:"read"`
}

type MutedStory struct {
	StoryID    string `json:"storyId"`
	StoryTitle string `json:"storyTitle"`
}

type SocialRepository interface {
	ListMutedStories(ctx context.Context, userID string) ([]MutedStory, error)
	DismissNotification(ctx context.Context, userID, chapterID string) error
	ClearNotifications(ctx context.Context, userID string) error
	MutedNotificationStoryIDs(ctx context.Context, userID string) ([]string, error)
	MuteStoryNotifications(ctx context.Context, userID, storyID string) error
	UnmuteStoryNotifications(ctx context.Context, userID, storyID string) error
	FollowIDs(ctx context.Context, userID string) ([]string, error)
	Follow(ctx context.Context, userID, storyID string) (bool, error)
	Unfollow(ctx context.Context, userID, storyID string) error
	ListNotifications(ctx context.Context, userID string, limit int) ([]StoryNotification, error)
	MarkNotificationRead(ctx context.Context, userID, chapterID string) error
	ListComments(ctx context.Context, storyID, userID string) ([]StoryComment, error)
	CreateComment(ctx context.Context, storyID, userID, body string) (*StoryComment, error)
	DeleteComment(ctx context.Context, userID, commentID string) error
	ListPendingComments(ctx context.Context, limit int) ([]StoryComment, error)
	ModerateComment(ctx context.Context, id, status string) error
}
