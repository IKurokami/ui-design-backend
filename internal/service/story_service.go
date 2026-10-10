package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"ui-design-backend/internal/domain"
)

var (
	ErrStoryTitleRequired   = errors.New("story title is required")
	ErrChapterTitleRequired = errors.New("chapter title is required")
	ErrForbidden            = errors.New("forbidden")
)

type StoryService struct {
	repo  domain.StoryRepository
	users domain.UserRepository
}

func NewStoryService(repo domain.StoryRepository, users ...domain.UserRepository) *StoryService {
	s := &StoryService{repo: repo}
	if len(users) > 0 {
		s.users = users[0]
	}
	return s
}

func (s *StoryService) canManage(ctx context.Context, userID, ownerID string) bool {
	if userID == ownerID {
		return true
	}
	if s.users == nil {
		return false
	}
	user, err := s.users.FindByID(ctx, userID)
	return err == nil && user.Role.IsStaff()
}

func (s *StoryService) CreateStory(ctx context.Context, userID string, req *domain.CreateStoryDTO) (*domain.Story, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrStoryTitleRequired
	}

	coverURL := strings.TrimSpace(req.CoverURL)
	// If cover is empty, but initial chapter has illustrations, automatically use the first illustration!
	if coverURL == "" && req.InitialChapter != nil && req.InitialChapter.Illustrations != "" {
		var imgList []string
		if err := json.Unmarshal([]byte(req.InitialChapter.Illustrations), &imgList); err == nil && len(imgList) > 0 {
			coverURL = imgList[0]
		}
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "published"
	}

	now := time.Now().UTC()
	story := &domain.Story{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: req.Description,
		CoverURL:    coverURL,
		Category:    req.Category,
		Tags:        req.Tags,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateStory(ctx, story); err != nil {
		return nil, err
	}

	if req.InitialChapter != nil {
		chTitle := strings.TrimSpace(req.InitialChapter.Title)
		if chTitle == "" {
			chTitle = "Chương 1"
		}
		chNum := req.InitialChapter.ChapterNumber
		if chNum <= 0 {
			chNum = 1
		}

		chapter := &domain.Chapter{
			ID:            uuid.New().String(),
			StoryID:       story.ID,
			ChapterNumber: chNum,
			Title:         chTitle,
			Content:       req.InitialChapter.Content,
			Illustrations: req.InitialChapter.Illustrations,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := s.repo.CreateChapter(ctx, chapter); err != nil {
			return nil, err
		}
		story.Chapters = append(story.Chapters, *chapter)
	}

	story.RefreshPreviewImages()
	return story, nil
}

func (s *StoryService) GetStory(ctx context.Context, id string) (*domain.Story, error) {
	return s.repo.GetStory(ctx, id)
}

func (s *StoryService) ListStories(ctx context.Context, userID string) ([]*domain.Story, error) {
	return s.repo.ListStories(ctx, userID)
}

func (s *StoryService) ListAll(ctx context.Context) ([]*domain.Story, error) {
	return s.repo.ListAll(ctx)
}
func (s *StoryService) ListAllPage(ctx context.Context, page, pageSize int, query string) ([]*domain.Story, int, error) {
	return s.repo.ListAllPage(ctx, page, pageSize, query)
}

func (s *StoryService) UpdateStory(ctx context.Context, userID, id string, req *domain.UpdateStoryDTO) (*domain.Story, error) {
	story, err := s.repo.GetStory(ctx, id)
	if err != nil {
		return nil, err
	}

	if !s.canManage(ctx, userID, story.UserID) {
		return nil, ErrForbidden
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed != "" {
			story.Title = trimmed
		}
	}
	if req.Description != nil {
		story.Description = *req.Description
	}
	if req.CoverURL != nil {
		story.CoverURL = *req.CoverURL
	}
	if req.Category != nil {
		story.Category = *req.Category
	}
	if req.Tags != nil {
		story.Tags = *req.Tags
	}
	if req.Status != nil {
		story.Status = *req.Status
	}
	story.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateStory(ctx, story); err != nil {
		return nil, err
	}
	story.RefreshPreviewImages()
	return story, nil
}

func (s *StoryService) DeleteStory(ctx context.Context, userID, id string) error {
	story, err := s.repo.GetStory(ctx, id)
	if err != nil {
		return err
	}
	if !s.canManage(ctx, userID, story.UserID) {
		return ErrForbidden
	}
	return s.repo.DeleteStory(ctx, id, story.UserID)
}

func (s *StoryService) CreateChapter(ctx context.Context, userID, storyID string, req *domain.CreateChapterDTO) (*domain.Chapter, error) {
	story, err := s.repo.GetStory(ctx, storyID)
	if err != nil {
		return nil, err
	}

	if !s.canManage(ctx, userID, story.UserID) {
		return nil, ErrForbidden
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrChapterTitleRequired
	}

	now := time.Now().UTC()
	chNum := req.ChapterNumber
	if chNum <= 0 {
		chNum = len(story.Chapters) + 1
	}

	chapter := &domain.Chapter{
		ID:            uuid.New().String(),
		StoryID:       storyID,
		ChapterNumber: chNum,
		Title:         title,
		Content:       req.Content,
		Illustrations: req.Illustrations,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateChapter(ctx, chapter); err != nil {
		return nil, err
	}

	// Auto-cover: if story cover is empty and chapter has an illustration, set it as cover!
	if story.CoverURL == "" && req.Illustrations != "" {
		var imgList []string
		if err := json.Unmarshal([]byte(req.Illustrations), &imgList); err == nil && len(imgList) > 0 {
			story.CoverURL = imgList[0]
			story.UpdatedAt = now
			_ = s.repo.UpdateStory(ctx, story)
		}
	}

	return chapter, nil
}

func (s *StoryService) GetChapter(ctx context.Context, id string) (*domain.Chapter, error) {
	return s.repo.GetChapter(ctx, id)
}

func (s *StoryService) UpdateChapter(ctx context.Context, userID, id string, req *domain.UpdateChapterDTO) (*domain.Chapter, error) {
	chapter, err := s.repo.GetChapter(ctx, id)
	if err != nil {
		return nil, err
	}

	story, err := s.repo.GetStory(ctx, chapter.StoryID)
	if err != nil {
		return nil, err
	}

	if !s.canManage(ctx, userID, story.UserID) {
		return nil, ErrForbidden
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed != "" {
			chapter.Title = trimmed
		}
	}
	if req.Content != nil {
		chapter.Content = *req.Content
	}
	if req.Illustrations != nil {
		chapter.Illustrations = *req.Illustrations
	}
	chapter.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateChapter(ctx, chapter); err != nil {
		return nil, err
	}

	// Auto-cover: if story cover is empty and illustrations added, update cover
	if story.CoverURL == "" && chapter.Illustrations != "" {
		var imgList []string
		if err := json.Unmarshal([]byte(chapter.Illustrations), &imgList); err == nil && len(imgList) > 0 {
			story.CoverURL = imgList[0]
			story.UpdatedAt = time.Now().UTC()
			_ = s.repo.UpdateStory(ctx, story)
		}
	}

	return chapter, nil
}

func (s *StoryService) DeleteChapter(ctx context.Context, userID, id string) error {
	chapter, err := s.repo.GetChapter(ctx, id)
	if err != nil {
		return err
	}

	story, err := s.repo.GetStory(ctx, chapter.StoryID)
	if err != nil {
		return err
	}

	if !s.canManage(ctx, userID, story.UserID) {
		return ErrForbidden
	}

	return s.repo.DeleteChapter(ctx, id)
}

func (s *StoryService) GetReadingProgress(ctx context.Context, userID, storyID string) (*domain.ReadingProgress, error) {
	return s.repo.GetReadingProgress(ctx, userID, storyID)
}
func (s *StoryService) SaveReadingProgress(ctx context.Context, userID, storyID, chapterID string) error {
	return s.repo.SaveReadingProgress(ctx, userID, storyID, chapterID)
}

func (s *StoryService) ListReadingStories(ctx context.Context, userID string) ([]*domain.Story, error) {
	return s.repo.ListReadingStories(ctx, userID)
}
func (s *StoryService) SearchPage(ctx context.Context, query string, page, pageSize int) ([]*domain.Story, error) {
	return s.repo.SearchPage(ctx, query, page, pageSize)
}
func (s *StoryService) ListFavorites(ctx context.Context, userID string, page, pageSize int) ([]*domain.Story, error) {
	return s.repo.ListFavorites(ctx, userID, page, pageSize)
}
func (s *StoryService) FavoriteIDs(ctx context.Context, userID string) ([]string, error) {
	return s.repo.FavoriteIDs(ctx, userID)
}
func (s *StoryService) SaveFavorite(ctx context.Context, userID, storyID string) (bool, error) {
	return s.repo.SaveFavorite(ctx, userID, storyID)
}
func (s *StoryService) RemoveFavorite(ctx context.Context, userID, storyID string) error {
	return s.repo.RemoveFavorite(ctx, userID, storyID)
}
func (s *StoryService) ListPage(ctx context.Context, userID, readerID string, page, pageSize int) ([]*domain.Story, error) {
	return s.repo.ListPage(ctx, userID, readerID, page, pageSize)
}

func (s *StoryService) FollowIDs(ctx context.Context, userID string) ([]string, error) {
	return s.repo.FollowIDs(ctx, userID)
}
func (s *StoryService) Follow(ctx context.Context, userID, storyID string) (bool, error) {
	return s.repo.Follow(ctx, userID, storyID)
}
func (s *StoryService) Unfollow(ctx context.Context, userID, storyID string) error {
	return s.repo.Unfollow(ctx, userID, storyID)
}
func (s *StoryService) ListNotifications(ctx context.Context, userID string, limit int) ([]domain.StoryNotification, error) {
	return s.repo.ListNotifications(ctx, userID, limit)
}
func (s *StoryService) MarkNotificationRead(ctx context.Context, userID, chapterID string) error {
	return s.repo.MarkNotificationRead(ctx, userID, chapterID)
}
func (s *StoryService) ListComments(ctx context.Context, storyID, userID string) ([]domain.StoryComment, error) {
	return s.repo.ListComments(ctx, storyID, userID)
}
func (s *StoryService) CreateComment(ctx context.Context, storyID, userID, body string) (*domain.StoryComment, error) {
	return s.repo.CreateComment(ctx, storyID, userID, body)
}
func (s *StoryService) DeleteComment(ctx context.Context, userID, commentID string) error {
	return s.repo.DeleteComment(ctx, userID, commentID)
}
func (s *StoryService) ListPendingComments(ctx context.Context, limit int) ([]domain.StoryComment, error) {
	return s.repo.ListPendingComments(ctx, limit)
}
func (s *StoryService) ModerateComment(ctx context.Context, id, status string) error {
	return s.repo.ModerateComment(ctx, id, status)
}

func (s *StoryService) MutedNotificationStoryIDs(ctx context.Context, userID string) ([]string, error) {
	return s.repo.MutedNotificationStoryIDs(ctx, userID)
}
func (s *StoryService) MuteStoryNotifications(ctx context.Context, userID, storyID string) error {
	return s.repo.MuteStoryNotifications(ctx, userID, storyID)
}
func (s *StoryService) UnmuteStoryNotifications(ctx context.Context, userID, storyID string) error {
	return s.repo.UnmuteStoryNotifications(ctx, userID, storyID)
}

func (s *StoryService) ListMutedStories(ctx context.Context, userID string) ([]domain.MutedStory, error) {
	return s.repo.ListMutedStories(ctx, userID)
}
func (s *StoryService) DismissNotification(ctx context.Context, userID, chapterID string) error {
	return s.repo.DismissNotification(ctx, userID, chapterID)
}
func (s *StoryService) ClearNotifications(ctx context.Context, userID string) error {
	return s.repo.ClearNotifications(ctx, userID)
}
