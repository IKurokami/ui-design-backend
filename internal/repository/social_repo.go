package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"ui-design-backend/internal/domain"
)

var ErrCommentNotFound = errors.New("comment not found")

// Subscription visibility matches the public catalog. Drafts and incomplete
// books never become public through the follower or notification endpoints.
const followEligible = "s.status IN ('published','completed') AND (SELECT COUNT(*) FROM chapters c WHERE c.story_id=s.id AND c.has_content=1)>1"

func (r *StoryRepo) FollowIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT f.story_id FROM story_follows f JOIN stories s ON s.id=f.story_id
		WHERE f.user_id=$1 AND `+followEligible+` ORDER BY f.created_at DESC, f.story_id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, rows.Err()
}

func (r *StoryRepo) Follow(ctx context.Context, userID, storyID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO story_follows(user_id,story_id,created_at)
		SELECT $1,s.id,$3 FROM stories s WHERE s.id=$2 AND `+followEligible+`
		ON CONFLICT(user_id,story_id) DO NOTHING`, userID, storyID, time.Now().UTC())
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	var found int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM story_follows f JOIN stories s ON s.id=f.story_id
		WHERE f.user_id=$1 AND f.story_id=$2 AND `+followEligible, userID, storyID).Scan(&found)
	return found > 0, err
}
func (r *StoryRepo) Unfollow(ctx context.Context, userID, storyID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM story_follows WHERE user_id=$1 AND story_id=$2", userID, storyID)
	return err
}

