package domain

import (
	"context"
	"time"
)

type UserProfile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	AvatarURL   *string   `json:"avatarUrl"`
	BannerURL   *string   `json:"bannerUrl"`
	Bio         *string   `json:"bio"`
	Headline    *string   `json:"headline"`
	ThemeColor  string    `json:"themeColor"`
	CustomCSS   *string   `json:"customCss"`
	WebsiteURL  *string   `json:"websiteUrl"`
	Location    *string   `json:"location"`
	SocialLinks *string   `json:"socialLinks"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProfileResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email,omitempty"`
	AvatarURL   *string   `json:"avatarUrl"`
	BannerURL   *string   `json:"bannerUrl"`
	Bio         *string   `json:"bio"`
	Headline    *string   `json:"headline"`
	ThemeColor  string    `json:"themeColor"`
	CustomCSS   *string   `json:"customCss"`
	WebsiteURL  *string   `json:"websiteUrl"`
	Location    *string   `json:"location"`
	SocialLinks *string   `json:"socialLinks"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateProfileDTO struct {
	DisplayName *string `json:"displayName,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	BannerURL   *string `json:"bannerUrl,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Headline    *string `json:"headline,omitempty"`
	ThemeColor  *string `json:"themeColor,omitempty"`
	CustomCSS   *string `json:"customCss,omitempty"`
	WebsiteURL  *string `json:"websiteUrl,omitempty"`
	Location    *string `json:"location,omitempty"`
	SocialLinks *string `json:"socialLinks,omitempty"`
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *UserProfile) error
	FindByUserID(ctx context.Context, userID string) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
}
