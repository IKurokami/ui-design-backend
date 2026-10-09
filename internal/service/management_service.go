package service

import (
	"context"
	"errors"

	"ui-design-backend/internal/domain"
)

var (
	ErrInvalidRole = errors.New("invalid management role")
	ErrSelfDelete  = errors.New("cannot delete current admin")
)

type ManagementService struct {
	users   domain.UserRepository
	stories *StoryService
}

func NewManagementService(users domain.UserRepository, stories *StoryService) *ManagementService {
	return &ManagementService{users: users, stories: stories}
}

func (s *ManagementService) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.users.List(ctx)
}
func (s *ManagementService) ListUsersPage(ctx context.Context, page, pageSize int, query, role string) ([]*domain.User, int, error) {
	return s.users.ListPage(ctx, page, pageSize, query, role)
}
func (s *ManagementService) ListStories(ctx context.Context) ([]*domain.Story, error) {
	return s.stories.ListAll(ctx)
}
func (s *ManagementService) ListStoriesPage(ctx context.Context, page, pageSize int, query string) ([]*domain.Story, int, error) {
	return s.stories.ListAllPage(ctx, page, pageSize, query)
}
func (s *ManagementService) SetRole(ctx context.Context, id string, role domain.Role) error {
	if !role.IsStaff() {
		return ErrInvalidRole
	}
	return s.users.UpdateRole(ctx, id, role)
}
func (s *ManagementService) DeleteUser(ctx context.Context, actorID, id string) error {
	if actorID == id {
		return ErrSelfDelete
	}
	return s.users.Delete(ctx, id)
}