func (r *StoryRepo) ListNotifications(ctx context.Context, userID string, limit int) ([]domain.StoryNotification, error) {
	if limit < 1 || limit > 100 {
		limit = 40
	}
	rows, err := r.db.QueryContext(ctx, `SELECT n.story_id,n.chapter_id,s.title,n.chapter_title,n.created_at,
		CASE WHEN n.read_at IS NULL THEN 0 ELSE 1 END FROM story_notifications n
		JOIN stories s ON s.id=n.story_id WHERE n.user_id=$1 AND s.status IN ('published','completed')
		ORDER BY n.created_at DESC,n.chapter_id DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.StoryNotification{}
	for rows.Next() {
		var n domain.StoryNotification
		var read int
		if err := rows.Scan(&n.StoryID, &n.ChapterID, &n.StoryTitle, &n.ChapterTitle, &n.CreatedAt, &read); err != nil {
			return nil, err
		}
		n.Read = read != 0
		items = append(items, n)
	}
	return items, rows.Err()
}
func (r *StoryRepo) MarkNotificationRead(ctx context.Context, userID, chapterID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE story_notifications SET read_at=$1 WHERE user_id=$2 AND chapter_id=$3 AND read_at IS NULL", time.Now().UTC(), userID, chapterID)
	return err
}

func scanSocialComments(rows *sql.Rows) ([]domain.StoryComment, error) {
	items := []domain.StoryComment{}
	for rows.Next() {
		var c domain.StoryComment
		if err := rows.Scan(&c.ID, &c.StoryID, &c.UserID, &c.AuthorName, &c.Body, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

const commentFields = "c.id,c.story_id,c.user_id,COALESCE(u.display_name,''),c.body,c.status,c.created_at"

func (r *StoryRepo) ListComments(ctx context.Context, storyID, userID string) ([]domain.StoryComment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+commentFields+` FROM story_comments c JOIN users u ON u.id=c.user_id
		JOIN stories s ON s.id=c.story_id
		WHERE c.story_id=$1 AND s.status IN ('published','completed')
		AND (c.status='approved' OR (c.user_id=$2 AND c.status='pending'))
		ORDER BY c.created_at DESC LIMIT 100`, storyID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSocialComments(rows)
}
func (r *StoryRepo) CreateComment(ctx context.Context, storyID, userID, body string) (*domain.StoryComment, error) {
	comment := &domain.StoryComment{ID: uuid.NewString(), StoryID: storyID, UserID: userID, Body: body, Status: "pending", CreatedAt: time.Now().UTC()}
	// Rate-limit at the database boundary and never expose draft content.
	res, err := r.db.ExecContext(ctx, `INSERT INTO story_comments(id,story_id,user_id,body,status,created_at)
		SELECT $1,s.id,$3,$4,'pending',$5 FROM stories s
		WHERE s.id=$2 AND s.status IN ('published','completed')
		AND (SELECT COUNT(*) FROM story_comments c WHERE c.user_id=$3 AND c.story_id=$2 AND c.created_at>$6)<5`,
		comment.ID, storyID, userID, body, comment.CreatedAt, comment.CreatedAt.Add(-time.Hour))
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrCommentNotFound
	}
	var name string
	if err := r.db.QueryRowContext(ctx, "SELECT display_name FROM users WHERE id=$1", userID).Scan(&name); err != nil {
		return nil, err
	}
	comment.AuthorName = name
	return comment, nil
}
func (r *StoryRepo) DeleteComment(ctx context.Context, userID, commentID string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM story_comments WHERE id=$1 AND user_id=$2", commentID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrCommentNotFound
	}
	return nil
}
func (r *StoryRepo) ListPendingComments(ctx context.Context, limit int) ([]domain.StoryComment, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+commentFields+` FROM story_comments c JOIN users u ON u.id=c.user_id
		WHERE c.status='pending' ORDER BY c.created_at ASC,c.id ASC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSocialComments(rows)
}
func (r *StoryRepo) ModerateComment(ctx context.Context, id, status string) error {
	if status != "approved" && status != "rejected" {
		return ErrCommentNotFound
	}
	res, err := r.db.ExecContext(ctx, "UPDATE story_comments SET status=$1 WHERE id=$2 AND status='pending'", status, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrCommentNotFound
	}
	return nil
}

func (r *StoryRepo) MutedNotificationStoryIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT story_id FROM story_notification_mutes WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
func (r *StoryRepo) MuteStoryNotifications(ctx context.Context, userID, storyID string) error {
    // Muting means permanently blocking new notifications until re-enabled,
    // and immediately cancelling all existing alerts for this story.
    // Both actions use one transaction to avoid a half-applied cancellation.
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()
    _, err = tx.ExecContext(ctx, `INSERT INTO story_notification_mutes(user_id, story_id, created_at)
        SELECT $1, s.id, $3 FROM stories s
        WHERE s.id = $2 AND s.status IN ('published','completed')
        ON CONFLICT(user_id,story_id) DO NOTHING`, userID, storyID, time.Now().UTC())
    if err != nil { return err }
    _, err = tx.ExecContext(ctx, `DELETE FROM story_notifications
        WHERE user_id = $1 AND story_id = $2`, userID, storyID)
    if err != nil { return err }
    return tx.Commit()
}
func (r *StoryRepo) UnmuteStoryNotifications(ctx context.Context, userID, storyID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM story_notification_mutes WHERE user_id = $1 AND story_id = $2", userID, storyID)
	return err
}

func (r *StoryRepo) ListMutedStories(ctx context.Context, userID string) ([]domain.MutedStory, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT m.story_id, s.title
        FROM story_notification_mutes m
        JOIN stories s ON s.id = m.story_id
        WHERE m.user_id = $1
        ORDER BY m.created_at DESC, m.story_id DESC`, userID)
    if err != nil { return nil, err }
    defer rows.Close()
    items := []domain.MutedStory{}
    for rows.Next() {
        var item domain.MutedStory
        if err := rows.Scan(&item.StoryID, &item.StoryTitle); err != nil { return nil, err }
        items = append(items, item)
    }
    return items, rows.Err()
}

// Dismissal changes only this user's inbox; new chapters from the same
// story will still notify them unless they explicitly mute the story.
func (r *StoryRepo) DismissNotification(ctx context.Context, userID, chapterID string) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM story_notifications
        WHERE user_id = $1 AND chapter_id = $2`, userID, chapterID)
    return err
}
func (r *StoryRepo) ClearNotifications(ctx context.Context, userID string) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM story_notifications WHERE user_id = $1`, userID)
    return err
}
