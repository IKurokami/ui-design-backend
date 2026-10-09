package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/repository"
)

type ProfileService struct {
	userRepo    domain.UserRepository
	profileRepo domain.ProfileRepository
}

func NewProfileService(userRepo domain.UserRepository, profileRepo domain.ProfileRepository) *ProfileService {
	return &ProfileService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

func (s *ProfileService) GetOrCreateProfile(ctx context.Context, userID string) (*domain.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrProfileNotFound) {
			now := time.Now().UTC()
			profile = &domain.UserProfile{
				ID:         uuid.New().String(),
				UserID:     userID,
				ThemeColor: "#3B82F6",
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			if err := s.profileRepo.Create(ctx, profile); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &domain.ProfileResponse{
		ID:          profile.ID,
		UserID:      user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   profile.AvatarURL,
		BannerURL:   profile.BannerURL,
		Bio:         profile.Bio,
		Headline:    profile.Headline,
		ThemeColor:  profile.ThemeColor,
		CustomCSS:   profile.CustomCSS,
		WebsiteURL:  profile.WebsiteURL,
		Location:    profile.Location,
		SocialLinks: profile.SocialLinks,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}

func (s *ProfileService) GetPublicProfile(ctx context.Context, userID string) (*domain.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrProfileNotFound) {
			now := time.Now().UTC()
			profile = &domain.UserProfile{
				ID:         uuid.New().String(),
				UserID:     userID,
				ThemeColor: "#3B82F6",
				CreatedAt:  now,
				UpdatedAt:  now,
			}
		} else {
			return nil, err
		}
	}

	return &domain.ProfileResponse{
		ID:          profile.ID,
		UserID:      user.ID,
		DisplayName: user.DisplayName,
		AvatarURL:   profile.AvatarURL,
		BannerURL:   profile.BannerURL,
		Bio:         profile.Bio,
		Headline:    profile.Headline,
		ThemeColor:  profile.ThemeColor,
		CustomCSS:   profile.CustomCSS,
		WebsiteURL:  profile.WebsiteURL,
		Location:    profile.Location,
		SocialLinks: profile.SocialLinks,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, dto *domain.UpdateProfileDTO) (*domain.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if dto.DisplayName != nil && *dto.DisplayName != "" {
		if err := s.userRepo.UpdateDisplayName(ctx, userID, *dto.DisplayName); err != nil {
			return nil, err
		}
		user.DisplayName = *dto.DisplayName
	}

	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrProfileNotFound) {
			now := time.Now().UTC()
			defaultTheme := "#3B82F6"
			if dto.ThemeColor != nil && *dto.ThemeColor != "" {
				defaultTheme = *dto.ThemeColor
			}
			profile = &domain.UserProfile{
				ID:          uuid.New().String(),
				UserID:      userID,
				AvatarURL:   dto.AvatarURL,
				BannerURL:   dto.BannerURL,
				Bio:         dto.Bio,
				Headline:    dto.Headline,
				ThemeColor:  defaultTheme,
				CustomCSS:   dto.CustomCSS,
				WebsiteURL:  dto.WebsiteURL,
				Location:    dto.Location,
				SocialLinks: dto.SocialLinks,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			if err := s.profileRepo.Create(ctx, profile); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if dto.AvatarURL != nil {
			profile.AvatarURL = dto.AvatarURL
		}
		if dto.BannerURL != nil {
			profile.BannerURL = dto.BannerURL
		}
		if dto.Bio != nil {
			profile.Bio = dto.Bio
		}
		if dto.Headline != nil {
			profile.Headline = dto.Headline
		}
		if dto.ThemeColor != nil && *dto.ThemeColor != "" {
			profile.ThemeColor = *dto.ThemeColor
		}
		if dto.CustomCSS != nil {
			profile.CustomCSS = dto.CustomCSS
		}
		if dto.WebsiteURL != nil {
			profile.WebsiteURL = dto.WebsiteURL
		}
		if dto.Location != nil {
			profile.Location = dto.Location
		}
		if dto.SocialLinks != nil {
			profile.SocialLinks = dto.SocialLinks
		}

		if err := s.profileRepo.Update(ctx, profile); err != nil {
			return nil, err
		}
	}

	return &domain.ProfileResponse{
		ID:          profile.ID,
		UserID:      user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   profile.AvatarURL,
		BannerURL:   profile.BannerURL,
		Bio:         profile.Bio,
		Headline:    profile.Headline,
		ThemeColor:  profile.ThemeColor,
		CustomCSS:   profile.CustomCSS,
		WebsiteURL:  profile.WebsiteURL,
		Location:    profile.Location,
		SocialLinks: profile.SocialLinks,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}
