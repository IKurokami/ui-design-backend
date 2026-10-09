package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ui-design-backend/internal/domain"
)

var ErrProfileNotFound = errors.New("profile not found")

type ProfileRepo struct {
	db *DB
}

func NewProfileRepo(db *DB) *ProfileRepo {
	return &ProfileRepo{db: db}
}

func (r *ProfileRepo) Create(ctx context.Context, p *domain.UserProfile) error {
	query := `INSERT INTO user_profiles (id, user_id, avatar_url, banner_url, bio, headline, theme_color, custom_css, website_url, location, social_links, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO user_profiles (id, user_id, avatar_url, banner_url, bio, headline, theme_color, custom_css, website_url, location, social_links, created_at, updated_at) 
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.UserID,
		p.AvatarURL,
		p.BannerURL,
		p.Bio,
		p.Headline,
		p.ThemeColor,
		p.CustomCSS,
		p.WebsiteURL,
		p.Location,
		p.SocialLinks,
		p.CreatedAt,
		p.UpdatedAt,
	)
	return err
}

func (r *ProfileRepo) FindByUserID(ctx context.Context, userID string) (*domain.UserProfile, error) {
	query := `SELECT id, user_id, avatar_url, banner_url, bio, headline, theme_color, custom_css, website_url, location, social_links, created_at, updated_at 
	          FROM user_profiles WHERE user_id = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, user_id, avatar_url, banner_url, bio, headline, theme_color, custom_css, website_url, location, social_links, created_at, updated_at 
		         FROM user_profiles WHERE user_id = ?`
	}

	row := r.db.QueryRowContext(ctx, query, userID)
	return r.scanProfile(row)
}

func (r *ProfileRepo) Update(ctx context.Context, p *domain.UserProfile) error {
	query := `UPDATE user_profiles SET 
	          avatar_url = $1, 
	          banner_url = $2, 
	          bio = $3, 
	          headline = $4, 
	          theme_color = $5, 
	          custom_css = $6, 
	          website_url = $7, 
	          location = $8, 
	          social_links = $9, 
	          updated_at = $10 
	          WHERE user_id = $11`
	if r.db.Driver == "sqlite" {
		query = `UPDATE user_profiles SET 
		          avatar_url = ?, 
		          banner_url = ?, 
		          bio = ?, 
		          headline = ?, 
		          theme_color = ?, 
		          custom_css = ?, 
		          website_url = ?, 
		          location = ?, 
		          social_links = ?, 
		          updated_at = ? 
		          WHERE user_id = ?`
	}

	p.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		p.AvatarURL,
		p.BannerURL,
		p.Bio,
		p.Headline,
		p.ThemeColor,
		p.CustomCSS,
		p.WebsiteURL,
		p.Location,
		p.SocialLinks,
		p.UpdatedAt,
		p.UserID,
	)
	return err
}

func (r *ProfileRepo) scanProfile(row *sql.Row) (*domain.UserProfile, error) {
	var p domain.UserProfile
	var avatarURL, bannerURL, bio, headline, customCSS, websiteURL, location, socialLinks sql.NullString

	err := row.Scan(
		&p.ID,
		&p.UserID,
		&avatarURL,
		&bannerURL,
		&bio,
		&headline,
		&p.ThemeColor,
		&customCSS,
		&websiteURL,
		&location,
		&socialLinks,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	if avatarURL.Valid {
		p.AvatarURL = &avatarURL.String
	}
	if bannerURL.Valid {
		p.BannerURL = &bannerURL.String
	}
	if bio.Valid {
		p.Bio = &bio.String
	}
	if headline.Valid {
		p.Headline = &headline.String
	}
	if customCSS.Valid {
		p.CustomCSS = &customCSS.String
	}
	if websiteURL.Valid {
		p.WebsiteURL = &websiteURL.String
	}
	if location.Valid {
		p.Location = &location.String
	}
	if socialLinks.Valid {
		p.SocialLinks = &socialLinks.String
	}

	return &p, nil
}
